package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	tokenjwt "go-admin-full/internal/token/jwt"
)

// RedisCheckMiddleware Redis连接检查中间件
func RedisCheckMiddleware(mgr *tokenjwt.Manager) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 对于重要的API端点，检查Redis连接
		if c.Request.URL.Path == "/auth/refresh" || c.Request.URL.Path == "/auth/logout" {
			if mgr.Store != nil {
				if err := mgr.Store.Ping(); err != nil {
					c.JSON(http.StatusServiceUnavailable, gin.H{
						"code": 503,
						"msg":  "缓存服务不可用",
					})
					c.Abort()
					return
				}
			}
		}

		c.Next()
	}
}
