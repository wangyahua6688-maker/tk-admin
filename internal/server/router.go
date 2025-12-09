package server

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"go-admin-full/internal/middleware"
	"go-admin-full/internal/routes"
	"go-admin-full/internal/tokenpkg"
	"go-admin-full/internal/utils"
	"gorm.io/gorm"
	"time"
)

func SetupRouter(mgr *tokenpkg.Manager, db *gorm.DB, redisClient *redis.Client, logger *utils.Logger) *gin.Engine {
	r := gin.Default()

	// 添加全局中间件
	// 1. 全局日志中间件（使用自定义logger）
	r.Use(middleware.JSONLoggerMiddleware())

	// 2. 全局错误恢复中间件
	r.Use(middleware.RecoveryMiddleware())

	// 3. 全局CORS中间件
	r.Use(middleware.CORSMiddleware())

	// 4. 数据库上下文中间件（将db注入到请求上下文中）
	r.Use(func(c *gin.Context) {
		c.Request = c.Request.WithContext(context.WithValue(c.Request.Context(), "db", db))
		c.Next()
	})

	// 5. 如果manager有Redis存储，添加Redis连接上下文
	if redisClient != nil {
		r.Use(func(c *gin.Context) {
			c.Request = c.Request.WithContext(context.WithValue(c.Request.Context(), "redis", redisClient))
			c.Next()
		})

		// 添加Redis连接健康检查中间件
		r.Use(middleware.RedisCheckMiddleware(mgr))
	}

	// 注册路由
	routes.AuthRoutes(r, db, mgr)
	routes.UserRoutes(r, db, mgr)
	routes.MenuRoutes(r, db, mgr)
	routes.PermissionRoutes(r, db, mgr)

	// 健康检查路由
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

	// 404处理
	r.NoRoute(func(c *gin.Context) {
		logger.Warn("404 Not Found: %s %s", c.Request.Method, c.Request.URL.Path)
		c.JSON(404, gin.H{
			"code": 404,
			"msg":  "请求的资源不存在",
		})
	})

	return r
}
