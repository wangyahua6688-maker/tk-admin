package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	gormx "github.com/wangyahua6688-maker/tk-common/utils/dbx/gormx"
	"go-admin/config"
	"go-admin/internal/auth/sessioncookie"
	"go-admin/internal/constants"
	"go-admin/internal/models"
	tokenjwt "go-admin/internal/token/jwt"
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
	cookieOpt := sessioncookie.FromConfig(config.GetConfig())
	// 返回当前处理结果。
	return func(c *gin.Context) {
		// 1) 先从 Authorization 头提取，再回退到 HttpOnly Cookie。
		tokenStr, fromCookie := resolveAccessToken(c, cookieOpt.AccessTokenName)
		// 判断条件并进入对应分支逻辑。
		if tokenStr == "" {
			// 调用c.AbortWithStatusJSON完成当前处理。
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"code": 401, "msg": "authentication required"})
			// 返回当前处理结果。
			return
		}
		// Cookie 鉴权场景增加轻量 CSRF 防护：
		// 非安全方法必须带 X-Device-ID（浏览器跨站表单无法伪造该头）。
		if fromCookie && !isSafeMethod(c.Request.Method) && strings.TrimSpace(c.GetHeader("X-Device-ID")) == "" {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"code": 403, "msg": "csrf validation failed"})
			return
		}

		// 2) 校验 access token 签名、过期时间、发行者、黑名单状态。
		claims, err := mgr.ValidateAccessToken(tokenStr)
		// 判断条件并进入对应分支逻辑。
		if err != nil {
			msg := "invalid token"
			if err == constants.ErrExpiredToken {
				msg = "token is expired"
			}
			// 调用c.AbortWithStatusJSON完成当前处理。
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"code": 401, "msg": msg})
			// 返回当前处理结果。
			return
		}

		// 安全防护：每次请求校验用户状态，禁用用户立即失效。
		if db := gormx.DBFromContext(c.Request.Context()); db != nil {
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

// resolveAccessToken 从 Authorization 头或认证 Cookie 提取 access token。
func resolveAccessToken(c *gin.Context, accessCookieName string) (string, bool) {
	auth := strings.TrimSpace(c.GetHeader("Authorization"))
	if strings.HasPrefix(strings.ToLower(auth), "bearer ") {
		return strings.TrimSpace(auth[7:]), false
	}
	if auth != "" {
		return auth, false
	}

	raw, err := c.Cookie(accessCookieName)
	if err != nil {
		return "", false
	}
	return strings.TrimSpace(raw), true
}

// isSafeMethod 判断请求方法是否为“无副作用”方法。
func isSafeMethod(method string) bool {
	switch strings.ToUpper(strings.TrimSpace(method)) {
	case http.MethodGet, http.MethodHead, http.MethodOptions:
		return true
	default:
		return false
	}
}
