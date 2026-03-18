package server

import (
	"context"
	"time"

	"go-admin-full/config"
	"go-admin-full/internal/middleware"
	"go-admin-full/internal/routes"
	tokenjwt "go-admin-full/internal/token/jwt"
	"go-admin-full/internal/utils"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"
)

// SetupRouter 初始化Router。
func SetupRouter(mgr *tokenjwt.Manager, db *gorm.DB, redisClient *redis.Client, logger *utils.Logger, cfg config.Config) *gin.Engine {
	// 定义并初始化当前变量。
	r := gin.Default()

	// 全局中间件链（顺序很重要）：
	// 1) 访问日志 -> 2) panic恢复 -> 3) CORS安全策略 -> 4) DB/Redis上下文注入。
	// 说明：JWT/RBAC 中间件在具体业务路由中单独挂载，不在这里全局挂载。
	r.Use(middleware.JSONLoggerMiddleware())

	// 调用r.Use完成当前处理。
	r.Use(middleware.RecoveryMiddleware())

	// 调用r.Use完成当前处理。
	r.Use(middleware.CORSMiddleware(middleware.CORSOptions{
		// 处理当前语句逻辑。
		AllowedOrigins: cfg.CORS.AllowedOrigins,
		// 处理当前语句逻辑。
		AllowCredentials: cfg.CORS.AllowCredentials,
	}))

	// 注入数据库连接到 request context，便于 DAO 在无侵入情况下获取 db。
	r.Use(func(c *gin.Context) {
		// 更新当前变量或字段值。
		c.Request = c.Request.WithContext(context.WithValue(c.Request.Context(), "db", db))
		// 调用c.Next完成当前处理。
		c.Next()
	})

	// 若 Redis 可用，注入 Redis 客户端并追加健康检查。
	if redisClient != nil {
		// 调用r.Use完成当前处理。
		r.Use(func(c *gin.Context) {
			// 更新当前变量或字段值。
			c.Request = c.Request.WithContext(context.WithValue(c.Request.Context(), "redis", redisClient))
			// 调用c.Next完成当前处理。
			c.Next()
		})

		// 调用r.Use完成当前处理。
		r.Use(middleware.RedisCheckMiddleware(mgr))
	}

	// 路由注册分为三段：
	// 1) 认证路由（/auth）
	// 2) 用户路由（/api/users）
	// 3) RBAC 聚合路由（角色/权限/菜单/用户角色/审计）
	// 4) 上传路由（/api/upload）
	routes.AuthRoutes(r, db, mgr, cfg.Auth.AllowPublicRegister)
	// 调用routes.UserRoutes完成当前处理。
	routes.UserRoutes(r, db, mgr)
	// 调用routes.RBACRoutes完成当前处理。
	routes.RBACRoutes(r, db, mgr)
	// 调用routes.UploadRoutes完成当前处理。
	routes.UploadRoutes(r, mgr)

	// 静态资源路由（用于图片回显）
	// 注意：如果 save_path 是前端目录，前端开发服务器通常会自动处理静态文件
	// 但为了后端独立可用，这里也挂载一下。
	if cfg.Upload.BaseURL != "" && cfg.Upload.SavePath != "" {
		// 调用r.Static完成当前处理。
		r.Static(cfg.Upload.BaseURL, cfg.Upload.SavePath)
	}

	// 统一健康检查：输出数据库、Redis、token 存储状态。
	r.GET("/health", func(c *gin.Context) {
		// 定义并初始化当前变量。
		healthStatus := gin.H{
			// 处理当前语句逻辑。
			"status": "ok",
			// 处理当前语句逻辑。
			"service": "go-admin-full",
			// 调用time.Now完成当前处理。
			"time": time.Now().Format(time.RFC3339),
		}

		// 检查数据库连接
		sqlDB, err := db.DB()
		// 判断条件并进入对应分支逻辑。
		if err != nil {
			// 更新当前变量或字段值。
			healthStatus["database"] = "error"
			// 更新当前变量或字段值。
			healthStatus["status"] = "degraded"
			// 定义并初始化当前变量。
		} else if err := sqlDB.Ping(); err != nil {
			// 更新当前变量或字段值。
			healthStatus["database"] = "unhealthy"
			// 更新当前变量或字段值。
			healthStatus["status"] = "degraded"
			// 进入新的代码块进行处理。
		} else {
			// 更新当前变量或字段值。
			healthStatus["database"] = "healthy"
		}

		// 检查Redis连接
		if redisClient != nil {
			// 判断条件并进入对应分支逻辑。
			if err := redisClient.Ping(c.Request.Context()).Err(); err != nil {
				// 更新当前变量或字段值。
				healthStatus["redis"] = "unhealthy"
				// 更新当前变量或字段值。
				healthStatus["status"] = "degraded"
				// 进入新的代码块进行处理。
			} else {
				// 更新当前变量或字段值。
				healthStatus["redis"] = "healthy"
			}
		}

		// 检查token存储
		if mgr.Store != nil {
			// 判断条件并进入对应分支逻辑。
			if err := mgr.Store.Ping(); err != nil {
				// 更新当前变量或字段值。
				healthStatus["token_store"] = "unhealthy"
				// 更新当前变量或字段值。
				healthStatus["status"] = "degraded"
				// 进入新的代码块进行处理。
			} else {
				// 更新当前变量或字段值。
				healthStatus["token_store"] = "healthy"
			}
		}

		// 调用c.JSON完成当前处理。
		c.JSON(200, healthStatus)
	})

	// 统一 404 输出结构，便于前端和日志系统识别。
	r.NoRoute(func(c *gin.Context) {
		// 调用logger.Warn完成当前处理。
		logger.Warn("404 Not Found: %s %s", c.Request.Method, c.Request.URL.Path)
		// 调用c.JSON完成当前处理。
		c.JSON(404, gin.H{
			// 处理当前语句逻辑。
			"code": 404,
			// 处理当前语句逻辑。
			"msg": "请求的资源不存在",
		})
	})

	// 返回当前处理结果。
	return r
}
