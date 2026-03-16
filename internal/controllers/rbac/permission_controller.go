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

// PermissionController 负责权限的增删改查与变更通知。
type PermissionController struct {
	svc    *rbacsvc.PermissionService    // 权限业务服务
	msgSvc *rbacsvc.SystemMessageService // 系统消息推送服务
}

// permissionUpsertReq 定义权限新增/更新的请求载体。
type permissionUpsertReq struct {
	Name   string `json:"name" binding:"required"`
	Code   string `json:"code" binding:"required"`
	Type   string `json:"type"`
	Method string `json:"method"`
	Path   string `json:"path"`
}

// NewPermissionController 构造权限控制器并注入依赖。
func NewPermissionController(db *gorm.DB) *PermissionController {
	return &PermissionController{
		svc:    rbacsvc.NewPermissionService(rbacdao.NewPermissionDAO(db)),       // 注入权限 DAO
		msgSvc: rbacsvc.NewSystemMessageService(rbacdao.NewSystemMessageDao(db)), // 注入系统消息 DAO
	}
}

// List 返回权限列表。
func (pc *PermissionController) List(c *gin.Context) {
	// 调用服务层获取全部权限
	items, err := pc.svc.List(c.Request.Context())
	if err != nil {
		// 服务异常时直接返回 500
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	// 返回权限列表
	c.JSON(http.StatusOK, gin.H{"data": items})
}

// Create 新增权限。
func (pc *PermissionController) Create(c *gin.Context) {
	var req permissionUpsertReq
	// 绑定并校验请求体
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	// 组装权限模型并做字段清理
	p := models.Permission{
		Name:   strings.TrimSpace(req.Name),
		Code:   strings.TrimSpace(req.Code),
		Type:   strings.TrimSpace(req.Type),
		Method: strings.ToUpper(strings.TrimSpace(req.Method)),
		Path:   strings.TrimSpace(req.Path),
	}
	// 再次校验必填字段，防止空值写入
	if p.Name == "" || p.Code == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "name/code不能为空"})
		return
	}
	// 交给服务层创建
	if err := pc.svc.Create(c.Request.Context(), &p); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	// 返回创建结果
	c.JSON(http.StatusOK, gin.H{"data": p})
}

// Get 返回单条权限信息。
func (pc *PermissionController) Get(c *gin.Context) {
	// 解析并校验路由参数
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil || id <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid permission id"})
		return
	}

	// 查询权限详情
	p, err := pc.svc.Get(c.Request.Context(), uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	// 返回权限详情
	c.JSON(http.StatusOK, gin.H{"data": p})
}

// Update 更新权限。
func (pc *PermissionController) Update(c *gin.Context) {
	// 解析并校验路由参数
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil || id <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid permission id"})
		return
	}

	var req permissionUpsertReq
	// 绑定并校验请求体
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 读取旧数据用于更新
	p, err := pc.svc.Get(c.Request.Context(), uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	// 覆盖字段并做清理/规范化
	p.Name = strings.TrimSpace(req.Name)
	p.Code = strings.TrimSpace(req.Code)
	p.Type = strings.TrimSpace(req.Type)
	p.Method = strings.ToUpper(strings.TrimSpace(req.Method))
	p.Path = strings.TrimSpace(req.Path)
	// 必填字段校验
	if p.Name == "" || p.Code == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "name/code不能为空"})
		return
	}

	// 写入更新
	if err := pc.svc.Update(c.Request.Context(), p); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 变更后通知拥有该权限的用户
	_ = pc.msgSvc.PushToUsersByPermissionIDs(
		c.Request.Context(),
		[]uint{p.ID},
		"权限策略变更通知",
		fmt.Sprintf("权限【%s】已被管理员更新，请确认你的菜单与接口访问范围。", p.Name),
		"warning",
		"permission",
		p.ID,
		c.GetUint("uid"),
	)

	// 返回更新后的权限
	c.JSON(http.StatusOK, gin.H{"data": p})
}

// Delete 删除权限。
func (pc *PermissionController) Delete(c *gin.Context) {
	// 解析并校验路由参数
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil || id <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid permission id"})
		return
	}

	// 先查询权限详情用于消息提示
	permission, err := pc.svc.Get(c.Request.Context(), uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	// 查询受影响用户列表
	affectedUserIDs, _ := pc.msgSvc.ListUserIDsByPermissionIDs(c.Request.Context(), []uint{uint(id)})

	// 执行删除
	if err := pc.svc.Delete(c.Request.Context(), uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 删除后通知相关用户
	_ = pc.msgSvc.PushToUsers(
		c.Request.Context(),
		affectedUserIDs,
		"权限删除通知",
		fmt.Sprintf("权限【%s】已被管理员删除，你的可访问资源可能发生变化。", permission.Name),
		"warning",
		"permission",
		uint(id),
		c.GetUint("uid"),
	)

	// 返回删除结果
	c.JSON(http.StatusOK, gin.H{"msg": "deleted"})
}
