package rbac

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	rbacdao "go-admin-full/internal/dao/rbac"
	"go-admin-full/internal/models"
	rbacsvc "go-admin-full/internal/services/rbac"
	"gorm.io/gorm"
)

// RoleController 负责角色管理及角色权限绑定。
type RoleController struct {
	svc         *rbacsvc.RoleService           // 角色业务服务
	rolePermSvc *rbacsvc.RolePermissionService // 角色-权限绑定服务
	msgSvc      *rbacsvc.SystemMessageService  // 系统消息推送服务
}

// roleUpsertReq 定义角色新增/更新的请求载体。
type roleUpsertReq struct {
	// 处理当前语句逻辑。
	Name string `json:"name" binding:"required"`
	// 处理当前语句逻辑。
	Code string `json:"code" binding:"required"`
}

// NewRoleController 构造角色控制器并注入依赖。
func NewRoleController(db *gorm.DB) *RoleController {
	// 初始化角色服务
	svc := rbacsvc.NewRoleService(rbacdao.NewRoleDAO(db))
	// 初始化角色权限绑定服务
	rolePermSvc := rbacsvc.NewRolePermissionService(rbacdao.NewRolePermissionDao(db))
	// 初始化系统消息服务
	msgSvc := rbacsvc.NewSystemMessageService(rbacdao.NewSystemMessageDao(db))
	// 返回当前处理结果。
	return &RoleController{
		// 处理当前语句逻辑。
		svc: svc,
		// 处理当前语句逻辑。
		rolePermSvc: rolePermSvc,
		// 处理当前语句逻辑。
		msgSvc: msgSvc,
	}
}

// List 返回角色列表。
func (rc *RoleController) List(c *gin.Context) {
	// 获取全部角色
	roles, err := rc.svc.List(c.Request.Context())
	// 判断条件并进入对应分支逻辑。
	if err != nil {
		// 调用c.JSON完成当前处理。
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		// 返回当前处理结果。
		return
	}
	// 返回列表
	c.JSON(http.StatusOK, gin.H{"data": roles})
}

// Create 新增角色。
func (rc *RoleController) Create(c *gin.Context) {
	// 声明当前变量。
	var req roleUpsertReq
	// 绑定并校验请求体
	if err := c.ShouldBindJSON(&req); err != nil {
		// 调用c.JSON完成当前处理。
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		// 返回当前处理结果。
		return
	}

	// 组装角色模型并清理字段
	r := models.Role{
		// 调用strings.TrimSpace完成当前处理。
		Name: strings.TrimSpace(req.Name),
		// 调用strings.TrimSpace完成当前处理。
		Code: strings.TrimSpace(req.Code),
	}
	// 必填字段校验
	if r.Name == "" || r.Code == "" {
		// 调用c.JSON完成当前处理。
		c.JSON(http.StatusBadRequest, gin.H{"error": "name/code不能为空"})
		// 返回当前处理结果。
		return
	}
	// 写入数据库
	if err := rc.svc.Create(c.Request.Context(), &r); err != nil {
		// 调用c.JSON完成当前处理。
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		// 返回当前处理结果。
		return
	}
	// 返回结果
	c.JSON(http.StatusOK, gin.H{"data": r})
}

// Get 返回单个角色信息。
func (rc *RoleController) Get(c *gin.Context) {
	// 解析并校验角色 ID
	idStr := c.Param("id")
	// 定义并初始化当前变量。
	id, err := strconv.Atoi(idStr)
	// 判断条件并进入对应分支逻辑。
	if err != nil || id <= 0 {
		// 调用c.JSON完成当前处理。
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid role id"})
		// 返回当前处理结果。
		return
	}

	// 查询角色详情
	r, err := rc.svc.Get(c.Request.Context(), uint(id))
	// 判断条件并进入对应分支逻辑。
	if err != nil {
		// 调用c.JSON完成当前处理。
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		// 返回当前处理结果。
		return
	}
	// 返回角色详情
	c.JSON(http.StatusOK, gin.H{"data": r})
}

