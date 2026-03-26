package rbac

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	rbacdao "go-admin/internal/dao/rbac"
	"go-admin/internal/models"
	rbacsvc "go-admin/internal/services/rbac"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// PermissionController 负责权限的增删改查与变更通知。
type PermissionController struct {
	svc    *rbacsvc.PermissionService    // 权限业务服务
	msgSvc *rbacsvc.SystemMessageService // 系统消息推送服务
}

// permissionUpsertReq 定义权限新增/更新的请求载体。
type permissionUpsertReq struct {
	// 处理当前语句逻辑。
	Name string `json:"name" binding:"required"`
	// 处理当前语句逻辑。
	Code string `json:"code" binding:"required"`
	// 处理当前语句逻辑。
	Type string `json:"type"`
	// 处理当前语句逻辑。
	Method string `json:"method"`
	// 处理当前语句逻辑。
	Path string `json:"path"`
}

// NewPermissionController 构造权限控制器并注入依赖。
func NewPermissionController(db *gorm.DB) *PermissionController {
	// 返回当前处理结果。
	return &PermissionController{
		svc:    rbacsvc.NewPermissionService(rbacdao.NewPermissionDAO(db)),       // 注入权限 DAO
		msgSvc: rbacsvc.NewSystemMessageService(rbacdao.NewSystemMessageDao(db)), // 注入系统消息 DAO
	}
}

// List 返回权限列表。
func (pc *PermissionController) List(c *gin.Context) {
	// 调用服务层获取全部权限
	items, err := pc.svc.List(c.Request.Context())
	// 判断条件并进入对应分支逻辑。
	if err != nil {
		// 服务异常时直接返回 500
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		// 返回当前处理结果。
		return
	}
	// 返回权限列表
	c.JSON(http.StatusOK, gin.H{"data": items})
}

// Create 新增权限。
func (pc *PermissionController) Create(c *gin.Context) {
	// 声明当前变量。
	var req permissionUpsertReq
	// 绑定并校验请求体
	if err := c.ShouldBindJSON(&req); err != nil {
		// 调用c.JSON完成当前处理。
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		// 返回当前处理结果。
		return
	}
	// 组装权限模型并做字段清理
	p := models.Permission{
		// 调用strings.TrimSpace完成当前处理。
		Name: strings.TrimSpace(req.Name),
		// 调用strings.TrimSpace完成当前处理。
		Code: strings.TrimSpace(req.Code),
		// 调用strings.TrimSpace完成当前处理。
		Type: strings.TrimSpace(req.Type),
		// 调用strings.ToUpper完成当前处理。
		Method: strings.ToUpper(strings.TrimSpace(req.Method)),
		// 调用strings.TrimSpace完成当前处理。
		Path: strings.TrimSpace(req.Path),
	}
	// 再次校验必填字段，防止空值写入
	if p.Name == "" || p.Code == "" {
		// 调用c.JSON完成当前处理。
		c.JSON(http.StatusBadRequest, gin.H{"error": "name/code不能为空"})
		// 返回当前处理结果。
		return
	}
	// 交给服务层创建
	if err := pc.svc.Create(c.Request.Context(), &p); err != nil {
		// 调用c.JSON完成当前处理。
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		// 返回当前处理结果。
		return
	}
	// 返回创建结果
	c.JSON(http.StatusOK, gin.H{"data": p})
}

// Get 返回单条权限信息。
func (pc *PermissionController) Get(c *gin.Context) {
	// 解析并校验路由参数
	id, err := strconv.Atoi(c.Param("id"))
	// 判断条件并进入对应分支逻辑。
	if err != nil || id <= 0 {
		// 调用c.JSON完成当前处理。
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid permission id"})
		// 返回当前处理结果。
		return
	}

	// 查询权限详情
	p, err := pc.svc.Get(c.Request.Context(), uint(id))
	// 判断条件并进入对应分支逻辑。
	if err != nil {
		// 调用c.JSON完成当前处理。
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		// 返回当前处理结果。
		return
	}
	// 返回权限详情
	c.JSON(http.StatusOK, gin.H{"data": p})
}

