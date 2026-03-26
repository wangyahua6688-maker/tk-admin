package jwt

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"go-admin-full/internal/constants"
	storepkg "go-admin-full/internal/token/store"
	"strings"
	"time"

	gjwt "github.com/golang-jwt/jwt/v5"
)

// 声明当前常量。
const (
	// TokenTypeAccess 用于接口访问鉴权，生命周期较短。
	TokenTypeAccess = "access"
	// TokenTypeRefresh 用于换发新的 access token，生命周期较长。
	TokenTypeRefresh = "refresh"
	// 更新当前变量或字段值。
	defaultDeviceID = "default"
)

// Manager 统一封装 JWT 的签发、校验、刷新和撤销逻辑。
type Manager struct {
	// 处理当前语句逻辑。
	Config *Config
	// 处理当前语句逻辑。
	Store storepkg.Store
}

// NewManager 创建 token 管理器。
// store 可选：传入后可启用 refresh token 持久化、access token 黑名单等能力。
func NewManager(cfg *Config, store storepkg.Store) *Manager {
	// 返回当前处理结果。
	return &Manager{
		// 处理当前语句逻辑。
		Config: cfg,
		// 处理当前语句逻辑。
		Store: store,
	}
}

// GenerateTokensWithDevice 按设备粒度生成 access/refresh token。
// 设计目的：
// - 同一账号可以多端登录；
// - 登出时只失效“当前设备”的 refresh token，而不是全部设备。
func (m *Manager) GenerateTokensWithDevice(userID uint, deviceID string) (access string, refresh string, err error) {
	// 定义并初始化当前变量。
	now := time.Now()
	// 更新当前变量或字段值。
	deviceID = normalizeDeviceID(deviceID)

	// 1) 生成 access token（短期凭证）。
	accessJTI, err := newJTI()
	// 判断条件并进入对应分支逻辑。
	if err != nil {
		// 返回当前处理结果。
		return "", "", fmt.Errorf("%w: %v", constants.ErrSigningToken, err)
	}

	// 定义并初始化当前变量。
	accessClaims := Claims{
		// 处理当前语句逻辑。
		UserID: userID,
		// 处理当前语句逻辑。
		TokenType: TokenTypeAccess,
		// 处理当前语句逻辑。
		DeviceID: deviceID,
		// 进入新的代码块进行处理。
		RegisteredClaims: gjwt.RegisteredClaims{
			// 处理当前语句逻辑。
			Issuer: m.Config.Issuer,
			// 处理当前语句逻辑。
			ID: accessJTI,
			// 调用gjwt.NewNumericDate完成当前处理。
			IssuedAt: gjwt.NewNumericDate(now),
			// 调用gjwt.NewNumericDate完成当前处理。
			ExpiresAt: gjwt.NewNumericDate(now.Add(m.Config.AccessExpire)),
		},
	}
	// 更新当前变量或字段值。
	access, err = m.signClaims(accessClaims)
	// 判断条件并进入对应分支逻辑。
	if err != nil {
		// 返回当前处理结果。
		return "", "", fmt.Errorf("%w: %v", constants.ErrSigningToken, err)
	}

	// 2) 生成 refresh token（长期凭证）。
	refreshJTI, err := newJTI()
	// 判断条件并进入对应分支逻辑。
	if err != nil {
		// 返回当前处理结果。
		return "", "", fmt.Errorf("%w: %v", constants.ErrSigningToken, err)
	}

	// 定义并初始化当前变量。
	refreshClaims := Claims{
		// 处理当前语句逻辑。
		UserID: userID,
		// 处理当前语句逻辑。
		TokenType: TokenTypeRefresh,
		// 处理当前语句逻辑。
		DeviceID: deviceID,
		// 进入新的代码块进行处理。
		RegisteredClaims: gjwt.RegisteredClaims{
			// 处理当前语句逻辑。
			Issuer: m.Config.Issuer,
			// 处理当前语句逻辑。
			ID: refreshJTI,
			// 调用gjwt.NewNumericDate完成当前处理。
			IssuedAt: gjwt.NewNumericDate(now),
			// 调用gjwt.NewNumericDate完成当前处理。
			ExpiresAt: gjwt.NewNumericDate(now.Add(m.Config.RefreshExpire)),
		},
	}
	// 更新当前变量或字段值。
	refresh, err = m.signClaims(refreshClaims)
	// 判断条件并进入对应分支逻辑。
	if err != nil {
		// 返回当前处理结果。
		return "", "", fmt.Errorf("%w: %v", constants.ErrSigningToken, err)
	}

	// 3) 将 refresh token 持久化（若有 Store）。
	// key 维度：user_id + device_id。
	if m.Store != nil {
		// 定义并初始化当前变量。
		key := m.getRefreshTokenKey(userID, deviceID)
		// 判断条件并进入对应分支逻辑。
		if serr := m.Store.Set(key, refresh, m.Config.RefreshExpire); serr != nil {
			// 返回当前处理结果。
			return "", "", constants.ErrTokenStoreFailed
		}

		// 同步刷新“会话活跃标记”，用于 1 小时空闲自动失效。
		if m.Config.SessionIdleTimeout > 0 {
			// 定义并初始化当前变量。
			sessionKey := m.getSessionActivityKey(userID, deviceID)
			// 判断条件并进入对应分支逻辑。
			if serr := m.Store.Set(sessionKey, "1", m.Config.SessionIdleTimeout); serr != nil {
				// 返回当前处理结果。
				return "", "", constants.ErrTokenStoreFailed
			}
		}
	}
	// 返回当前处理结果。
	return access, refresh, nil
}

