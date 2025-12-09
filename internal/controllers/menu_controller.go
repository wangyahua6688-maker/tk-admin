package controllers

import (
	"go-admin-full/internal/dao"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"go-admin-full/internal/models"
	"go-admin-full/internal/services"
	"gorm.io/gorm"
)

type MenuController struct {
	svc      *services.MenuService
	svcRoles *services.UserRoleService
}

func NewMenuController(db *gorm.DB) *MenuController {
	menuDao := dao.NewMenuDao(db)
	menuSvc := services.NewMenuService(menuDao)

	userRoleDao := dao.NewUserRoleDao(db)
	userRoleSvc := services.NewUserRoleService(userRoleDao)

	return &MenuController{
		svc:      menuSvc,
		svcRoles: userRoleSvc,
	}
}
func (mc *MenuController) List(c *gin.Context) {
	items, err := mc.svc.List(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": items})
}

func (mc *MenuController) Create(c *gin.Context) {
	var m models.Menu
	if err := c.ShouldBindJSON(&m); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := mc.svc.Create(c.Request.Context(), &m); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": m})
}

func (mc *MenuController) Get(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	m, err := mc.svc.Get(c.Request.Context(), uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": m})
}

func (mc *MenuController) Delete(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	if err := mc.svc.Delete(c.Request.Context(), uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"msg": "deleted"})
}

func (mc *MenuController) FrontendRoutes(c *gin.Context) {
	uid := c.GetUint("uid")

	roles, err := mc.svcRoles.GetUserRoles(c.Request.Context(), uid)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	var roleCodes []string
	for _, r := range roles {
		roleCodes = append(roleCodes, r.Code)
	}

	tree, err := mc.svc.ListForRoles(c.Request.Context(), roleCodes)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"data": tree})
}