// Update 更新角色。
func (rc *RoleController) Update(c *gin.Context) {
	// 解析并校验角色 ID
	idStr := c.Param("id")
	// 定义并初始化当前变量。
	id, err := strconv.Atoi(idStr)
	// 判断条件并进入对应分支逻辑。
	if err != nil || id <= 0 {
		// 调用c.JSON完成当前处理。
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid role id"})
		// 返回当前处理结果。
		return
	}

	// 声明当前变量。
	var req roleUpsertReq
	// 绑定并校验请求体
	if err := c.ShouldBindJSON(&req); err != nil {
		// 调用c.JSON完成当前处理。
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		// 返回当前处理结果。
		return
	}

	// 读取旧数据用于更新
	role, err := rc.svc.Get(c.Request.Context(), uint(id))
	// 判断条件并进入对应分支逻辑。
	if err != nil {
		// 调用c.JSON完成当前处理。
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		// 返回当前处理结果。
		return
	}

	// 更新字段并清理
	role.Name = strings.TrimSpace(req.Name)
	// 更新当前变量或字段值。
	role.Code = strings.TrimSpace(req.Code)
	// 必填字段校验
	if role.Name == "" || role.Code == "" {
		// 调用c.JSON完成当前处理。
		c.JSON(http.StatusBadRequest, gin.H{"error": "name/code不能为空"})
		// 返回当前处理结果。
		return
	}

	// 提交更新
	if err := rc.svc.Update(c.Request.Context(), role); err != nil {
		// 调用c.JSON完成当前处理。
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		// 返回当前处理结果。
		return
	}

	// 角色信息变更后，通知该角色关联用户。
	_ = rc.msgSvc.PushToUsersByRoleIDs(
		// 调用c.Request.Context完成当前处理。
		c.Request.Context(),
		// 处理当前语句逻辑。
		[]uint{role.ID},
		// 处理当前语句逻辑。
		"角色信息更新通知",
		// 调用fmt.Sprintf完成当前处理。
		fmt.Sprintf("角色【%s】信息已更新，你的角色权限可能受到影响。", role.Name),
		// 处理当前语句逻辑。
		"warning",
		// 处理当前语句逻辑。
		"role",
		// 处理当前语句逻辑。
		role.ID,
		// 调用c.GetUint完成当前处理。
		c.GetUint("uid"),
	)

	// 调用c.JSON完成当前处理。
	c.JSON(http.StatusOK, gin.H{"data": role})
}

// Delete 删除角色。
func (rc *RoleController) Delete(c *gin.Context) {
	// 解析并校验角色 ID
	idStr := c.Param("id")
	// 定义并初始化当前变量。
	id, err := strconv.Atoi(idStr)
	// 判断条件并进入对应分支逻辑。
	if err != nil || id <= 0 {
		// 调用c.JSON完成当前处理。
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid role id"})
		// 返回当前处理结果。
		return
	}

	// 安全防护：禁止删除内置超级管理员角色
	role, err := rc.svc.Get(c.Request.Context(), uint(id))
	// 判断条件并进入对应分支逻辑。
	if err != nil {
		// 调用c.JSON完成当前处理。
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		// 返回当前处理结果。
		return
	}
	// 判断条件并进入对应分支逻辑。
	if strings.EqualFold(role.Code, "admin") {
		// 调用c.JSON完成当前处理。
		c.JSON(http.StatusForbidden, gin.H{"error": "admin角色不可删除"})
		// 返回当前处理结果。
		return
	}

	// 删除角色前先缓存受影响用户，避免角色删除后关联关系被清空无法通知。
	affectedUserIDs, _ := rc.msgSvc.ListUserIDsByRoleIDs(c.Request.Context(), []uint{role.ID})

	// 执行删除
	if err := rc.svc.Delete(c.Request.Context(), uint(id)); err != nil {
		// 调用c.JSON完成当前处理。
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		// 返回当前处理结果。
		return
	}

	// 删除后通知相关用户
	_ = rc.msgSvc.PushToUsers(
		// 调用c.Request.Context完成当前处理。
		c.Request.Context(),
		// 处理当前语句逻辑。
		affectedUserIDs,
		// 处理当前语句逻辑。
		"角色删除通知",
		// 调用fmt.Sprintf完成当前处理。
		fmt.Sprintf("角色【%s】已被管理员删除，请检查你的账号权限是否仍满足业务需求。", role.Name),
		// 处理当前语句逻辑。
		"warning",
		// 处理当前语句逻辑。
		"role",
		// 处理当前语句逻辑。
		role.ID,
		// 调用c.GetUint完成当前处理。
		c.GetUint("uid"),
	)

	// 调用c.JSON完成当前处理。
	c.JSON(http.StatusOK, gin.H{"msg": "deleted"})
}

