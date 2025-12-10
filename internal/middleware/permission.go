package middleware

import (
	"github.com/gin-gonic/gin"
	"go-admin-full/internal/services"
)

func PermissionRequired(code string, userRoleSvc *services.UserRoleService) gin.HandlerFunc {
	return func(c *gin.Context) {
		uid := c.GetUint("uid")
		if uid == 0 {
			c.AbortWithStatusJSON(401, gin.H{"msg": "unauthorized"})
			return
		}

		roles, err := userRoleSvc.GetUserRoles(c.Request.Context(), uid)
		if err != nil {
			c.AbortWithStatusJSON(403, gin.H{"msg": "no roles"})
			return
		}
		if len(roles) == 0 {
			c.AbortWithStatusJSON(403, gin.H{"msg": "no roles"})
			return
		}

		// Super Admin bypass
		for _, r := range roles {
			if r.Code == "admin" {
				c.Next()
				return
			}
		}

		permSet := map[string]bool{}
		for _, r := range roles {
			for _, p := range r.Permissions {
				permSet[p.Code] = true
			}
		}

		if !permSet[code] {
			c.AbortWithStatusJSON(403, gin.H{"msg": "permission denied"})
			return
		}

		c.Next()
	}
}
