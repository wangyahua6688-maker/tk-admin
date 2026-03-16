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
	Name string `json:"name" binding:"required"`
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
	return &RoleController{
		svc:         svc,
		rolePermSvc: rolePermSvc,
		msgSvc:      msgSvc,
	}
}

// List 返回角色列表。
func (rc *RoleController) List(c *gin.Context) {
	// 获取全部角色
	roles, err := rc.svc.List(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	// 返回列表
	c.JSON(http.StatusOK, gin.H{"data": roles})
}

// Create 新增角色。
func (rc *RoleController) Create(c *gin.Context) {
	var req roleUpsertReq
	// 绑定并校验请求体
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 组装角色模型并清理字段
	r := models.Role{
		Name: strings.TrimSpace(req.Name),
		Code: strings.TrimSpace(req.Code),
	}
	// 必填字段校验
	if r.Name == "" || r.Code == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "name/code不能为空"})
		return
	}
	// 写入数据库
	if err := rc.svc.Create(c.Request.Context(), &r); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	// 返回结果
	c.JSON(http.StatusOK, gin.H{"data": r})
}

// Get 返回单个角色信息。
func (rc *RoleController) Get(c *gin.Context) {
	// 解析并校验角色 ID
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil || id <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid role id"})
		return
	}

	// 查询角色详情
	r, err := rc.svc.Get(c.Request.Context(), uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	// 返回角色详情
	c.JSON(http.StatusOK, gin.H{"data": r})
}

// Update 更新角色。
func (rc *RoleController) Update(c *gin.Context) {
	// 解析并校验角色 ID
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil || id <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid role id"})
		return
	}

	var req roleUpsertReq
	// 绑定并校验请求体
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 读取旧数据用于更新
	role, err := rc.svc.Get(c.Request.Context(), uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	// 更新字段并清理
	role.Name = strings.TrimSpace(req.Name)
	role.Code = strings.TrimSpace(req.Code)
	// 必填字段校验
	if role.Name == "" || role.Code == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "name/code不能为空"})
		return
	}

	// 提交更新
	if err := rc.svc.Update(c.Request.Context(), role); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 角色信息变更后，通知该角色关联用户。
	_ = rc.msgSvc.PushToUsersByRoleIDs(
		c.Request.Context(),
		[]uint{role.ID},
		"角色信息更新通知",
		fmt.Sprintf("角色【%s】信息已更新，你的角色权限可能受到影响。", role.Name),
		"warning",
		"role",
		role.ID,
		c.GetUint("uid"),
	)

	c.JSON(http.StatusOK, gin.H{"data": role})
}

// Delete 删除角色。
func (rc *RoleController) Delete(c *gin.Context) {
	// 解析并校验角色 ID
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil || id <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid role id"})
		return
	}

	// 安全防护：禁止删除内置超级管理员角色
	role, err := rc.svc.Get(c.Request.Context(), uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	if strings.EqualFold(role.Code, "admin") {
		c.JSON(http.StatusForbidden, gin.H{"error": "admin角色不可删除"})
		return
	}

	// 删除角色前先缓存受影响用户，避免角色删除后关联关系被清空无法通知。
	affectedUserIDs, _ := rc.msgSvc.ListUserIDsByRoleIDs(c.Request.Context(), []uint{role.ID})

	// 执行删除
	if err := rc.svc.Delete(c.Request.Context(), uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 删除后通知相关用户
	_ = rc.msgSvc.PushToUsers(
		c.Request.Context(),
		affectedUserIDs,
		"角色删除通知",
		fmt.Sprintf("角色【%s】已被管理员删除，请检查你的账号权限是否仍满足业务需求。", role.Name),
		"warning",
		"role",
		role.ID,
		c.GetUint("uid"),
	)

	c.JSON(http.StatusOK, gin.H{"msg": "deleted"})
}

// BindPermissions 绑定角色权限（全量替换）。
func (rc *RoleController) BindPermissions(c *gin.Context) {
	// 解析并校验角色 ID
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil || id <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid role id"})
		return
	}

	var req struct {
		PermissionIDs []uint `json:"permission_ids"`
	}
	// 绑定并校验请求体
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 校验角色是否存在
	role, err := rc.svc.Get(c.Request.Context(), uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	// 执行权限绑定（全量替换）
	if err := rc.rolePermSvc.BindPermissions(c.Request.Context(), uint(id), req.PermissionIDs); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 绑定完成后通知该角色的所有用户
	_ = rc.msgSvc.PushToUsersByRoleIDs(
		c.Request.Context(),
		[]uint{uint(id)},
		"角色权限变更通知",
		fmt.Sprintf("角色【%s】的权限策略已调整，你的菜单与接口访问范围可能发生变化。", role.Name),
		"warning",
		"role_permission",
		uint(id),
		c.GetUint("uid"),
	)

	// 返回绑定成功
	c.JSON(http.StatusOK, gin.H{"msg": "permissions bound"})
}

// GetPermissions 查询角色已绑定权限。
func (rc *RoleController) GetPermissions(c *gin.Context) {
	// 解析并校验角色 ID
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil || id <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid role id"})
		return
	}

	// 查询角色对应的权限列表
	perms, err := rc.rolePermSvc.GetRolePermissions(c.Request.Context(), uint(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	// 返回权限列表
	c.JSON(http.StatusOK, gin.H{"data": perms})
}