// BindPermissions 绑定角色权限（全量替换）。
func (rc *RoleController) BindPermissions(c *gin.Context) {
	// 解析并校验角色 ID
	idStr := c.Param("id")
	// 定义并初始化当前变量。
	id, err := strconv.Atoi(idStr)
	// 判断条件并进入对应分支逻辑。
	if err != nil || id <= 0 {
		// 调用c.JSON完成当前处理。
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid role id"})
		// 返回当前处理结果。
		return
	}

	// 声明当前变量。
	var req struct {
		// 处理当前语句逻辑。
		PermissionIDs []uint `json:"permission_ids"`
	}
	// 绑定并校验请求体
	if err := c.ShouldBindJSON(&req); err != nil {
		// 调用c.JSON完成当前处理。
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		// 返回当前处理结果。
		return
	}

	// 校验角色是否存在
	role, err := rc.svc.Get(c.Request.Context(), uint(id))
	// 判断条件并进入对应分支逻辑。
	if err != nil {
		// 调用c.JSON完成当前处理。
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		// 返回当前处理结果。
		return
	}

	// 执行权限绑定（全量替换）
	if err := rc.rolePermSvc.BindPermissions(c.Request.Context(), uint(id), req.PermissionIDs); err != nil {
		// 调用c.JSON完成当前处理。
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		// 返回当前处理结果。
		return
	}

	// 绑定完成后通知该角色的所有用户
	_ = rc.msgSvc.PushToUsersByRoleIDs(
		// 调用c.Request.Context完成当前处理。
		c.Request.Context(),
		// 调用uint完成当前处理。
		[]uint{uint(id)},
		// 处理当前语句逻辑。
		"角色权限变更通知",
		// 调用fmt.Sprintf完成当前处理。
		fmt.Sprintf("角色【%s】的权限策略已调整，你的菜单与接口访问范围可能发生变化。", role.Name),
		// 处理当前语句逻辑。
		"warning",
		// 处理当前语句逻辑。
		"role_permission",
		// 调用uint完成当前处理。
		uint(id),
		// 调用c.GetUint完成当前处理。
		c.GetUint("uid"),
	)

	// 返回绑定成功
	c.JSON(http.StatusOK, gin.H{"msg": "permissions bound"})
}

// GetPermissions 查询角色已绑定权限。
func (rc *RoleController) GetPermissions(c *gin.Context) {
	// 解析并校验角色 ID
	idStr := c.Param("id")
	// 定义并初始化当前变量。
	id, err := strconv.Atoi(idStr)
	// 判断条件并进入对应分支逻辑。
	if err != nil || id <= 0 {
		// 调用c.JSON完成当前处理。
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid role id"})
		// 返回当前处理结果。
		return
	}

	// 查询角色对应的权限列表
	perms, err := rc.rolePermSvc.GetRolePermissions(c.Request.Context(), uint(id))
	// 判断条件并进入对应分支逻辑。
	if err != nil {
		// 调用c.JSON完成当前处理。
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		// 返回当前处理结果。
		return
	}
	// 返回权限列表
	c.JSON(http.StatusOK, gin.H{"data": perms})
}
