package routes

import (
	"github.com/gin-gonic/gin"
	"go-admin-full/internal/constants"
	rbac "go-admin-full/internal/controllers/rbac"
	rbacdao "go-admin-full/internal/dao/rbac"
	"go-admin-full/internal/middleware"
	rbacsvc "go-admin-full/internal/services/rbac"
	tokenjwt "go-admin-full/internal/token/jwt"
	"gorm.io/gorm"
)

func RoleRoutes(r *gin.Engine, db *gorm.DB, mgr *tokenjwt.Manager) {
	rc := rbac.NewRoleController(db)
	userRoleSvc := rbacsvc.NewUserRoleService(rbacdao.NewUserRoleDao(db))

	// 角色管理：包含角色 CRUD + 角色权限绑定。
	rg := r.Group("/api/roles")
	rg.Use(middleware.NewJWTMiddleware(mgr))
	{
		rg.GET("/", middleware.PermissionRequired(constants.PermRoleList, userRoleSvc), rc.List)
		rg.POST("/", middleware.PermissionRequired(constants.PermRoleCreate, userRoleSvc), rc.Create)
		rg.PUT("/:id", middleware.PermissionRequired(constants.PermRoleUpdate, userRoleSvc), rc.Update)
		rg.GET("/:id", middleware.PermissionRequired(constants.PermRoleView, userRoleSvc), rc.Get)
		rg.DELETE("/:id", middleware.PermissionRequired(constants.PermRoleDelete, userRoleSvc), rc.Delete)
		// 角色权限管理
		rg.GET("/:id/permissions", middleware.PermissionRequired(constants.PermRolePermissionView, userRoleSvc), rc.GetPermissions)
		rg.PUT("/:id/permissions", middleware.PermissionRequired(constants.PermRolePermissionBind, userRoleSvc), rc.BindPermissions)
	}
}
