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

// UserRoleController handles binding roles to users.
type UserRoleController struct {
	svc    *services.UserRoleService
	msgSvc *services.SystemMessageService
}

// NewUserRoleController constructs controller by wiring DAO -> Service
func NewUserRoleController(db *gorm.DB) *UserRoleController {
	d := dao.NewUserRoleDao(db)
	s := services.NewUserRoleService(d)
	msgSvc := services.NewSystemMessageService(dao.NewSystemMessageDao(db))
	return &UserRoleController{svc: s, msgSvc: msgSvc}
}

type bindRolesReq struct {
	UserID  uint   `json:"user_id" binding:"required"`
	RoleIDs []uint `json:"role_ids"`
}

func (uc *UserRoleController) BindRoles(c *gin.Context) {
	var req bindRolesReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	// 安全防护：禁止修改当前登录用户自身角色，避免误操作导致权限锁死或提权。
	if req.UserID == c.GetUint("uid") {
		c.JSON(http.StatusForbidden, gin.H{"error": "不可修改当前登录用户角色"})
		return
	}

	if err := uc.svc.BindRoles(c.Request.Context(), req.UserID, req.RoleIDs); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	roles, _ := uc.svc.GetUserRoles(c.Request.Context(), req.UserID)
	_ = uc.msgSvc.PushToUser(
		c.Request.Context(),
		req.UserID,
		"角色分配变更通知",
		"管理员已调整你的角色分配，当前角色："+joinRoleNames(roles),
		"warning",
		"user_role",
		req.UserID,
		c.GetUint("uid"),
	)

	c.JSON(http.StatusOK, gin.H{"msg": "roles bound"})
}

func (uc *UserRoleController) AddRoles(c *gin.Context) {
	var req bindRolesReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if len(req.RoleIDs) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "role_ids不能为空"})
		return
	}
	if req.UserID == c.GetUint("uid") {
		c.JSON(http.StatusForbidden, gin.H{"error": "不可修改当前登录用户角色"})
		return
	}

	if err := uc.svc.AddRoles(c.Request.Context(), req.UserID, req.RoleIDs); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	roles, _ := uc.svc.GetUserRoles(c.Request.Context(), req.UserID)
	_ = uc.msgSvc.PushToUser(
		c.Request.Context(),
		req.UserID,
		"角色新增通知",
		"管理员已给你新增角色，当前角色："+joinRoleNames(roles),
		"info",
		"user_role",
		req.UserID,
		c.GetUint("uid"),
	)

	c.JSON(http.StatusOK, gin.H{"msg": "roles added"})
}

func (uc *UserRoleController) RemoveRoles(c *gin.Context) {
	var req bindRolesReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if len(req.RoleIDs) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "role_ids不能为空"})
		return
	}
	if req.UserID == c.GetUint("uid") {
		c.JSON(http.StatusForbidden, gin.H{"error": "不可修改当前登录用户角色"})
		return
	}

	if err := uc.svc.RemoveRoles(c.Request.Context(), req.UserID, req.RoleIDs); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	roles, _ := uc.svc.GetUserRoles(c.Request.Context(), req.UserID)
	_ = uc.msgSvc.PushToUser(
		c.Request.Context(),
		req.UserID,
		"角色移除通知",
		"管理员已移除你的部分角色，当前角色："+joinRoleNames(roles),
		"warning",
		"user_role",
		req.UserID,
		c.GetUint("uid"),
	)

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

func joinRoleNames(roles []models.Role) string {
	if len(roles) == 0 {
		return "无"
	}
	names := make([]string, 0, len(roles))
	for _, role := range roles {
		if strings.TrimSpace(role.Name) == "" {
			continue
		}
		names = append(names, role.Name)
	}
	if len(names) == 0 {
		return "无"
	}
	return strings.Join(names, "、")
}
