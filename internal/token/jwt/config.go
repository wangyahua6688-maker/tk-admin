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
	// 返回当前处理结果。
	return &Config{
		// 处理当前语句逻辑。
		SigningKey: "change-this-secret",
		// 处理当前语句逻辑。
		AccessExpire: 1 * time.Hour,
		// 处理当前语句逻辑。
		RefreshExpire: 7 * 24 * time.Hour,
		// 处理当前语句逻辑。
		Issuer: "go-admin-full",
	}
}