// RefreshTokenPair 校验并轮换 refresh token，返回新的 access+refresh。
// 安全要点：
// - 强制 token_type=refresh；
// - 强制 issuer 匹配；
// - 与 Store 中保存值比对，防止伪造或被旧值重放。
func (m *Manager) RefreshTokenPair(refreshToken string) (string, string, error) {
	// 定义并初始化当前变量。
	claims, err := ParseTokenClaims(refreshToken, m.Config.SigningKey)
	// 判断条件并进入对应分支逻辑。
	if err != nil {
		// 返回当前处理结果。
		return "", "", err
	}
	// 判断条件并进入对应分支逻辑。
	if claims.TokenType != TokenTypeRefresh {
		// 返回当前处理结果。
		return "", "", constants.ErrInvalidToken
	}
	// 判断条件并进入对应分支逻辑。
	if m.Config.Issuer != "" && claims.Issuer != m.Config.Issuer {
		// 返回当前处理结果。
		return "", "", constants.ErrInvalidToken
	}

	// 定义并初始化当前变量。
	deviceID := normalizeDeviceID(claims.DeviceID)
	// 判断条件并进入对应分支逻辑。
	if m.Store != nil {
		// 定义并初始化当前变量。
		key := m.getRefreshTokenKey(claims.UserID, deviceID)
		// 定义并初始化当前变量。
		stored, err := m.Store.Get(key)
		// 判断条件并进入对应分支逻辑。
		if err != nil {
			// 返回当前处理结果。
			return "", "", constants.ErrTokenNotFound
		}
		// 判断条件并进入对应分支逻辑。
		if stored != refreshToken {
			// 返回当前处理结果。
			return "", "", constants.ErrInvalidToken
		}

		// 空闲超时控制：refresh 前必须仍在活跃会话窗口内。
		if m.Config.SessionIdleTimeout > 0 {
			// 定义并初始化当前变量。
			sessionKey := m.getSessionActivityKey(claims.UserID, deviceID)
			// 判断条件并进入对应分支逻辑。
			if _, err := m.Store.Get(sessionKey); err != nil {
				// 判断条件并进入对应分支逻辑。
				if err == constants.ErrTokenNotFound {
					// 返回当前处理结果。
					return "", "", constants.ErrExpiredToken
				}
				// 返回当前处理结果。
				return "", "", constants.ErrTokenStoreFailed
			}
		}
	}

	// 定义并初始化当前变量。
	access, refresh, err := m.GenerateTokensWithDevice(claims.UserID, deviceID)
	// 判断条件并进入对应分支逻辑。
	if err != nil {
		// 返回当前处理结果。
		return "", "", err
	}
	// 返回当前处理结果。
	return access, refresh, nil
}

