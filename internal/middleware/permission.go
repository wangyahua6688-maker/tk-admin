package middleware

import (
	"github.com/gin-gonic/gin"
	"go-admin-full/internal/services"
)

func PermissionRequired(code string, svc *services.UserRoleService) gin.HandlerFunc {
	return func(c *gin.Context) {
		uid := c.GetUint("uid")

		roles, err := svc.GetUserRoles(c.Request.Context(), uid)
		if err != nil {
			c.JSON(403, gin.H{"msg": "no roles"})
			c.Abort()
			return
		}

		// 收集角色所有权限
		permSet := map[string]bool{}
		// 超级管理员
		for _, r := range roles {
			if r.Code == "admin" {
				c.Next()
				return
			}
		}

		if !permSet[code] {
			c.JSON(403, gin.H{"msg": "permission denied"})
			c.Abort()
			return
		}

		c.Next()
	}
}
