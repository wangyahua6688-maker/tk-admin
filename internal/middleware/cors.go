package middleware

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

type CORSOptions struct {
	AllowedOrigins   []string
	AllowCredentials bool
}

// CORSMiddleware CORS中间件（白名单模式）
func CORSMiddleware(opt CORSOptions) gin.HandlerFunc {
	allowed := make(map[string]struct{}, len(opt.AllowedOrigins))
	for _, origin := range opt.AllowedOrigins {
		o := strings.TrimSpace(origin)
		if o != "" {
			allowed[o] = struct{}{}
		}
	}

	return func(c *gin.Context) {
		// 安全响应头（基础防护）
		c.Header("X-Content-Type-Options", "nosniff")
		c.Header("X-Frame-Options", "DENY")
		c.Header("Referrer-Policy", "no-referrer")
		c.Header("X-XSS-Protection", "1; mode=block")

		origin := c.GetHeader("Origin")
		if origin != "" {
			if _, ok := allowed[origin]; ok {
				c.Writer.Header().Set("Access-Control-Allow-Origin", origin)
				c.Writer.Header().Set("Vary", "Origin")
				if opt.AllowCredentials {
					c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
				}
				c.Writer.Header().Set("Access-Control-Allow-Headers", "Authorization, Content-Type, X-Requested-With, X-Device-ID")
				c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, PATCH, OPTIONS")
			} else if c.Request.Method == http.MethodOptions {
				c.AbortWithStatus(http.StatusForbidden)
				return
			}
		}

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}
