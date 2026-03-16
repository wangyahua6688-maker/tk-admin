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

func UserRoutes(r *gin.Engine, db *gorm.DB, mgr *tokenjwt.Manager) {
	// 说明：用户接口统一走 “JWT认证 + 权限码校验” 双层防护。
	userCtrl := rbac.NewUserController(db)
	userRoleSvc := rbacsvc.NewUserRoleService(rbacdao.NewUserRoleDao(db))

	userGroup := r.Group("/api/users")
	userGroup.Use(middleware.NewJWTMiddleware(mgr))
	{
		userGroup.GET("/", middleware.PermissionRequired(constants.PermUserList, userRoleSvc), userCtrl.List)
		userGroup.POST("/", middleware.PermissionRequired(constants.PermUserCreate, userRoleSvc), userCtrl.Create)
		userGroup.PUT("/:id", middleware.PermissionRequired(constants.PermUserUpdate, userRoleSvc), userCtrl.Update)
		userGroup.DELETE("/:id", middleware.PermissionRequired(constants.PermUserDelete, userRoleSvc), userCtrl.Delete)

		// 个人信息接口也纳入 RBAC 权限控制
		userGroup.GET("/profile", middleware.PermissionRequired(constants.PermUserProfile, userRoleSvc), userCtrl.Profile)
		userGroup.GET("/:id", middleware.PermissionRequired(constants.PermUserView, userRoleSvc), userCtrl.Get)
	}
}
