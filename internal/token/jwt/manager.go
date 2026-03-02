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

const (
	// TokenTypeAccess 用于接口访问鉴权，生命周期较短。
	TokenTypeAccess = "access"
	// TokenTypeRefresh 用于换发新的 access token，生命周期较长。
	TokenTypeRefresh = "refresh"
	defaultDeviceID  = "default"
)

// Manager 统一封装 JWT 的签发、校验、刷新和撤销逻辑。
type Manager struct {
	Config *Config
	Store  storepkg.Store
}

// NewManager 创建 token 管理器。
// store 可选：传入后可启用 refresh token 持久化、access token 黑名单等能力。
func NewManager(cfg *Config, store storepkg.Store) *Manager {
	return &Manager{
		Config: cfg,
		Store:  store,
	}
}

// GenerateTokens 为指定用户生成 access/refresh 一对 token。
// 默认使用 defaultDeviceID，当你不关心设备粒度时可直接调用。
func (m *Manager) GenerateTokens(userID uint) (access string, refresh string, err error) {
	return m.GenerateTokensWithDevice(userID, defaultDeviceID)
}

// GenerateTokensWithDevice 按设备粒度生成 access/refresh token。
// 设计目的：
// - 同一账号可以多端登录；
// - 登出时只失效“当前设备”的 refresh token，而不是全部设备。
func (m *Manager) GenerateTokensWithDevice(userID uint, deviceID string) (access string, refresh string, err error) {
	now := time.Now()
	deviceID = normalizeDeviceID(deviceID)

	// 1) 生成 access token（短期凭证）。
	accessJTI, err := newJTI()
	if err != nil {
		return "", "", fmt.Errorf("%w: %v", constants.ErrSigningToken, err)
	}

	accessClaims := Claims{
		UserID:    userID,
		TokenType: TokenTypeAccess,
		DeviceID:  deviceID,
		RegisteredClaims: gjwt.RegisteredClaims{
			Issuer:    m.Config.Issuer,
			ID:        accessJTI,
			IssuedAt:  gjwt.NewNumericDate(now),
			ExpiresAt: gjwt.NewNumericDate(now.Add(m.Config.AccessExpire)),
		},
	}
	access, err = m.signClaims(accessClaims)
	if err != nil {
		return "", "", fmt.Errorf("%w: %v", constants.ErrSigningToken, err)
	}

	// 2) 生成 refresh token（长期凭证）。
	refreshJTI, err := newJTI()
	if err != nil {
		return "", "", fmt.Errorf("%w: %v", constants.ErrSigningToken, err)
	}

	refreshClaims := Claims{
		UserID:    userID,
		TokenType: TokenTypeRefresh,
		DeviceID:  deviceID,
		RegisteredClaims: gjwt.RegisteredClaims{
			Issuer:    m.Config.Issuer,
			ID:        refreshJTI,
			IssuedAt:  gjwt.NewNumericDate(now),
			ExpiresAt: gjwt.NewNumericDate(now.Add(m.Config.RefreshExpire)),
		},
	}
	refresh, err = m.signClaims(refreshClaims)
	if err != nil {
		return "", "", fmt.Errorf("%w: %v", constants.ErrSigningToken, err)
	}

	// 3) 将 refresh token 持久化（若有 Store）。
	// key 维度：user_id + device_id。
	if m.Store != nil {
		key := m.getRefreshTokenKey(userID, deviceID)
		if serr := m.Store.Set(key, refresh, m.Config.RefreshExpire); serr != nil {
			return "", "", constants.ErrTokenStoreFailed
		}
	}
	return access, refresh, nil
}

// RefreshToken 仅返回新的 access token（兼容旧调用方式）。
func (m *Manager) RefreshToken(refreshToken string) (string, error) {
	access, _, err := m.RefreshTokenPair(refreshToken)
	return access, err
}

