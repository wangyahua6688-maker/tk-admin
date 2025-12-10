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

type MenuController struct {
	svc         *services.MenuService
	menuPermSvc *services.MenuPermissionService
}

func NewMenuController(db *gorm.DB) *MenuController {
	menuDao := dao.NewMenuDAO(db)
	menuSvc := services.NewMenuService(menuDao)
	menuPermSvc := services.NewMenuPermissionService(dao.NewMenuPermissionDao(db))

	return &MenuController{
		svc:         menuSvc,
		menuPermSvc: menuPermSvc,
	}
}

type menuUpsertReq struct {
	Title     string `json:"title" binding:"required"`
	Path      string `json:"path" binding:"required"`
	Icon      string `json:"icon"`
	ParentID  uint   `json:"parent_id"`
	Component string `json:"component"`
	OrderNum  int    `json:"order_num"`
}

// List 菜单列表查询。
func (mc *MenuController) List(c *gin.Context) {
	items, err := mc.svc.List(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": items})
}

// Create 创建菜单。
func (mc *MenuController) Create(c *gin.Context) {
	var req menuUpsertReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	m := models.Menu{
		Title:     strings.TrimSpace(req.Title),
		Path:      strings.TrimSpace(req.Path),
		Icon:      strings.TrimSpace(req.Icon),
		ParentID:  req.ParentID,
		Component: strings.TrimSpace(req.Component),
		OrderNum:  req.OrderNum,
	}

	if m.Title == "" || m.Path == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "title/path不能为空"})
		return
	}
	if m.ParentID > 0 {
		if _, err := mc.svc.Get(c.Request.Context(), m.ParentID); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "parent_id无效"})
			return
		}
	}

	if err := mc.svc.Create(c.Request.Context(), &m); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": m})
}

// Get 查询菜单详情。
func (mc *MenuController) Get(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil || id <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid menu id"})
		return
	}

	m, err := mc.svc.Get(c.Request.Context(), uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": m})
}

// Update 更新菜单。
func (mc *MenuController) Update(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil || id <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid menu id"})
		return
	}

	var req menuUpsertReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	menu, err := mc.svc.Get(c.Request.Context(), uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	menu.Title = strings.TrimSpace(req.Title)
	menu.Path = strings.TrimSpace(req.Path)
	menu.Icon = strings.TrimSpace(req.Icon)
	menu.ParentID = req.ParentID
	menu.Component = strings.TrimSpace(req.Component)
	menu.OrderNum = req.OrderNum

	if menu.Title == "" || menu.Path == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "title/path不能为空"})
		return
	}
	if menu.ParentID == uint(id) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "parent_id不能为自身"})
		return
	}
	if menu.ParentID > 0 {
		if _, err := mc.svc.Get(c.Request.Context(), menu.ParentID); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "parent_id无效"})
			return
		}
	}

	if err := mc.svc.Update(c.Request.Context(), menu); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": menu})
}

// Delete 删除菜单。
func (mc *MenuController) Delete(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil || id <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid menu id"})
		return
	}

	if err := mc.svc.Delete(c.Request.Context(), uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"msg": "deleted"})
}

// FrontendTree 用于前端动态路由，返回当前用户权限范围内的菜单树。
func (mc *MenuController) FrontendTree(c *gin.Context) {
	uid := c.GetUint("uid")
	menus, err := mc.svc.ListForUser(c.Request.Context(), uid)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	tree := mc.svc.BuildMenuTreeFromList(menus)
	c.JSON(200, gin.H{"data": tree})
}

// GetPermissions 查询菜单已绑定权限。
func (mc *MenuController) GetPermissions(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil || id <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid menu id"})
		return
	}

	perms, err := mc.menuPermSvc.GetMenuPermissions(c.Request.Context(), uint(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": perms})
}

// BindPermissions 绑定菜单权限（全量替换）。
func (mc *MenuController) BindPermissions(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil || id <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid menu id"})
		return
	}

	var req struct {
		PermissionIDs []uint `json:"permission_ids"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := mc.menuPermSvc.BindPermissions(c.Request.Context(), uint(id), req.PermissionIDs); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"msg": "menu permissions bound"})
}
