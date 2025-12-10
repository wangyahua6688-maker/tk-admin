package controllers

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"go-admin-full/internal/dao"
	"go-admin-full/internal/models"
	"go-admin-full/internal/services"
	"gorm.io/gorm"
)

type RoleController struct {
	svc         *services.RoleService
	rolePermSvc *services.RolePermissionService
}

type roleUpsertReq struct {
	Name string `json:"name" binding:"required"`
	Code string `json:"code" binding:"required"`
}

func NewRoleController(db *gorm.DB) *RoleController {
	svc := services.NewRoleService(dao.NewRoleDAO(db))
	rolePermSvc := services.NewRolePermissionService(dao.NewRolePermissionDao(db))
	return &RoleController{
		svc:         svc,
		rolePermSvc: rolePermSvc,
	}
}

func (rc *RoleController) List(c *gin.Context) {
	roles, err := rc.svc.List(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": roles})
}

func (rc *RoleController) Create(c *gin.Context) {
	var req roleUpsertReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	r := models.Role{
		Name: strings.TrimSpace(req.Name),
		Code: strings.TrimSpace(req.Code),
	}
	if r.Name == "" || r.Code == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "name/code不能为空"})
		return
	}
	if err := rc.svc.Create(c.Request.Context(), &r); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": r})
}

func (rc *RoleController) Get(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil || id <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid role id"})
		return
	}

	r, err := rc.svc.Get(c.Request.Context(), uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": r})
}

// Update 更新角色。
func (rc *RoleController) Update(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil || id <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid role id"})
		return
	}

	var req roleUpsertReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	role, err := rc.svc.Get(c.Request.Context(), uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	role.Name = strings.TrimSpace(req.Name)
	role.Code = strings.TrimSpace(req.Code)
	if role.Name == "" || role.Code == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "name/code不能为空"})
		return
	}

	if err := rc.svc.Update(c.Request.Context(), role); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": role})
}

func (rc *RoleController) Delete(c *gin.Context) {
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

	if err := rc.svc.Delete(c.Request.Context(), uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"msg": "deleted"})
}

// BindPermissions 绑定角色权限（全量替换）。
func (rc *RoleController) BindPermissions(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil || id <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid role id"})
		return
	}

	var req struct {
		PermissionIDs []uint `json:"permission_ids"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := rc.rolePermSvc.BindPermissions(c.Request.Context(), uint(id), req.PermissionIDs); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"msg": "permissions bound"})
}

// GetPermissions 查询角色已绑定权限。
func (rc *RoleController) GetPermissions(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil || id <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid role id"})
		return
	}

	perms, err := rc.rolePermSvc.GetRolePermissions(c.Request.Context(), uint(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": perms})
}
