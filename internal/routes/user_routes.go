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

// UserRoutes 处理UserRoutes相关逻辑。
func UserRoutes(r *gin.Engine, db *gorm.DB, mgr *tokenjwt.Manager) {
	// 说明：用户接口统一走 “JWT认证 + 权限码校验” 双层防护。
	userCtrl := rbac.NewUserController(db)
	// 定义并初始化当前变量。
	userRoleSvc := rbacsvc.NewUserRoleService(rbacdao.NewUserRoleDao(db))

	// 定义并初始化当前变量。
	userGroup := r.Group("/api/users")
	// 调用userGroup.Use完成当前处理。
	userGroup.Use(middleware.NewJWTMiddleware(mgr))
	{
		// 调用userGroup.GET完成当前处理。
		userGroup.GET("/", middleware.PermissionRequired(constants.PermUserList, userRoleSvc), userCtrl.List)
		// 调用userGroup.POST完成当前处理。
		userGroup.POST("/", middleware.PermissionRequired(constants.PermUserCreate, userRoleSvc), userCtrl.Create)
		// 调用userGroup.PUT完成当前处理。
		userGroup.PUT("/:id", middleware.PermissionRequired(constants.PermUserUpdate, userRoleSvc), userCtrl.Update)
		// 调用userGroup.DELETE完成当前处理。
		userGroup.DELETE("/:id", middleware.PermissionRequired(constants.PermUserDelete, userRoleSvc), userCtrl.Delete)

		// 个人信息接口也纳入 RBAC 权限控制
		userGroup.GET("/profile", middleware.PermissionRequired(constants.PermUserProfile, userRoleSvc), userCtrl.Profile)
		// 调用userGroup.GET完成当前处理。
		userGroup.GET("/:id", middleware.PermissionRequired(constants.PermUserView, userRoleSvc), userCtrl.Get)
	}
}
