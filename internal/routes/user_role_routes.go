package routes

import (
	"github.com/gin-gonic/gin"
	"go-admin-full/internal/constants"
	"go-admin-full/internal/controllers"
	"go-admin-full/internal/dao"
	"go-admin-full/internal/middleware"
	"go-admin-full/internal/services"
	tokenjwt "go-admin-full/internal/token/jwt"
	"gorm.io/gorm"
)

func UserRoleRoutes(r *gin.Engine, db *gorm.DB, mgr *tokenjwt.Manager) {
	// 用户角色管理：绑定/新增/移除/查询。
	userRoleCtrl := controllers.NewUserRoleController(db)
	userRoleSvc := services.NewUserRoleService(dao.NewUserRoleDao(db))

	userRoleGroup := r.Group("/api/users/role")
	userRoleGroup.Use(middleware.NewJWTMiddleware(mgr))
	{
		userRoleGroup.POST("/bind", middleware.PermissionRequired(constants.PermUserRoleBind, userRoleSvc), userRoleCtrl.BindRoles)
		userRoleGroup.POST("/add", middleware.PermissionRequired(constants.PermUserRoleAdd, userRoleSvc), userRoleCtrl.AddRoles)
		userRoleGroup.POST("/remove", middleware.PermissionRequired(constants.PermUserRoleRemove, userRoleSvc), userRoleCtrl.RemoveRoles)
		userRoleGroup.GET("/:id", middleware.PermissionRequired(constants.PermUserRoleView, userRoleSvc), userRoleCtrl.GetUserRoles)
	}
}
