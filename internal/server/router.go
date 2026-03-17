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

func SetupRouter(mgr *tokenjwt.Manager, db *gorm.DB, redisClient *redis.Client, logger *utils.Logger, cfg config.Config) *gin.Engine {
	r := gin.Default()

	// 全局中间件链（顺序很重要）：
	// 1) 访问日志 -> 2) panic恢复 -> 3) CORS安全策略 -> 4) DB/Redis上下文注入。
	// 说明：JWT/RBAC 中间件在具体业务路由中单独挂载，不在这里全局挂载。
	r.Use(middleware.JSONLoggerMiddleware())

	r.Use(middleware.RecoveryMiddleware())

	r.Use(middleware.CORSMiddleware(middleware.CORSOptions{
		AllowedOrigins:   cfg.CORS.AllowedOrigins,
		AllowCredentials: cfg.CORS.AllowCredentials,
	}))

	// 注入数据库连接到 request context，便于 DAO 在无侵入情况下获取 db。
	r.Use(func(c *gin.Context) {
		c.Request = c.Request.WithContext(context.WithValue(c.Request.Context(), "db", db))
		c.Next()
	})

	// 若 Redis 可用，注入 Redis 客户端并追加健康检查。
	if redisClient != nil {
		r.Use(func(c *gin.Context) {
			c.Request = c.Request.WithContext(context.WithValue(c.Request.Context(), "redis", redisClient))
			c.Next()
		})

		r.Use(middleware.RedisCheckMiddleware(mgr))
	}

	// 路由注册分为三段：
	// 1) 认证路由（/auth）
	// 2) 用户路由（/api/users）
	// 3) RBAC 聚合路由（角色/权限/菜单/用户角色/审计）
	// 4) 上传路由（/api/upload）
	routes.AuthRoutes(r, db, mgr, cfg.Auth.AllowPublicRegister)
	routes.UserRoutes(r, db, mgr)
	routes.RBACRoutes(r, db, mgr)
	routes.UploadRoutes(r, mgr)

	// 静态资源路由（用于图片回显）
	// 注意：如果 save_path 是前端目录，前端开发服务器通常会自动处理静态文件
	// 但为了后端独立可用，这里也挂载一下。
	if cfg.Upload.BaseURL != "" && cfg.Upload.SavePath != "" {
		r.Static(cfg.Upload.BaseURL, cfg.Upload.SavePath)
	}

	// 统一健康检查：输出数据库、Redis、token 存储状态。
	r.GET("/health", func(c *gin.Context) {
		healthStatus := gin.H{
			"status":  "ok",
			"service": "go-admin-full",
			"time":    time.Now().Format(time.RFC3339),
		}

		// 检查数据库连接
		sqlDB, err := db.DB()
		if err != nil {
			healthStatus["database"] = "error"
			healthStatus["status"] = "degraded"
		} else if err := sqlDB.Ping(); err != nil {
			healthStatus["database"] = "unhealthy"
			healthStatus["status"] = "degraded"
		} else {
			healthStatus["database"] = "healthy"
		}

		// 检查Redis连接
		if redisClient != nil {
			if err := redisClient.Ping(c.Request.Context()).Err(); err != nil {
				healthStatus["redis"] = "unhealthy"
				healthStatus["status"] = "degraded"
			} else {
				healthStatus["redis"] = "healthy"
			}
		}

		// 检查token存储
		if mgr.Store != nil {
			if err := mgr.Store.Ping(); err != nil {
				healthStatus["token_store"] = "unhealthy"
				healthStatus["status"] = "degraded"
			} else {
				healthStatus["token_store"] = "healthy"
			}
		}

		c.JSON(200, healthStatus)
	})

	// 统一 404 输出结构，便于前端和日志系统识别。
	r.NoRoute(func(c *gin.Context) {
		logger.Warn("404 Not Found: %s %s", c.Request.Method, c.Request.URL.Path)
		c.JSON(404, gin.H{
			"code": 404,
			"msg":  "请求的资源不存在",
		})
	})

	return r
}
