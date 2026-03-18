package middleware

import (
	"github.com/gin-gonic/gin"
	rbacsvc "go-admin-full/internal/services/rbac"
)

// PermissionRequired 按“权限码”做 RBAC 鉴权。
// 鉴权流程：
// 1. 从 JWT 中间件注入的 uid 中获取用户身份；
// 2. 查询用户全部角色（并预加载角色权限）；
// 3. admin 角色直接放行；
// 4. 聚合权限码集合，判断目标权限码是否存在。
func PermissionRequired(code string, userRoleSvc *rbacsvc.UserRoleService) gin.HandlerFunc {
	// 返回当前处理结果。
	return func(c *gin.Context) {
		// 1) 必须有已认证用户。
		uid := c.GetUint("uid")
		// 判断条件并进入对应分支逻辑。
		if uid == 0 {
			// 调用c.AbortWithStatusJSON完成当前处理。
			c.AbortWithStatusJSON(401, gin.H{"msg": "unauthorized"})
			// 返回当前处理结果。
			return
		}

		// 2) 读取用户角色与角色权限。
		roles, err := userRoleSvc.GetUserRoles(c.Request.Context(), uid)
		// 判断条件并进入对应分支逻辑。
		if err != nil {
			// 调用c.AbortWithStatusJSON完成当前处理。
			c.AbortWithStatusJSON(403, gin.H{"msg": "no roles"})
			// 返回当前处理结果。
			return
		}
		// 判断条件并进入对应分支逻辑。
		if len(roles) == 0 {
			// 调用c.AbortWithStatusJSON完成当前处理。
			c.AbortWithStatusJSON(403, gin.H{"msg": "no roles"})
			// 返回当前处理结果。
			return
		}

		// 3) 超级管理员直通。
		for _, r := range roles {
			// 判断条件并进入对应分支逻辑。
			if r.Code == "admin" {
				// 调用c.Next完成当前处理。
				c.Next()
				// 返回当前处理结果。
				return
			}
		}

		// 4) 将权限列表拍平成集合，O(1) 判断目标权限。
		permSet := map[string]bool{}
		// 循环处理当前数据集合。
		for _, r := range roles {
			// 循环处理当前数据集合。
			for _, p := range r.Permissions {
				// 更新当前变量或字段值。
				permSet[p.Code] = true
			}
		}

		// 5) 不具备目标权限则拒绝访问。
		if !permSet[code] {
			// 调用c.AbortWithStatusJSON完成当前处理。
			c.AbortWithStatusJSON(403, gin.H{"msg": "permission denied"})
			// 返回当前处理结果。
			return
		}

		// 调用c.Next完成当前处理。
		c.Next()
	}
}
