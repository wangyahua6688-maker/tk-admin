package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"go-admin-full/internal/models"
	tokenjwt "go-admin-full/internal/token/jwt"
	"go-admin-full/internal/utils"
)

// NewJWTMiddleware JWT 认证中间件。
// 核心职责：
// 1. 解析 Authorization 头并校验 access token；
// 2. 校验 token 是否被撤销（黑名单）；
// 3. 二次校验用户状态（禁用用户立即失效）；
// 4. 将 uid/claims/access_token 注入 gin context 供后续 RBAC 使用。
func NewJWTMiddleware(mgr *tokenjwt.Manager) gin.HandlerFunc {
	if mgr == nil {
		// 为了安全，避免 nil 引发 panic，返回不通过任何请求的中间件
		return func(c *gin.Context) {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"code": 500, "msg": "jwt manager is not initialized"})
		}
	}
	return func(c *gin.Context) {
		// 1) 提取并规范化 Bearer Token。
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

		// 2) 校验 access token 签名、过期时间、发行者、黑名单状态。
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

		// 3) 把认证结果写入上下文，供后续 RBAC 与业务控制器复用。
		c.Set("uid", claims.UserID)
		c.Set("claims", claims)
		c.Set("access_token", tokenStr)
		c.Next()
	}
}
