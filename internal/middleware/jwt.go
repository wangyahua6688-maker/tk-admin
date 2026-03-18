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
	// 判断条件并进入对应分支逻辑。
	if mgr == nil {
		// 为了安全，避免 nil 引发 panic，返回不通过任何请求的中间件
		return func(c *gin.Context) {
			// 调用c.AbortWithStatusJSON完成当前处理。
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"code": 500, "msg": "jwt manager is not initialized"})
		}
	}
	// 返回当前处理结果。
	return func(c *gin.Context) {
		// 1) 提取并规范化 Bearer Token。
		auth := c.GetHeader("Authorization")
		// 判断条件并进入对应分支逻辑。
		if auth == "" {
			// 调用c.AbortWithStatusJSON完成当前处理。
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"code": 401, "msg": "authorization header required"})
			// 返回当前处理结果。
			return
		}
		// 定义并初始化当前变量。
		tokenStr := strings.TrimSpace(auth)
		// 判断条件并进入对应分支逻辑。
		if strings.HasPrefix(strings.ToLower(tokenStr), "bearer ") {
			// 更新当前变量或字段值。
			tokenStr = tokenStr[7:]
		}
		// 判断条件并进入对应分支逻辑。
		if tokenStr == "" {
			// 调用c.AbortWithStatusJSON完成当前处理。
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"code": 401, "msg": "token is empty"})
			// 返回当前处理结果。
			return
		}

		// 2) 校验 access token 签名、过期时间、发行者、黑名单状态。
		claims, err := mgr.ValidateAccessToken(tokenStr)
		// 判断条件并进入对应分支逻辑。
		if err != nil {
			// 调用c.AbortWithStatusJSON完成当前处理。
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"code": 401, "msg": "invalid token"})
			// 返回当前处理结果。
			return
		}

		// 安全防护：每次请求校验用户状态，禁用用户立即失效。
		if db := utils.DBFromContext(c.Request.Context()); db != nil {
			// 声明当前变量。
			var user models.User
			// 判断条件并进入对应分支逻辑。
			if err := db.Select("id", "status").First(&user, claims.UserID).Error; err != nil || user.Status != 1 {
				// 调用c.AbortWithStatusJSON完成当前处理。
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"code": 401, "msg": "user disabled or not found"})
				// 返回当前处理结果。
				return
			}
		}

		// 3) 把认证结果写入上下文，供后续 RBAC 与业务控制器复用。
		c.Set("uid", claims.UserID)
		// 调用c.Set完成当前处理。
		c.Set("claims", claims)
		// 调用c.Set完成当前处理。
		c.Set("access_token", tokenStr)
		// 调用c.Next完成当前处理。
		c.Next()
	}
}
