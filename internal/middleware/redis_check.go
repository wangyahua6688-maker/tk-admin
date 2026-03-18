package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	tokenjwt "go-admin-full/internal/token/jwt"
)

// RedisCheckMiddleware Redis连接检查中间件
func RedisCheckMiddleware(mgr *tokenjwt.Manager) gin.HandlerFunc {
	// 返回当前处理结果。
	return func(c *gin.Context) {
		// 对于重要的API端点，检查Redis连接
		if c.Request.URL.Path == "/auth/refresh" || c.Request.URL.Path == "/auth/logout" {
			// 判断条件并进入对应分支逻辑。
			if mgr.Store != nil {
				// 判断条件并进入对应分支逻辑。
				if err := mgr.Store.Ping(); err != nil {
					// 调用c.JSON完成当前处理。
					c.JSON(http.StatusServiceUnavailable, gin.H{
						// 处理当前语句逻辑。
						"code": 503,
						// 处理当前语句逻辑。
						"msg": "缓存服务不可用",
					})
					// 调用c.Abort完成当前处理。
					c.Abort()
					// 返回当前处理结果。
					return
				}
			}
		}

		// 调用c.Next完成当前处理。
		c.Next()
	}
}