// ValidateAccessToken 校验 access token，并检查是否在黑名单。
// 黑名单用于“主动失效”场景（如用户登出后立即让 access token 无效）。
func (m *Manager) ValidateAccessToken(accessToken string) (*Claims, error) {
	// 定义并初始化当前变量。
	claims, err := ParseTokenClaims(accessToken, m.Config.SigningKey)
	// 判断条件并进入对应分支逻辑。
	if err != nil {
		// 返回当前处理结果。
		return nil, err
	}
	// 判断条件并进入对应分支逻辑。
	if claims.TokenType != TokenTypeAccess {
		// 返回当前处理结果。
		return nil, constants.ErrInvalidToken
	}
	// 判断条件并进入对应分支逻辑。
	if m.Config.Issuer != "" && claims.Issuer != m.Config.Issuer {
		// 返回当前处理结果。
		return nil, constants.ErrInvalidToken
	}

	// 判断条件并进入对应分支逻辑。
	if claims.ID == "" {
		// 返回当前处理结果。
		return nil, constants.ErrInvalidToken
	}

	// 判断条件并进入对应分支逻辑。
	if m.Store != nil {
		// 判断条件并进入对应分支逻辑。
		if _, err := m.Store.Get(m.getAccessBlacklistKey(claims.ID)); err == nil {
			// 返回当前处理结果。
			return nil, constants.ErrInvalidToken
			// 进入新的代码块进行处理。
		} else if err != constants.ErrTokenNotFound {
			// 返回当前处理结果。
			return nil, constants.ErrInvalidToken
		}

		// 会话空闲续期：
		// 1. access token 每次通过鉴权都触发会话续约；
		// 2. 若活跃键不存在，视为超过空闲窗口，直接判定失效。
		if m.Config.SessionIdleTimeout > 0 {
			// 定义并初始化当前变量。
			deviceID := normalizeDeviceID(claims.DeviceID)
			// 定义并初始化当前变量。
			sessionKey := m.getSessionActivityKey(claims.UserID, deviceID)
			// 判断条件并进入对应分支逻辑。
			if _, err := m.Store.Get(sessionKey); err != nil {
				// 判断条件并进入对应分支逻辑。
				if err == constants.ErrTokenNotFound {
					// 返回当前处理结果。
					return nil, constants.ErrExpiredToken
				}
				// 返回当前处理结果。
				return nil, constants.ErrTokenStoreFailed
			}
			// 判断条件并进入对应分支逻辑。
			if err := m.Store.Set(sessionKey, "1", m.Config.SessionIdleTimeout); err != nil {
				// 返回当前处理结果。
				return nil, constants.ErrTokenStoreFailed
			}
		}
	}

	// 返回当前处理结果。
	return claims, nil
}

// RevokeAccessToken 将 access token 的 jti 写入黑名单，TTL 为 token 剩余生命周期。
func (m *Manager) RevokeAccessToken(accessToken string) error {
	// 判断条件并进入对应分支逻辑。
	if m.Store == nil {
		// 返回当前处理结果。
		return nil
	}
	// 定义并初始化当前变量。
	claims, err := ParseTokenClaims(accessToken, m.Config.SigningKey)
	// 判断条件并进入对应分支逻辑。
	if err != nil {
		// 返回当前处理结果。
		return err
	}
	// 判断条件并进入对应分支逻辑。
	if claims.TokenType != TokenTypeAccess || claims.ID == "" || claims.ExpiresAt == nil {
		// 返回当前处理结果。
		return constants.ErrInvalidToken
	}

	// 定义并初始化当前变量。
	ttl := time.Until(claims.ExpiresAt.Time)
	// 判断条件并进入对应分支逻辑。
	if ttl <= 0 {
		// 返回当前处理结果。
		return nil
	}

	// 返回当前处理结果。
	return m.Store.Set(m.getAccessBlacklistKey(claims.ID), "1", ttl)
}

