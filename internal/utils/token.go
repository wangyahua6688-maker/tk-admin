package utils

import (
	"errors"
	"go-admin-full/internal/tokenpkg"
)

// ParseToken 解析 token 并返回 userID（委托给 tokenpkg）
// 注意：这里示例使用 DefaultConfig 的 signing key；生产请把真实配置注入并避免使用 DefaultConfig 的默认 secret。
func ParseToken(tokenString string) (uint, error) {
	cfg := tokenpkg.DefaultConfig()
	if cfg.SigningKey == "" || cfg.SigningKey == "change-this-secret" {
		return 0, errors.New("jwt signing key is default or empty; set real key in config")
	}
	return tokenpkg.ParseToken(tokenString, cfg.SigningKey)
}

// GenerateTokens 是 utils 层的简单包装：使用你传入的 Manager 生成 tokens
// 注意：真正的实现逻辑在 tokenpkg.Manager.GenerateTokens 中。
func GenerateTokens(mgr *tokenpkg.Manager, userID uint) (access string, refresh string, err error) {
	if mgr == nil {
		return "", "", errors.New("token manager is nil")
	}
	return mgr.GenerateTokens(userID)
}
