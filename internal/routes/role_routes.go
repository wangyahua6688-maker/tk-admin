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

// RoleRoutes 处理RoleRoutes相关逻辑。
func RoleRoutes(r *gin.Engine, db *gorm.DB, mgr *tokenjwt.Manager) {
	// 注入 mgr 使 RoleController 能在权限变更后清除 Redis 权限缓存
	rc := rbac.NewRoleController(db, mgr)
	// 定义并初始化当前变量。
	userRoleSvc := rbacsvc.NewUserRoleService(rbacdao.NewUserRoleDao(db))

	// 角色管理：包含角色 CRUD + 角色权限绑定。
	rg := r.Group("/api/roles")
	// 调用rg.Use完成当前处理。
	rg.Use(middleware.NewJWTMiddleware(mgr))
	{
		// 调用rg.GET完成当前处理。
		rg.GET("/", middleware.PermissionRequired(constants.PermRoleList, userRoleSvc, mgr), rc.List)
		// 调用rg.POST完成当前处理。
		rg.POST("/", middleware.PermissionRequired(constants.PermRoleCreate, userRoleSvc, mgr), rc.Create)
		// 调用rg.PUT完成当前处理。
		rg.PUT("/:id", middleware.PermissionRequired(constants.PermRoleUpdate, userRoleSvc, mgr), rc.Update)
		// 调用rg.GET完成当前处理。
		rg.GET("/:id", middleware.PermissionRequired(constants.PermRoleView, userRoleSvc, mgr), rc.Get)
		// 调用rg.DELETE完成当前处理。
		rg.DELETE("/:id", middleware.PermissionRequired(constants.PermRoleDelete, userRoleSvc, mgr), rc.Delete)
		// 角色权限管理
		rg.GET("/:id/permissions", middleware.PermissionRequired(constants.PermRolePermissionView, userRoleSvc, mgr), rc.GetPermissions)
		// 调用rg.PUT完成当前处理。
		rg.PUT("/:id/permissions", middleware.PermissionRequired(constants.PermRolePermissionBind, userRoleSvc, mgr), rc.BindPermissions)
	}
}
