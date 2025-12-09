package middleware

import (
	"github.com/gin-gonic/gin"
	"go-admin-full/internal/services"
)

func PermissionRequired(code string, userRoleSvc *services.UserRoleService) gin.HandlerFunc {
	return func(c *gin.Context) {
		uid := c.GetUint("uid")

		roles, err := userRoleSvc.GetUserRoles(c.Request.Context(), uid)
		if err != nil {
			c.JSON(403, gin.H{"msg": "no roles"})
			c.Abort()
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
			c.JSON(403, gin.H{"msg": "permission denied"})
			c.Abort()
			return
		}

		c.Next()
	}
}
