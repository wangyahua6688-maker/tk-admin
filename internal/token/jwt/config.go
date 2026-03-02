package jwt

import "time"

// Config holds token related configuration.
type Config struct {
	SigningKey    string        // JWT 签名密钥（请从 env/config 注入）
	AccessExpire  time.Duration // access token 过期时间
	RefreshExpire time.Duration // refresh token 过期时间
	Issuer        string        // issuer claim
}

// DefaultConfig returns a safe default (must replace SigningKey in production).
func DefaultConfig() *Config {
	return &Config{
		SigningKey:    "change-this-secret",
		AccessExpire:  15 * time.Minute,
		RefreshExpire: 7 * 24 * time.Hour,
		Issuer:        "go-admin-full",
	}
}
