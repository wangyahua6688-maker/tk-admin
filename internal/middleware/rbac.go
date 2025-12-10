package middleware

import (
	"go-admin-full/internal/models"
	"net/http"

	"github.com/gin-gonic/gin"
	"go-admin-full/internal/services"
)

type RBACMiddleware struct {
	authService       *services.AuthService
	roleService       *services.RoleService
	permissionService *services.PermissionService
}

func NewRBACMiddleware(
	auth *services.AuthService,
	role *services.RoleService,
	perm *services.PermissionService,
) *RBACMiddleware {
	return &RBACMiddleware{
		authService:       auth,
		roleService:       role,
		permissionService: perm,
	}
}

func (m *RBACMiddleware) Handler() gin.HandlerFunc {
	return func(c *gin.Context) {
		// RBAC 中间件要求先经过 JWT 中间件（从上下文获取 uid）
		uid, ok := c.Get("uid")
		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"msg": "缺少认证上下文"})
			return
		}
		userID, ok := uid.(uint)
		if !ok || userID == 0 {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"msg": "无效认证上下文"})
			return
		}

		// 3. 查询用户拥有的角色
		roles, err := m.roleService.GetRolesByUserID(c.Request.Context(), userID)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"msg": "无法加载角色"})
			return
		}

		// 4. 查询角色拥有的权限
		permissions, err := m.permissionService.GetPermissionsByRoleIDs(
			c.Request.Context(),
			extractRoleIDs(roles),
		)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"msg": "无法加载权限"})
			return
		}

		// 5. 判断当前请求路径是否允许访问
		reqPath := c.FullPath()

		if !isAllowed(reqPath, permissions) {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"msg": "无访问权限"})
			return
		}

		// 6. 放行
		c.Next()
	}
}

// 提取角色 ID
func extractRoleIDs(list []models.Role) []uint {
	ids := make([]uint, len(list))
	for i, r := range list {
		ids[i] = r.ID
	}
	return ids
}

// 校验接口是否在权限表内
func isAllowed(path string, perms []models.Permission) bool {
	for _, p := range perms {
		if p.Path == path { // Permission.Path = "/api/user/list"
			return true
		}
	}
	return false
}
