// Package rediskey 定义 tk-admin 服务使用的全部 Redis key 模板。
package rediskey

import (
	"fmt"
	"time"
)

// ─────────────────────────────────────────────
// JWT Token 存储（与 token/store 模块配合）
// ─────────────────────────────────────────────

// KeyRefreshToken RefreshToken 存储 key（user_id + device_id 维度）。
// 格式: refresh_token:{userID}:{deviceID}
// TTL:  对齐 RefreshExpire 配置
func KeyRefreshToken(userID uint, deviceID string) string {
	return fmt.Sprintf("refresh_token:%d:%s", userID, deviceID)
}

// KeyAccessBlacklist 已撤销 AccessToken 的 JTI 黑名单。
// 格式: blacklist:access:{jti}
// TTL:  对齐 Token 剩余有效期
func KeyAccessBlacklist(jti string) string {
	return fmt.Sprintf("blacklist:access:%s", jti)
}

// KeySessionActivity 会话空闲活跃标记（用于空闲超时控制）。
// 格式: session:active:{userID}:{deviceID}
// TTL:  SessionIdleTimeout 配置项
func KeySessionActivity(userID uint, deviceID string) string {
	return fmt.Sprintf("session:active:%d:%s", userID, deviceID)
}

// ─────────────────────────────────────────────
// 登录失败频控（管理员账号）
// ─────────────────────────────────────────────

const (
	// LoginFailMax 管理员登录失败最大次数（超出后锁定）
	LoginFailMax = 5
	// LoginFailWindow 登录失败计数窗口（15 分钟）
	LoginFailWindow = 15 * time.Minute
)

// KeyLoginFail 管理员登录失败计数。
// 格式: auth:login_fail:{username}:{clientIP}
// TTL:  LoginFailWindow
func KeyLoginFail(username, ip string) string {
	return fmt.Sprintf("auth:login_fail:%s:%s", username, ip)
}

// ─────────────────────────────────────────────
// RBAC 权限缓存
// ─────────────────────────────────────────────

const (
	// RBACPermsTTL 用户权限集合缓存有效期（5 分钟）
	// 角色/权限变更后需主动删除此 key
	RBACPermsTTL = 5 * time.Minute
	// RBACMenusTTL 用户菜单树缓存有效期（5 分钟）
	RBACMenusTTL = 5 * time.Minute
)

// KeyRBACPerms 管理员用户权限码集合缓存（JSON 序列化的 []string）。
// 格式: rbac:perms:{userID}
// TTL:  RBACPermsTTL
func KeyRBACPerms(userID uint) string {
	return fmt.Sprintf("rbac:perms:%d", userID)
}

// KeyRBACMenus 管理员用户菜单树缓存。
// 格式: rbac:menus:{userID}
// TTL:  RBACMenusTTL
func KeyRBACMenus(userID uint) string {
	return fmt.Sprintf("rbac:menus:%d", userID)
}

// ─────────────────────────────────────────────
// 幂等去重
// ─────────────────────────────────────────────

const (
	// IdempotentTTL 幂等 key 有效期（60 秒）
	IdempotentTTL = 60 * time.Second
)

// KeyIdempotent 接口幂等去重（基于 X-Request-ID）。
// 格式: idempotent:admin:{requestID}
// TTL:  IdempotentTTL
func KeyIdempotent(requestID string) string {
	return fmt.Sprintf("idempotent:admin:%s", requestID)
}