// RefreshTokenPair 校验并轮换 refresh token，返回新的 access+refresh。
// 安全要点：
// - 强制 token_type=refresh；
// - 强制 issuer 匹配；
// - 与 Store 中保存值比对，防止伪造或被旧值重放。
func (m *Manager) RefreshTokenPair(refreshToken string) (string, string, error) {
	claims, err := ParseTokenClaims(refreshToken, m.Config.SigningKey)
	if err != nil {
		return "", "", err
	}
	if claims.TokenType != TokenTypeRefresh {
		return "", "", constants.ErrInvalidToken
	}
	if m.Config.Issuer != "" && claims.Issuer != m.Config.Issuer {
		return "", "", constants.ErrInvalidToken
	}

	deviceID := normalizeDeviceID(claims.DeviceID)
	if m.Store != nil {
		key := m.getRefreshTokenKey(claims.UserID, deviceID)
		stored, err := m.Store.Get(key)
		if err != nil {
			return "", "", constants.ErrTokenNotFound
		}
		if stored != refreshToken {
			return "", "", constants.ErrInvalidToken
		}
	}

	access, refresh, err := m.GenerateTokensWithDevice(claims.UserID, deviceID)
	if err != nil {
		return "", "", err
	}
	return access, refresh, nil
}

// ValidateAccessToken 校验 access token，并检查是否在黑名单。
// 黑名单用于“主动失效”场景（如用户登出后立即让 access token 无效）。
func (m *Manager) ValidateAccessToken(accessToken string) (*Claims, error) {
	claims, err := ParseTokenClaims(accessToken, m.Config.SigningKey)
	if err != nil {
		return nil, err
	}
	if claims.TokenType != TokenTypeAccess {
		return nil, constants.ErrInvalidToken
	}
	if m.Config.Issuer != "" && claims.Issuer != m.Config.Issuer {
		return nil, constants.ErrInvalidToken
	}

	if claims.ID == "" {
		return nil, constants.ErrInvalidToken
	}

	if m.Store != nil {
		if _, err := m.Store.Get(m.getAccessBlacklistKey(claims.ID)); err == nil {
			return nil, constants.ErrInvalidToken
		} else if err != constants.ErrTokenNotFound {
			return nil, constants.ErrInvalidToken
		}
	}

	return claims, nil
}

// RevokeAccessToken 将 access token 的 jti 写入黑名单，TTL 为 token 剩余生命周期。
func (m *Manager) RevokeAccessToken(accessToken string) error {
	if m.Store == nil {
		return nil
	}
	claims, err := ParseTokenClaims(accessToken, m.Config.SigningKey)
	if err != nil {
		return err
	}
	if claims.TokenType != TokenTypeAccess || claims.ID == "" || claims.ExpiresAt == nil {
		return constants.ErrInvalidToken
	}

	ttl := time.Until(claims.ExpiresAt.Time)
	if ttl <= 0 {
		return nil
	}

	return m.Store.Set(m.getAccessBlacklistKey(claims.ID), "1", ttl)
}

// InvalidateRefresh 使某个用户某个设备的 refresh token 失效（用于登出）。
func (m *Manager) InvalidateRefresh(userID uint, deviceID string) error {
	if m.Store == nil {
		return nil
	}
	return m.Store.Delete(m.getRefreshTokenKey(userID, normalizeDeviceID(deviceID)))
}

func (m *Manager) signClaims(c Claims) (string, error) {
	token := gjwt.NewWithClaims(gjwt.SigningMethodHS256, c)
	return token.SignedString([]byte(m.Config.SigningKey))
}

func (m *Manager) getRefreshTokenKey(userID uint, deviceID string) string {
	return fmt.Sprintf("refresh_token:%d:%s", userID, deviceID)
}

func (m *Manager) getAccessBlacklistKey(jti string) string {
	return fmt.Sprintf("blacklist:access:%s", jti)
}

func normalizeDeviceID(deviceID string) string {
	trimmed := strings.TrimSpace(deviceID)
	if trimmed == "" {
		return defaultDeviceID
	}
	// 对 device 标识做哈希，避免把原始 UA/设备信息直接作为 Redis key 暴露。
	sum := sha256.Sum256([]byte(trimmed))
	return hex.EncodeToString(sum[:16])
}

func newJTI() (string, error) {
	buf := make([]byte, 16)
	if _, err := rand.Read(buf); err != nil {
		return "", err
	}
	return hex.EncodeToString(buf), nil
}
