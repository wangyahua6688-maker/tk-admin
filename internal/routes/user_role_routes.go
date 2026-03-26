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

// UserRoleRoutes 处理UserRoleRoutes相关逻辑。
func UserRoleRoutes(r *gin.Engine, db *gorm.DB, mgr *tokenjwt.Manager) {
	// 注入 mgr 使控制器能在角色变更后清除 Redis 权限缓存
	userRoleCtrl := rbac.NewUserRoleController(db, mgr)
	// 定义并初始化当前变量。
	userRoleSvc := rbacsvc.NewUserRoleService(rbacdao.NewUserRoleDao(db))

	// 定义并初始化当前变量。
	userRoleGroup := r.Group("/api/users/role")
	// 调用userRoleGroup.Use完成当前处理。
	userRoleGroup.Use(middleware.NewJWTMiddleware(mgr))
	{
		// 调用userRoleGroup.POST完成当前处理。
		userRoleGroup.POST("/bind", middleware.PermissionRequired(constants.PermUserRoleBind, userRoleSvc, mgr), userRoleCtrl.BindRoles)
		// 调用userRoleGroup.POST完成当前处理。
		userRoleGroup.POST("/add", middleware.PermissionRequired(constants.PermUserRoleAdd, userRoleSvc, mgr), userRoleCtrl.AddRoles)
		// 调用userRoleGroup.POST完成当前处理。
		userRoleGroup.POST("/remove", middleware.PermissionRequired(constants.PermUserRoleRemove, userRoleSvc, mgr), userRoleCtrl.RemoveRoles)
		// 调用userRoleGroup.GET完成当前处理。
		userRoleGroup.GET("/:id", middleware.PermissionRequired(constants.PermUserRoleView, userRoleSvc, mgr), userRoleCtrl.GetUserRoles)
	}
}
