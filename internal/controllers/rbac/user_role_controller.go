package rbac

import (
	"net/http"
	"strconv"
	"strings"

	rbacdao "go-admin/internal/dao/rbac"
	"go-admin/internal/middleware"
	"go-admin/internal/models"
	rbacsvc "go-admin/internal/services/rbac"
	tokenjwt "go-admin/internal/token/jwt"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// UserRoleController 负责用户角色绑定与变更通知。
type UserRoleController struct {
	svc      *rbacsvc.UserRoleService      // 用户角色关系服务
	msgSvc   *rbacsvc.SystemMessageService // 系统消息服务
	tokenMgr *tokenjwt.Manager             // Token 管理器（用于权限缓存失效）
}

// NewUserRoleController 构造控制器并注入依赖。
func NewUserRoleController(db *gorm.DB, mgr ...*tokenjwt.Manager) *UserRoleController {
	// 初始化 DAO
	d := rbacdao.NewUserRoleDao(db)
	// 初始化 Service
	s := rbacsvc.NewUserRoleService(d)
	// 初始化系统消息 Service
	msgSvc := rbacsvc.NewSystemMessageService(rbacdao.NewSystemMessageDao(db))
	// 可选注入 tokenMgr（不影响现有调用方式）
	var tokenMgr *tokenjwt.Manager
	if len(mgr) > 0 {
		tokenMgr = mgr[0]
	}
	return &UserRoleController{svc: s, msgSvc: msgSvc, tokenMgr: tokenMgr}
}

// bindRolesReq 定义用户角色绑定/增删的请求载体。
type bindRolesReq struct {
	// 处理当前语句逻辑。
	UserID uint `json:"user_id" binding:"required"`
	// 处理当前语句逻辑。
	RoleIDs []uint `json:"role_ids"`
}

// BindRoles 全量替换用户的角色列表。
func (uc *UserRoleController) BindRoles(c *gin.Context) {
	// 声明当前变量。
	var req bindRolesReq
	// 绑定并校验请求体
	if err := c.ShouldBindJSON(&req); err != nil {
		// 调用c.JSON完成当前处理。
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		// 返回当前处理结果。
		return
	}
	// 安全防护：禁止修改当前登录用户自身角色，避免误操作导致权限锁死或提权。
	if req.UserID == c.GetUint("uid") {
		// 调用c.JSON完成当前处理。
		c.JSON(http.StatusForbidden, gin.H{"error": "不可修改当前登录用户角色"})
		// 返回当前处理结果。
		return
	}

	// 执行全量绑定
	if err := uc.svc.BindRoles(c.Request.Context(), req.UserID, req.RoleIDs); err != nil {
		// 调用c.JSON完成当前处理。
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		// 返回当前处理结果。
		return
	}

	// 用户角色变更后立即清除该用户的 RBAC 权限缓存，
	// 保证用户下次请求时能加载到最新角色和权限配置。
	middleware.InvalidateUserPermCache(uc.tokenMgr, req.UserID)

	// 查询更新后的角色列表，用于消息提示
	roles, _ := uc.svc.GetUserRoles(c.Request.Context(), req.UserID)
	// 通知目标用户
	_ = uc.msgSvc.PushToUser(
		// 调用c.Request.Context完成当前处理。
		c.Request.Context(),
		// 处理当前语句逻辑。
		req.UserID,
		// 处理当前语句逻辑。
		"角色分配变更通知",
		// 调用joinRoleNames完成当前处理。
		"管理员已调整你的角色分配，当前角色："+joinRoleNames(roles),
		// 处理当前语句逻辑。
		"warning",
		// 处理当前语句逻辑。
		"user_role",
		// 处理当前语句逻辑。
		req.UserID,
		// 调用c.GetUint完成当前处理。
		c.GetUint("uid"),
	)

	// 返回绑定成功
	c.JSON(http.StatusOK, gin.H{"msg": "roles bound"})
}

// AddRoles 为用户追加角色（非全量替换）。
func (uc *UserRoleController) AddRoles(c *gin.Context) {
	// 声明当前变量。
	var req bindRolesReq
	// 绑定并校验请求体
	if err := c.ShouldBindJSON(&req); err != nil {
		// 调用c.JSON完成当前处理。
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		// 返回当前处理结果。
		return
	}
	// 校验角色列表不为空
	if len(req.RoleIDs) == 0 {
		// 调用c.JSON完成当前处理。
		c.JSON(http.StatusBadRequest, gin.H{"error": "role_ids不能为空"})
		// 返回当前处理结果。
		return
	}
	// 禁止修改当前登录用户角色
	if req.UserID == c.GetUint("uid") {
		// 调用c.JSON完成当前处理。
		c.JSON(http.StatusForbidden, gin.H{"error": "不可修改当前登录用户角色"})
		// 返回当前处理结果。
		return
	}

	// 追加角色
	if err := uc.svc.AddRoles(c.Request.Context(), req.UserID, req.RoleIDs); err != nil {
		// 调用c.JSON完成当前处理。
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		// 返回当前处理结果。
		return
	}

	// 角色变更后清除权限缓存
	middleware.InvalidateUserPermCache(uc.tokenMgr, req.UserID)

	// 查询更新后的角色列表
	roles, _ := uc.svc.GetUserRoles(c.Request.Context(), req.UserID)
	// 通知目标用户
	_ = uc.msgSvc.PushToUser(
		// 调用c.Request.Context完成当前处理。
		c.Request.Context(),
		// 处理当前语句逻辑。
		req.UserID,
		// 处理当前语句逻辑。
		"角色新增通知",
		// 调用joinRoleNames完成当前处理。
		"管理员已给你新增角色，当前角色："+joinRoleNames(roles),
		// 处理当前语句逻辑。
		"info",
		// 处理当前语句逻辑。
		"user_role",
		// 处理当前语句逻辑。
		req.UserID,
		// 调用c.GetUint完成当前处理。
		c.GetUint("uid"),
	)

	// 返回处理结果
	c.JSON(http.StatusOK, gin.H{"msg": "roles added"})
}

// RemoveRoles 从用户角色中移除指定角色。
func (uc *UserRoleController) RemoveRoles(c *gin.Context) {
	// 声明当前变量。
	var req bindRolesReq
	// 绑定并校验请求体
	if err := c.ShouldBindJSON(&req); err != nil {
		// 调用c.JSON完成当前处理。
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		// 返回当前处理结果。
		return
	}
	// 校验角色列表不为空
	if len(req.RoleIDs) == 0 {
		// 调用c.JSON完成当前处理。
		c.JSON(http.StatusBadRequest, gin.H{"error": "role_ids不能为空"})
		// 返回当前处理结果。
		return
	}
	// 禁止修改当前登录用户角色
	if req.UserID == c.GetUint("uid") {
		// 调用c.JSON完成当前处理。
		c.JSON(http.StatusForbidden, gin.H{"error": "不可修改当前登录用户角色"})
		// 返回当前处理结果。
		return
	}

	// 移除角色
	if err := uc.svc.RemoveRoles(c.Request.Context(), req.UserID, req.RoleIDs); err != nil {
		// 调用c.JSON完成当前处理。
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		// 返回当前处理结果。
		return
	}

	// 角色移除后立即清除该用户权限缓存，确保下次请求不再持有被撤销的权限
	middleware.InvalidateUserPermCache(uc.tokenMgr, req.UserID)

	// 查询更新后的角色列表
	roles, _ := uc.svc.GetUserRoles(c.Request.Context(), req.UserID)
	// 通知目标用户
	_ = uc.msgSvc.PushToUser(
		// 调用c.Request.Context完成当前处理。
		c.Request.Context(),
		// 处理当前语句逻辑。
		req.UserID,
		// 处理当前语句逻辑。
		"角色移除通知",
		// 调用joinRoleNames完成当前处理。
		"管理员已移除你的部分角色，当前角色："+joinRoleNames(roles),
		// 处理当前语句逻辑。
		"warning",
		// 处理当前语句逻辑。
		"user_role",
		// 处理当前语句逻辑。
		req.UserID,
		// 调用c.GetUint完成当前处理。
		c.GetUint("uid"),
	)

	// 返回处理结果
	c.JSON(http.StatusOK, gin.H{"msg": "roles removed"})
}

// GetUserRoles 获取指定用户的角色列表。
func (uc *UserRoleController) GetUserRoles(c *gin.Context) {
	// 兼容 /:id 或 query 参数 user_id
	idStr := c.Param("id")
	// 判断条件并进入对应分支逻辑。
	if idStr == "" {
		// 更新当前变量或字段值。
		idStr = c.Query("user_id")
	}
	// 校验用户 ID 是否存在
	if idStr == "" {
		// 调用c.JSON完成当前处理。
		c.JSON(http.StatusBadRequest, gin.H{"error": "user id required"})
		// 返回当前处理结果。
		return
	}
	// 转换为整数 ID
	id64, err := strconv.ParseUint(idStr, 10, 64)
	// 判断条件并进入对应分支逻辑。
	if err != nil {
		// 调用c.JSON完成当前处理。
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user id"})
		// 返回当前处理结果。
		return
	}
	// 查询角色列表
	roles, err := uc.svc.GetUserRoles(c.Request.Context(), uint(id64))
	// 判断条件并进入对应分支逻辑。
	if err != nil {
		// 调用c.JSON完成当前处理。
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		// 返回当前处理结果。
		return
	}
	// 返回数据
	c.JSON(http.StatusOK, gin.H{"data": roles})
}

// joinRoleNames 将角色名称拼接为友好字符串。
func joinRoleNames(roles []models.Role) string {
	// 无角色时返回“无”
	if len(roles) == 0 {
		return "无"
	}
	// 过滤空角色名
	names := make([]string, 0, len(roles))
	// 循环处理当前数据集合。
	for _, role := range roles {
		// 判断条件并进入对应分支逻辑。
		if strings.TrimSpace(role.Name) == "" {
			// 处理当前语句逻辑。
			continue
		}
		// 更新当前变量或字段值。
		names = append(names, role.Name)
	}
	// 二次判空
	if len(names) == 0 {
		return "无"
	}
	// 用中文顿号拼接
	return strings.Join(names, "、")
}