// InvalidateRefresh 使某个用户某个设备的 refresh token 失效（用于登出）。
func (m *Manager) InvalidateRefresh(userID uint, deviceID string) error {
	// 判断条件并进入对应分支逻辑。
	if m.Store == nil {
		// 返回当前处理结果。
		return nil
	}
	// 定义并初始化当前变量。
	normalized := normalizeDeviceID(deviceID)
	// 判断条件并进入对应分支逻辑。
	if err := m.Store.Delete(m.getRefreshTokenKey(userID, normalized)); err != nil {
		// 返回当前处理结果。
		return err
	}
	// 判断条件并进入对应分支逻辑。
	if m.Config.SessionIdleTimeout > 0 {
		// 判断条件并进入对应分支逻辑。
		if err := m.Store.Delete(m.getSessionActivityKey(userID, normalized)); err != nil {
			// 返回当前处理结果。
			return err
		}
	}
	// 返回当前处理结果。
	return nil
}

// signClaims 处理signClaims相关逻辑。
func (m *Manager) signClaims(c Claims) (string, error) {
	// 定义并初始化当前变量。
	token := gjwt.NewWithClaims(gjwt.SigningMethodHS256, c)
	// 返回当前处理结果。
	return token.SignedString([]byte(m.Config.SigningKey))
}

// getRefreshTokenKey 处理getRefreshTokenKey相关逻辑。
func (m *Manager) getRefreshTokenKey(userID uint, deviceID string) string {
	// 返回当前处理结果。
	return fmt.Sprintf("refresh_token:%d:%s", userID, deviceID)
}

// getAccessBlacklistKey 处理getAccessBlacklistKey相关逻辑。
func (m *Manager) getAccessBlacklistKey(jti string) string {
	// 返回当前处理结果。
	return fmt.Sprintf("blacklist:access:%s", jti)
}

// getSessionActivityKey 返回会话活跃标记存储 key。
func (m *Manager) getSessionActivityKey(userID uint, deviceID string) string {
	// 返回当前处理结果。
	return fmt.Sprintf("session:active:%d:%s", userID, deviceID)
}

// normalizeDeviceID 处理normalizeDeviceID相关逻辑。
func normalizeDeviceID(deviceID string) string {
	// 定义并初始化当前变量。
	trimmed := strings.ToLower(strings.TrimSpace(deviceID))
	// 判断条件并进入对应分支逻辑。
	if trimmed == "" {
		// 返回当前处理结果。
		return defaultDeviceID
	}
	// 若输入已经是 32 位十六进制（本项目归一化后的 deviceID），直接返回避免二次哈希。
	if isHexToken(trimmed, 32) {
		// 返回当前处理结果。
		return trimmed
	}
	// 对 device 标识做哈希，避免把原始 UA/设备信息直接作为 Redis key 暴露。
	sum := sha256.Sum256([]byte(trimmed))
	// 返回当前处理结果。
	return hex.EncodeToString(sum[:16])
}

// isHexToken 判断字符串是否为固定长度十六进制文本。
func isHexToken(v string, expectLen int) bool {
	// 判断条件并进入对应分支逻辑。
	if len(v) != expectLen {
		// 返回当前处理结果。
		return false
	}
	// 循环处理当前数据集合。
	for _, ch := range v {
		// 判断条件并进入对应分支逻辑。
		if !(ch >= '0' && ch <= '9') && !(ch >= 'a' && ch <= 'f') {
			// 返回当前处理结果。
			return false
		}
	}
	// 返回当前处理结果。
	return true
}

// newJTI 处理newJTI相关逻辑。
func newJTI() (string, error) {
	// 定义并初始化当前变量。
	buf := make([]byte, 16)
	// 判断条件并进入对应分支逻辑。
	if _, err := rand.Read(buf); err != nil {
		// 返回当前处理结果。
		return "", err
	}
	// 返回当前处理结果。
	return hex.EncodeToString(buf), nil
}
