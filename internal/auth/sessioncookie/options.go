package sessioncookie

import (
	"net/http"
	"strings"

	"go-admin/config"
)

// 默认 cookie 名称，避免配置缺失时退化到不安全或不兼容行为。
const (
	DefaultAccessTokenCookieName  = "tk_admin_access_token"
	DefaultRefreshTokenCookieName = "tk_admin_refresh_token"
)

// Options 统一承载认证 cookie 配置。
type Options struct {
	AccessTokenName  string
	RefreshTokenName string
	Domain           string
	Path             string
	Secure           bool
	HTTPOnly         bool
	SameSite         http.SameSite
}

// FromConfig 根据全局配置生成 cookie 选项并补齐安全默认值。
func FromConfig(cfg config.Config) Options {
	opt := Options{
		AccessTokenName:  strings.TrimSpace(cfg.Auth.Cookie.AccessTokenName),
		RefreshTokenName: strings.TrimSpace(cfg.Auth.Cookie.RefreshTokenName),
		Domain:           strings.TrimSpace(cfg.Auth.Cookie.Domain),
		Path:             strings.TrimSpace(cfg.Auth.Cookie.Path),
		Secure:           cfg.Auth.Cookie.Secure,
		HTTPOnly:         cfg.Auth.Cookie.HTTPOnly,
		SameSite:         parseSameSite(cfg.Auth.Cookie.SameSite),
	}

	if opt.AccessTokenName == "" {
		opt.AccessTokenName = DefaultAccessTokenCookieName
	}
	if opt.RefreshTokenName == "" {
		opt.RefreshTokenName = DefaultRefreshTokenCookieName
	}
	if opt.Path == "" {
		opt.Path = "/"
	}
	// 浏览器要求 SameSite=None 必须配合 Secure。
	if opt.SameSite == http.SameSiteNoneMode {
		opt.Secure = true
	}

	return opt
}

// parseSameSite 将字符串配置映射到 net/http SameSite 枚举。
func parseSameSite(raw string) http.SameSite {
	switch strings.ToLower(strings.TrimSpace(raw)) {
	case "strict":
		return http.SameSiteStrictMode
	case "none":
		return http.SameSiteNoneMode
	case "default":
		return http.SameSiteDefaultMode
	case "lax", "":
		return http.SameSiteLaxMode
	default:
		return http.SameSiteLaxMode
	}
}
