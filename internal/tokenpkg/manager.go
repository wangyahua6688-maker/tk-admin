package tokenpkg

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"go-admin-full/internal/constants"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

const (
	TokenTypeAccess  = "access"
	TokenTypeRefresh = "refresh"
	defaultDeviceID  = "default"
)

// Manager handles generation, refresh and invalidation of tokens.
type Manager struct {
	Config *Config
	Store  Store
}

// NewManager constructs a token manager with config and optional store.
func NewManager(cfg *Config, store Store) *Manager {
	return &Manager{
		Config: cfg,
		Store:  store,
	}
}

// GenerateTokens returns signed access and refresh tokens for a user id.
// It also persists the refresh token into Store if Store != nil.
func (m *Manager) GenerateTokens(userID uint) (access string, refresh string, err error) {
	return m.GenerateTokensWithDevice(userID, defaultDeviceID)
}

// GenerateTokensWithDevice returns signed access and refresh tokens scoped to a device.
func (m *Manager) GenerateTokensWithDevice(userID uint, deviceID string) (access string, refresh string, err error) {
	now := time.Now()
	deviceID = normalizeDeviceID(deviceID)

	accessJTI, err := newJTI()
	if err != nil {
		return "", "", fmt.Errorf("%w: %v", constants.ErrSigningToken, err)
	}

	accessClaims := Claims{
		UserID:    userID,
		TokenType: TokenTypeAccess,
		DeviceID:  deviceID,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    m.Config.Issuer,
			ID:        accessJTI,
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(m.Config.AccessExpire)),
		},
	}
	access, err = m.signClaims(accessClaims)
	if err != nil {
		return "", "", fmt.Errorf("%w: %v", constants.ErrSigningToken, err)
	}

	refreshJTI, err := newJTI()
	if err != nil {
		return "", "", fmt.Errorf("%w: %v", constants.ErrSigningToken, err)
	}

	refreshClaims := Claims{
		UserID:    userID,
		TokenType: TokenTypeRefresh,
		DeviceID:  deviceID,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    m.Config.Issuer,
			ID:        refreshJTI,
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(m.Config.RefreshExpire)),
		},
	}
	refresh, err = m.signClaims(refreshClaims)
	if err != nil {
		return "", "", fmt.Errorf("%w: %v", constants.ErrSigningToken, err)
	}

	if m.Store != nil {
		key := m.getRefreshTokenKey(userID, deviceID)
		if serr := m.Store.Set(key, refresh, m.Config.RefreshExpire); serr != nil {
			return "", "", constants.ErrTokenStoreFailed
		}
	}
	return access, refresh, nil
}

// RefreshToken validates the given refresh token and issues a new access token.
func (m *Manager) RefreshToken(refreshToken string) (string, error) {
	access, _, err := m.RefreshTokenPair(refreshToken)
	return access, err
}

// RefreshTokenPair validates and rotates refresh token, then returns new access+refresh token pair.
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

// ValidateAccessToken validates access token and checks token blacklist.
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

// RevokeAccessToken writes access token jti into blacklist until token expires.
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

// InvalidateRefresh deletes stored refresh token for a user (used on logout).
func (m *Manager) InvalidateRefresh(userID uint, deviceID string) error {
	if m.Store == nil {
		return nil
	}
	return m.Store.Delete(m.getRefreshTokenKey(userID, normalizeDeviceID(deviceID)))
}

func (m *Manager) signClaims(c Claims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, c)
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
