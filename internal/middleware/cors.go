package middleware

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

// CORSOptions 定义CORSOptions相关结构。
type CORSOptions struct {
	// 处理当前语句逻辑。
	AllowedOrigins []string
	// 处理当前语句逻辑。
	AllowCredentials bool
}

// CORSMiddleware CORS中间件（白名单模式）
func CORSMiddleware(opt CORSOptions) gin.HandlerFunc {
	// 定义并初始化当前变量。
	allowed := make(map[string]struct{}, len(opt.AllowedOrigins))
	// 循环处理当前数据集合。
	for _, origin := range opt.AllowedOrigins {
		// 定义并初始化当前变量。
		o := strings.TrimSpace(origin)
		// 判断条件并进入对应分支逻辑。
		if o != "" {
			// 更新当前变量或字段值。
			allowed[o] = struct{}{}
		}
	}

	// 返回当前处理结果。
	return func(c *gin.Context) {
		// 安全响应头（基础防护）
		c.Header("X-Content-Type-Options", "nosniff")
		// 调用c.Header完成当前处理。
		c.Header("X-Frame-Options", "DENY")
		// 调用c.Header完成当前处理。
		c.Header("Referrer-Policy", "no-referrer")
		// 更新当前变量或字段值。
		c.Header("X-XSS-Protection", "1; mode=block")

		// 定义并初始化当前变量。
		origin := c.GetHeader("Origin")
		// 判断条件并进入对应分支逻辑。
		if origin != "" {
			// 判断条件并进入对应分支逻辑。
			if _, ok := allowed[origin]; ok {
				// 调用c.Writer.Header完成当前处理。
				c.Writer.Header().Set("Access-Control-Allow-Origin", origin)
				// 调用c.Writer.Header完成当前处理。
				c.Writer.Header().Set("Vary", "Origin")
				// 判断条件并进入对应分支逻辑。
				if opt.AllowCredentials {
					// 调用c.Writer.Header完成当前处理。
					c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
				}
				// 调用c.Writer.Header完成当前处理。
				c.Writer.Header().Set("Access-Control-Allow-Headers", "Authorization, Content-Type, X-Requested-With, X-Device-ID")
				// 调用c.Writer.Header完成当前处理。
				c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, PATCH, OPTIONS")
				// 进入新的代码块进行处理。
			} else if c.Request.Method == http.MethodOptions {
				// 调用c.AbortWithStatus完成当前处理。
				c.AbortWithStatus(http.StatusForbidden)
				// 返回当前处理结果。
				return
			}
		}

		// 判断条件并进入对应分支逻辑。
		if c.Request.Method == "OPTIONS" {
			// 调用c.AbortWithStatus完成当前处理。
			c.AbortWithStatus(204)
			// 返回当前处理结果。
			return
		}

		// 调用c.Next完成当前处理。
		c.Next()
	}
}
