package routes

import (
	"github.com/gin-gonic/gin"
	"go-admin-full/internal/controllers"
	"go-admin-full/internal/middleware"
	tokenjwt "go-admin-full/internal/token/jwt"
	"gorm.io/gorm"
)

// AuthRoutes 注册认证相关路由。
// 分组原则：
// - /auth/login /auth/refresh（公共接口）
// - /auth/register（按配置决定是否开放）
// - /auth/logout（必须登录后访问）
func AuthRoutes(r *gin.Engine, db *gorm.DB, mgr *tokenjwt.Manager, allowPublicRegister bool) {
	// 创建认证控制器，传递数据库连接
	auth := controllers.NewAuthController(db, mgr)

	// 公共路由组：不走 JWT 中间件。
	public := r.Group("/auth")
	{
		public.POST("/login", auth.Login)
		if allowPublicRegister {
			public.POST("/register", auth.Register)
		}
		public.POST("/refresh", auth.Refresh)
	}

	// 受保护路由组：统一挂载 JWT 中间件。
	protected := r.Group("/auth")
	protected.Use(middleware.NewJWTMiddleware(mgr))
	{
		protected.POST("/logout", auth.Logout)
	}
}
