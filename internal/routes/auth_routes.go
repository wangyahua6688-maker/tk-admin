package routes

import (
	"github.com/gin-gonic/gin"
	"go-admin-full/internal/controllers"
	"go-admin-full/internal/middleware"
	"go-admin-full/internal/tokenpkg"
	"gorm.io/gorm"
)

// AuthRoutes 注册认证相关路由
func AuthRoutes(r *gin.Engine, db *gorm.DB, mgr *tokenpkg.Manager, allowPublicRegister bool) {
	// 创建认证控制器，传递数据库连接
	auth := controllers.NewAuthController(db, mgr)

	// 公共路由组
	public := r.Group("/auth")
	{
		public.POST("/login", auth.Login)
		if allowPublicRegister {
			public.POST("/register", auth.Register)
		}
		public.POST("/refresh", auth.Refresh)
	}

	// 需要认证的路由组
	protected := r.Group("/auth")
	protected.Use(middleware.NewJWTMiddleware(mgr))
	{
		protected.POST("/logout", auth.Logout)
	}
}
