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

func UserRoleRoutes(r *gin.Engine, db *gorm.DB, mgr *tokenjwt.Manager) {
	// 用户角色管理：绑定/新增/移除/查询。
	userRoleCtrl := rbac.NewUserRoleController(db)
	userRoleSvc := rbacsvc.NewUserRoleService(rbacdao.NewUserRoleDao(db))

	userRoleGroup := r.Group("/api/users/role")
	userRoleGroup.Use(middleware.NewJWTMiddleware(mgr))
	{
		userRoleGroup.POST("/bind", middleware.PermissionRequired(constants.PermUserRoleBind, userRoleSvc), userRoleCtrl.BindRoles)
		userRoleGroup.POST("/add", middleware.PermissionRequired(constants.PermUserRoleAdd, userRoleSvc), userRoleCtrl.AddRoles)
		userRoleGroup.POST("/remove", middleware.PermissionRequired(constants.PermUserRoleRemove, userRoleSvc), userRoleCtrl.RemoveRoles)
		userRoleGroup.GET("/:id", middleware.PermissionRequired(constants.PermUserRoleView, userRoleSvc), userRoleCtrl.GetUserRoles)
	}
}
