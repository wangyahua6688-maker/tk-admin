package controllers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"go-admin-full/internal/dao"
	"go-admin-full/internal/services"
	"gorm.io/gorm"
)

// UserRoleController handles binding roles to users.
type UserRoleController struct {
	svc *services.UserRoleService
}

// NewUserRoleController constructs controller by wiring DAO -> Service
func NewUserRoleController(db *gorm.DB) *UserRoleController {
	d := dao.NewUserRoleDao(db)
	s := services.NewUserRoleService(d)
	return &UserRoleController{svc: s}
}

type bindRolesReq struct {
	UserID  uint   `json:"user_id" binding:"required"`
	RoleIDs []uint `json:"role_ids" binding:"required"`
}

func (uc *UserRoleController) BindRoles(c *gin.Context) {
	var req bindRolesReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := uc.svc.BindRoles(c.Request.Context(), req.UserID, req.RoleIDs); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"msg": "roles bound"})
}

func (uc *UserRoleController) AddRoles(c *gin.Context) {
	var req bindRolesReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := uc.svc.AddRoles(c.Request.Context(), req.UserID, req.RoleIDs); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"msg": "roles added"})
}

func (uc *UserRoleController) RemoveRoles(c *gin.Context) {
	var req bindRolesReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := uc.svc.RemoveRoles(c.Request.Context(), req.UserID, req.RoleIDs); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"msg": "roles removed"})
}

func (uc *UserRoleController) GetUserRoles(c *gin.Context) {
	// accept either /:id or query param id
	idStr := c.Param("id")
	if idStr == "" {
		idStr = c.Query("user_id")
	}
	if idStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "user id required"})
		return
	}
	id64, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user id"})
		return
	}
	roles, err := uc.svc.GetUserRoles(c.Request.Context(), uint(id64))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": roles})
}
