package routes

import (
	"github.com/gin-gonic/gin"
	"go-admin-full/internal/constants"
	"go-admin-full/internal/controllers"
	"go-admin-full/internal/dao"
	"go-admin-full/internal/middleware"
	"go-admin-full/internal/services"
	"go-admin-full/internal/tokenpkg"
	"gorm.io/gorm"
)

func UserRoleRoutes(r *gin.Engine, db *gorm.DB, mgr *tokenpkg.Manager) {
	// 创建用户控制器
	userRoleCtrl := controllers.NewUserRoleController(db)
	userRoleSvc := services.NewUserRoleService(dao.NewUserRoleDao(db))

	// 用户相关路由组（需要认证）
	userRoleGroup := r.Group("/api/users/role")
	// 使用JWT中间件
	userRoleGroup.Use(middleware.NewJWTMiddleware(mgr))
	{
		userRoleGroup.POST("/bind", middleware.PermissionRequired(constants.PermUserRoleBind, userRoleSvc), userRoleCtrl.BindRoles)
		userRoleGroup.POST("/add", middleware.PermissionRequired(constants.PermUserRoleAdd, userRoleSvc), userRoleCtrl.AddRoles)
		userRoleGroup.POST("/remove", middleware.PermissionRequired(constants.PermUserRoleRemove, userRoleSvc), userRoleCtrl.RemoveRoles)
		userRoleGroup.GET("/:id", middleware.PermissionRequired(constants.PermUserRoleView, userRoleSvc), userRoleCtrl.GetUserRoles)
	}
}
