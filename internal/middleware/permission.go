package middleware

import (
	"encoding/json"

	"github.com/gin-gonic/gin"
	"go-admin/internal/constants/rediskey"
	rbacsvc "go-admin/internal/services/rbac"
	tokenjwt "go-admin/internal/token/jwt"
)

// PermissionRequired 按"权限码"做 RBAC 鉴权（含 Redis 权限缓存）。
//
// 鉴权流程：
//  1. 从 JWT 中间件注入的 uid 中获取用户身份
//  2. 先查 Redis 权限缓存（TTL=5min），命中直接判断
//  3. 缓存 miss 时查 DB，并将结果写入缓存
//  4. admin 角色直接放行
//  5. 聚合权限码集合，判断目标权限码是否存在
func PermissionRequired(code string, userRoleSvc *rbacsvc.UserRoleService, tokenMgr *tokenjwt.Manager) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 1) 必须有已认证用户
		uid := c.GetUint("uid")
		if uid == 0 {
			c.AbortWithStatusJSON(401, gin.H{"msg": "unauthorized"})
			return
		}

		ctx := c.Request.Context()

		// 2) 尝试从 Redis 缓存读取权限码集合
		if tokenMgr != nil && tokenMgr.Store != nil {
			permKey := rediskey.KeyRBACPerms(uid)
			cached, err := tokenMgr.Store.Get(permKey)
			if err == nil && cached != "" {
				// 缓存命中：反序列化权限码集合
				var permSet map[string]bool
				if jsonErr := json.Unmarshal([]byte(cached), &permSet); jsonErr == nil {
					// 检查 admin 超级权限（admin key 特殊处理）
					if permSet["__admin__"] {
						c.Next()
						return
					}
					if !permSet[code] {
						c.AbortWithStatusJSON(403, gin.H{"msg": "permission denied"})
						return
					}
					c.Next()
					return
				}
			}
		}

		// 3) 缓存 miss：查询 DB 获取角色和权限
		roles, err := userRoleSvc.GetUserRoles(ctx, uid)
		if err != nil || len(roles) == 0 {
			c.AbortWithStatusJSON(403, gin.H{"msg": "no roles"})
			return
		}

		// 4) 构建权限码集合
		permSet := map[string]bool{}
		isAdmin := false
		for _, r := range roles {
			if r.Code == "admin" {
				isAdmin = true
				// admin 标记写入缓存集合
				permSet["__admin__"] = true
			}
			for _, p := range r.Permissions {
				permSet[p.Code] = true
			}
		}

		// 5) 将权限集合写入 Redis 缓存（异步写，不阻塞鉴权流程）
		if tokenMgr != nil && tokenMgr.Store != nil {
			if raw, jsonErr := json.Marshal(permSet); jsonErr == nil {
				permKey := rediskey.KeyRBACPerms(uid)
				_ = tokenMgr.Store.Set(permKey, string(raw), rediskey.RBACPermsTTL)
			}
		}

		// 6) 权限判断
		if isAdmin {
			c.Next()
			return
		}
		if !permSet[code] {
			c.AbortWithStatusJSON(403, gin.H{"msg": "permission denied"})
			return
		}
		c.Next()
	}
}

// InvalidateUserPermCache 主动清除指定用户的权限缓存（角色/权限变更后调用）。
func InvalidateUserPermCache(tokenMgr *tokenjwt.Manager, userID uint) {
	if tokenMgr == nil || tokenMgr.Store == nil {
		return
	}
	// 同时清除权限缓存和菜单缓存
	_ = tokenMgr.Store.Delete(rediskey.KeyRBACPerms(userID))
	_ = tokenMgr.Store.Delete(rediskey.KeyRBACMenus(userID))
}
