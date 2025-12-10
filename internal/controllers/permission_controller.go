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

type PermissionController struct {
	svc *services.PermissionService
}

type permissionUpsertReq struct {
	Name   string `json:"name" binding:"required"`
	Code   string `json:"code" binding:"required"`
	Type   string `json:"type"`
	Method string `json:"method"`
	Path   string `json:"path"`
}

func NewPermissionController(db *gorm.DB) *PermissionController {
	return &PermissionController{svc: services.NewPermissionService(dao.NewPermissionDAO(db))}
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
	c.JSON(http.StatusOK, gin.H{"data": p})
}

func (pc *PermissionController) Delete(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil || id <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid permission id"})
		return
	}

	if err := pc.svc.Delete(c.Request.Context(), uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"msg": "deleted"})
}
