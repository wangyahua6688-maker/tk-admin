package routes

import (
	"github.com/gin-gonic/gin"
	"go-admin/internal/constants"
	rbac "go-admin/internal/controllers/rbac"
	rbacdao "go-admin/internal/dao/rbac"
	"go-admin/internal/middleware"
	rbacsvc "go-admin/internal/services/rbac"
	tokenjwt "go-admin/internal/token/jwt"
	"gorm.io/gorm"
)

// PermissionRoutes 处理PermissionRoutes相关逻辑。
func PermissionRoutes(r *gin.Engine, db *gorm.DB, mgr *tokenjwt.Manager) {
	// 定义并初始化当前变量。
	pc := rbac.NewPermissionController(db)
	// 定义并初始化当前变量。
	userRoleSvc := rbacsvc.NewUserRoleService(rbacdao.NewUserRoleDao(db))

	// 权限管理：提供权限点的 CRUD。
	pr := r.Group("/api/permissions")
	// 调用pr.Use完成当前处理。
	pr.Use(middleware.NewJWTMiddleware(mgr))
	{
		// 调用pr.GET完成当前处理。
		pr.GET("/", middleware.PermissionRequired(constants.PermPermissionList, userRoleSvc, mgr), pc.List)
		// 调用pr.POST完成当前处理。
		pr.POST("/", middleware.PermissionRequired(constants.PermPermissionCreate, userRoleSvc, mgr), pc.Create)
		// 调用pr.PUT完成当前处理。
		pr.PUT("/:id", middleware.PermissionRequired(constants.PermPermissionUpdate, userRoleSvc, mgr), pc.Update)
		// 调用pr.GET完成当前处理。
		pr.GET("/:id", middleware.PermissionRequired(constants.PermPermissionView, userRoleSvc, mgr), pc.Get)
		// 调用pr.DELETE完成当前处理。
		pr.DELETE("/:id", middleware.PermissionRequired(constants.PermPermissionDelete, userRoleSvc, mgr), pc.Delete)
	}
}
