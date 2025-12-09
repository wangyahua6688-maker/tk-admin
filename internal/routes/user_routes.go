package routes

import (
	"github.com/gin-gonic/gin"
	"go-admin-full/internal/controllers"
	"go-admin-full/internal/middleware"
	"go-admin-full/internal/tokenpkg"
	"gorm.io/gorm"
)

func UserRoutes(r *gin.Engine, mgr *tokenpkg.Manager, db *gorm.DB) {
	// 创建用户控制器
	userCtrl := controllers.NewUserController(db)

	// 用户相关路由组（需要认证）
	userGroup := r.Group("/api/users")
	// 使用JWT中间件
	userGroup.Use(middleware.NewJWTMiddleware(mgr))
	{
		//userGroup.GET("/", userCtrl.List)
		//userGroup.GET("/:id", userCtrl.Get)
		//userGroup.POST("/", userCtrl.Create)
		//userGroup.PUT("/:id", userCtrl.Update)
		//userGroup.DELETE("/:id", userCtrl.Delete)

		// 需要特定权限的路由
		userGroup.GET("/profile", userCtrl.Profile)
		//userGroup.PUT("/profile", userCtrl.UpdateProfile)
	}
}
