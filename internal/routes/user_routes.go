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

func UserRoutes(r *gin.Engine, db *gorm.DB, mgr *tokenpkg.Manager) {
	// 创建用户控制器
	userCtrl := controllers.NewUserController(db)
	userRoleSvc := services.NewUserRoleService(dao.NewUserRoleDao(db))

	// 用户相关路由组（需要认证）
	userGroup := r.Group("/api/users")
	// 使用JWT中间件
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
