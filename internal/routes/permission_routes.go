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

func PermissionRoutes(r *gin.Engine, db *gorm.DB, mgr *tokenjwt.Manager) {
	pc := rbac.NewPermissionController(db)
	userRoleSvc := rbacsvc.NewUserRoleService(rbacdao.NewUserRoleDao(db))

	// 权限管理：提供权限点的 CRUD。
	pr := r.Group("/api/permissions")
	pr.Use(middleware.NewJWTMiddleware(mgr))
	{
		pr.GET("/", middleware.PermissionRequired(constants.PermPermissionList, userRoleSvc), pc.List)
		pr.POST("/", middleware.PermissionRequired(constants.PermPermissionCreate, userRoleSvc), pc.Create)
		pr.PUT("/:id", middleware.PermissionRequired(constants.PermPermissionUpdate, userRoleSvc), pc.Update)
		pr.GET("/:id", middleware.PermissionRequired(constants.PermPermissionView, userRoleSvc), pc.Get)
		pr.DELETE("/:id", middleware.PermissionRequired(constants.PermPermissionDelete, userRoleSvc), pc.Delete)
	}
}
