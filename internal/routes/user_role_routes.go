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

// UserRoleRoutes 处理UserRoleRoutes相关逻辑。
func UserRoleRoutes(r *gin.Engine, db *gorm.DB, mgr *tokenjwt.Manager) {
	// 用户角色管理：绑定/新增/移除/查询。
	userRoleCtrl := rbac.NewUserRoleController(db)
	// 定义并初始化当前变量。
	userRoleSvc := rbacsvc.NewUserRoleService(rbacdao.NewUserRoleDao(db))

	// 定义并初始化当前变量。
	userRoleGroup := r.Group("/api/users/role")
	// 调用userRoleGroup.Use完成当前处理。
	userRoleGroup.Use(middleware.NewJWTMiddleware(mgr))
	{
		// 调用userRoleGroup.POST完成当前处理。
		userRoleGroup.POST("/bind", middleware.PermissionRequired(constants.PermUserRoleBind, userRoleSvc), userRoleCtrl.BindRoles)
		// 调用userRoleGroup.POST完成当前处理。
		userRoleGroup.POST("/add", middleware.PermissionRequired(constants.PermUserRoleAdd, userRoleSvc), userRoleCtrl.AddRoles)
		// 调用userRoleGroup.POST完成当前处理。
		userRoleGroup.POST("/remove", middleware.PermissionRequired(constants.PermUserRoleRemove, userRoleSvc), userRoleCtrl.RemoveRoles)
		// 调用userRoleGroup.GET完成当前处理。
		userRoleGroup.GET("/:id", middleware.PermissionRequired(constants.PermUserRoleView, userRoleSvc), userRoleCtrl.GetUserRoles)
	}
}
