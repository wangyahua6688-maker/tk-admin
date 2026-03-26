package routes

import (
	"github.com/gin-gonic/gin"
	"go-admin-full/config"
	rbac "go-admin-full/internal/controllers/rbac"
	"go-admin-full/internal/middleware"
	tokenjwt "go-admin-full/internal/token/jwt"
	"gorm.io/gorm"
)

// AuthRoutes 注册认证相关路由。
// 分组原则：
// - /auth/login /auth/refresh（公共接口）
// - /auth/register（按配置决定是否开放）
// - /auth/logout（必须登录后访问）
func AuthRoutes(r *gin.Engine, db *gorm.DB, mgr *tokenjwt.Manager, cfg config.Config) {
	// 创建认证控制器，传递数据库连接
	auth := rbac.NewAuthController(db, mgr, cfg)

	// 公共路由组：不走 JWT 中间件。
	public := r.Group("/auth")
	{
		// 调用public.POST完成当前处理。
		public.POST("/login", auth.Login)
		// 判断条件并进入对应分支逻辑。
		if cfg.Auth.AllowPublicRegister {
			// 调用public.POST完成当前处理。
			public.POST("/register", auth.Register)
		}
		// 调用public.POST完成当前处理。
		public.POST("/refresh", auth.Refresh)
	}

	// 受保护路由组：统一挂载 JWT 中间件。
	protected := r.Group("/auth")
	// 调用protected.Use完成当前处理。
	protected.Use(middleware.NewJWTMiddleware(mgr))
	{
		// 调用protected.POST完成当前处理。
		protected.POST("/logout", auth.Logout)
	}
}
