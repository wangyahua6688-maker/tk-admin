package middleware

import (
	"github.com/gin-gonic/gin"
	"go-admin-full/internal/services"
)

// PermissionRequired 按“权限码”做 RBAC 鉴权。
// 鉴权流程：
// 1. 从 JWT 中间件注入的 uid 中获取用户身份；
// 2. 查询用户全部角色（并预加载角色权限）；
// 3. admin 角色直接放行；
// 4. 聚合权限码集合，判断目标权限码是否存在。
func PermissionRequired(code string, userRoleSvc *services.UserRoleService) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 1) 必须有已认证用户。
		uid := c.GetUint("uid")
		if uid == 0 {
			c.AbortWithStatusJSON(401, gin.H{"msg": "unauthorized"})
			return
		}

		// 2) 读取用户角色与角色权限。
		roles, err := userRoleSvc.GetUserRoles(c.Request.Context(), uid)
		if err != nil {
			c.AbortWithStatusJSON(403, gin.H{"msg": "no roles"})
			return
		}
		if len(roles) == 0 {
			c.AbortWithStatusJSON(403, gin.H{"msg": "no roles"})
			return
		}

		// 3) 超级管理员直通。
		for _, r := range roles {
			if r.Code == "admin" {
				c.Next()
				return
			}
		}

		// 4) 将权限列表拍平成集合，O(1) 判断目标权限。
		permSet := map[string]bool{}
		for _, r := range roles {
			for _, p := range r.Permissions {
				permSet[p.Code] = true
			}
		}

		// 5) 不具备目标权限则拒绝访问。
		if !permSet[code] {
			c.AbortWithStatusJSON(403, gin.H{"msg": "permission denied"})
			return
		}

		c.Next()
	}
}