// Update 更新权限。
func (pc *PermissionController) Update(c *gin.Context) {
	// 解析并校验路由参数
	id, err := strconv.Atoi(c.Param("id"))
	// 判断条件并进入对应分支逻辑。
	if err != nil || id <= 0 {
		// 调用c.JSON完成当前处理。
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid permission id"})
		// 返回当前处理结果。
		return
	}

	// 声明当前变量。
	var req permissionUpsertReq
	// 绑定并校验请求体
	if err := c.ShouldBindJSON(&req); err != nil {
		// 调用c.JSON完成当前处理。
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		// 返回当前处理结果。
		return
	}

	// 读取旧数据用于更新
	p, err := pc.svc.Get(c.Request.Context(), uint(id))
	// 判断条件并进入对应分支逻辑。
	if err != nil {
		// 调用c.JSON完成当前处理。
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		// 返回当前处理结果。
		return
	}

	// 覆盖字段并做清理/规范化
	p.Name = strings.TrimSpace(req.Name)
	// 更新当前变量或字段值。
	p.Code = strings.TrimSpace(req.Code)
	// 更新当前变量或字段值。
	p.Type = strings.TrimSpace(req.Type)
	// 更新当前变量或字段值。
	p.Method = strings.ToUpper(strings.TrimSpace(req.Method))
	// 更新当前变量或字段值。
	p.Path = strings.TrimSpace(req.Path)
	// 必填字段校验
	if p.Name == "" || p.Code == "" {
		// 调用c.JSON完成当前处理。
		c.JSON(http.StatusBadRequest, gin.H{"error": "name/code不能为空"})
		// 返回当前处理结果。
		return
	}

	// 写入更新
	if err := pc.svc.Update(c.Request.Context(), p); err != nil {
		// 调用c.JSON完成当前处理。
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		// 返回当前处理结果。
		return
	}

	// 变更后通知拥有该权限的用户
	_ = pc.msgSvc.PushToUsersByPermissionIDs(
		// 调用c.Request.Context完成当前处理。
		c.Request.Context(),
		// 处理当前语句逻辑。
		[]uint{p.ID},
		// 处理当前语句逻辑。
		"权限策略变更通知",
		// 调用fmt.Sprintf完成当前处理。
		fmt.Sprintf("权限【%s】已被管理员更新，请确认你的菜单与接口访问范围。", p.Name),
		// 处理当前语句逻辑。
		"warning",
		// 处理当前语句逻辑。
		"permission",
		// 处理当前语句逻辑。
		p.ID,
		// 调用c.GetUint完成当前处理。
		c.GetUint("uid"),
	)

	// 返回更新后的权限
	c.JSON(http.StatusOK, gin.H{"data": p})
}

// Delete 删除权限。
func (pc *PermissionController) Delete(c *gin.Context) {
	// 解析并校验路由参数
	id, err := strconv.Atoi(c.Param("id"))
	// 判断条件并进入对应分支逻辑。
	if err != nil || id <= 0 {
		// 调用c.JSON完成当前处理。
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid permission id"})
		// 返回当前处理结果。
		return
	}

	// 先查询权限详情用于消息提示
	permission, err := pc.svc.Get(c.Request.Context(), uint(id))
	// 判断条件并进入对应分支逻辑。
	if err != nil {
		// 调用c.JSON完成当前处理。
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		// 返回当前处理结果。
		return
	}

	// 查询受影响用户列表
	affectedUserIDs, _ := pc.msgSvc.ListUserIDsByPermissionIDs(c.Request.Context(), []uint{uint(id)})

	// 执行删除
	if err := pc.svc.Delete(c.Request.Context(), uint(id)); err != nil {
		// 调用c.JSON完成当前处理。
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		// 返回当前处理结果。
		return
	}

	// 删除后通知相关用户
	_ = pc.msgSvc.PushToUsers(
		// 调用c.Request.Context完成当前处理。
		c.Request.Context(),
		// 处理当前语句逻辑。
		affectedUserIDs,
		// 处理当前语句逻辑。
		"权限删除通知",
		// 调用fmt.Sprintf完成当前处理。
		fmt.Sprintf("权限【%s】已被管理员删除，你的可访问资源可能发生变化。", permission.Name),
		// 处理当前语句逻辑。
		"warning",
		// 处理当前语句逻辑。
		"permission",
		// 调用uint完成当前处理。
		uint(id),
		// 调用c.GetUint完成当前处理。
		c.GetUint("uid"),
	)

	// 返回删除结果
	c.JSON(http.StatusOK, gin.H{"msg": "deleted"})
}
