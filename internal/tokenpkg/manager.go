package tokenpkg

import (
	"fmt"
	"go-admin-full/internal/constants"
	"time"

	"github.com/golang-jwt/jwt/v5"
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
	now := time.Now()

	accessClaims := Claims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    m.Config.Issuer,
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(m.Config.AccessExpire)),
		},
	}
	access, err = m.signClaims(accessClaims)
	if err != nil {
		return "", "", fmt.Errorf("%w: %v", constants.ErrSigningToken, err)
	}

	refreshClaims := Claims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    m.Config.Issuer,
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(m.Config.RefreshExpire)),
		},
	}
	refresh, err = m.signClaims(refreshClaims)
	if err != nil {
		return "", "", fmt.Errorf("%w: %v", constants.ErrSigningToken, err)
	}

	if m.Store != nil {
		key := m.getRefreshTokenKey(userID)
		if serr := m.Store.Set(key, refresh, m.Config.RefreshExpire); serr != nil {
			return "", "", constants.ErrTokenStoreFailed
		}
	}
	return access, refresh, nil
}

// RefreshToken validates the given refresh token and issues a new access token.
// This implementation does not rotate refresh tokens; rotation can be added if desired.
func (m *Manager) RefreshToken(refreshToken string) (string, error) {
	uid, err := ParseToken(refreshToken, m.Config.SigningKey)
	if err != nil {
		return "", err
	}
	if m.Store != nil {
		key := m.getRefreshTokenKey(uid)
		stored, err := m.Store.Get(key)
		if err != nil {
			return "", constants.ErrTokenNotFound
		}
		if stored != refreshToken {
			return "", constants.ErrInvalidToken
		}
	}

	now := time.Now()
	accessClaims := Claims{
		UserID: uid,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    m.Config.Issuer,
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(m.Config.AccessExpire)),
		},
	}
	return m.signClaims(accessClaims)
}

// InvalidateRefresh deletes stored refresh token for a user (used on logout).
func (m *Manager) InvalidateRefresh(userID uint) error {
	if m.Store == nil {
		return nil
	}
	return m.Store.Delete(m.getRefreshTokenKey(userID))
}

func (m *Manager) signClaims(c Claims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, c)
	return token.SignedString([]byte(m.Config.SigningKey))
}

func (m *Manager) getRefreshTokenKey(userID uint) string {
	return fmt.Sprintf("refresh_token:%d", userID)
}
