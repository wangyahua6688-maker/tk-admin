package controllers

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"go-admin-full/internal/dao"
	"go-admin-full/internal/models"
	"go-admin-full/internal/services"
	"gorm.io/gorm"
)

type PermissionController struct {
	svc    *services.PermissionService
	msgSvc *services.SystemMessageService
}

type permissionUpsertReq struct {
	Name   string `json:"name" binding:"required"`
	Code   string `json:"code" binding:"required"`
	Type   string `json:"type"`
	Method string `json:"method"`
	Path   string `json:"path"`
}

func NewPermissionController(db *gorm.DB) *PermissionController {
	return &PermissionController{
		svc:    services.NewPermissionService(dao.NewPermissionDAO(db)),
		msgSvc: services.NewSystemMessageService(dao.NewSystemMessageDao(db)),
	}
}

func (pc *PermissionController) List(c *gin.Context) {
	items, err := pc.svc.List(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": items})
}

func (pc *PermissionController) Create(c *gin.Context) {
	var req permissionUpsertReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	p := models.Permission{
		Name:   strings.TrimSpace(req.Name),
		Code:   strings.TrimSpace(req.Code),
		Type:   strings.TrimSpace(req.Type),
		Method: strings.ToUpper(strings.TrimSpace(req.Method)),
		Path:   strings.TrimSpace(req.Path),
	}
	if p.Name == "" || p.Code == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "name/code不能为空"})
		return
	}
	if err := pc.svc.Create(c.Request.Context(), &p); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": p})
}

func (pc *PermissionController) Get(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil || id <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid permission id"})
		return
	}

	p, err := pc.svc.Get(c.Request.Context(), uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": p})
}

// Update 更新权限。
func (pc *PermissionController) Update(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil || id <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid permission id"})
		return
	}

	var req permissionUpsertReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	p, err := pc.svc.Get(c.Request.Context(), uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	p.Name = strings.TrimSpace(req.Name)
	p.Code = strings.TrimSpace(req.Code)
	p.Type = strings.TrimSpace(req.Type)
	p.Method = strings.ToUpper(strings.TrimSpace(req.Method))
	p.Path = strings.TrimSpace(req.Path)
	if p.Name == "" || p.Code == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "name/code不能为空"})
		return
	}

	if err := pc.svc.Update(c.Request.Context(), p); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

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

	c.JSON(http.StatusOK, gin.H{"data": p})
}

func (pc *PermissionController) Delete(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil || id <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid permission id"})
		return
	}

	permission, err := pc.svc.Get(c.Request.Context(), uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	affectedUserIDs, _ := pc.msgSvc.ListUserIDsByPermissionIDs(c.Request.Context(), []uint{uint(id)})

	if err := pc.svc.Delete(c.Request.Context(), uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

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

	c.JSON(http.StatusOK, gin.H{"msg": "deleted"})
}
