package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"go-admin-full/internal/models"
)

func RBACMiddleware(db *gorm.DB, requiredCode string) gin.HandlerFunc {
	return func(c *gin.Context) {
		uidv, exists := c.Get("uid")
		if !exists {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"code": 401, "msg": "未认证用户"})
			return
		}
		uid := uidv.(uint)

		var user models.User
		if err := db.WithContext(c.Request.Context()).Preload("Roles.Permissions").First(&user, uid).Error; err != nil {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"code": 403, "msg": "用户不存在"})
			return
		}

		perm := map[string]bool{}
		for _, role := range user.Roles {
			for _, p := range role.Permissions {
				perm[p.Code] = true
			}
		}
		if !perm[requiredCode] {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"code": 403, "msg": "无权限访问"})
			return
		}
		c.Next()
	}
}
