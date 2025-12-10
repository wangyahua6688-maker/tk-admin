package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"go-admin-full/internal/models"
	"go-admin-full/internal/tokenpkg"
	"go-admin-full/internal/utils"
)

// NewJWTMiddleware 接收一个 *tokenpkg.Manager 并返回 gin.HandlerFunc
// 推荐在程序启动处构造 Manager 并注入到路由注册中： middleware.NewJWTMiddleware(manager)
func NewJWTMiddleware(mgr *tokenpkg.Manager) gin.HandlerFunc {
	if mgr == nil {
		// 为了安全，避免 nil 引发 panic，返回不通过任何请求的中间件
		return func(c *gin.Context) {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"code": 500, "msg": "jwt manager is not initialized"})
		}
	}
	return func(c *gin.Context) {
		auth := c.GetHeader("Authorization")
		if auth == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"code": 401, "msg": "authorization header required"})
			return
		}
		tokenStr := strings.TrimSpace(auth)
		if strings.HasPrefix(strings.ToLower(tokenStr), "bearer ") {
			tokenStr = tokenStr[7:]
		}
		if tokenStr == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"code": 401, "msg": "token is empty"})
			return
		}

		claims, err := mgr.ValidateAccessToken(tokenStr)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"code": 401, "msg": "invalid token"})
			return
		}

		// 安全防护：每次请求校验用户状态，禁用用户立即失效。
		if db := utils.DBFromContext(c.Request.Context()); db != nil {
			var user models.User
			if err := db.Select("id", "status").First(&user, claims.UserID).Error; err != nil || user.Status != 1 {
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"code": 401, "msg": "user disabled or not found"})
				return
			}
		}

		// 将解析到的 uid 写入 context（key: "uid"）
		c.Set("uid", claims.UserID)
		c.Set("claims", claims)
		c.Set("access_token", tokenStr)
		c.Next()
	}
}
