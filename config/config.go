package config

import (
	"fmt"
	"strings"

	"github.com/spf13/viper"
)

// Config 定义后台服务的统一配置结构。
type Config struct {
	Server struct { // 服务监听分组。
		Host string `yaml:"host"` // HTTP 监听地址。
		Port int    `yaml:"port"` // HTTP 监听端口。
	} `yaml:"server"` // 服务监听配置。

	Database struct { // 数据库分组。
		DSN      string `yaml:"dsn"`       // MySQL 连接串。
		LogLevel string `yaml:"log_level"` // GORM 日志级别。
	} `yaml:"database"` // 数据库配置。

	Redis struct { // Redis 分组。
		Addr     string `yaml:"addr"`     // Redis 地址。
		Password string `yaml:"password"` // Redis 密码。
		DB       int    `yaml:"db"`       // Redis 分库编号。
	} `yaml:"redis"` // Redis 配置。

	JWT struct { // 鉴权分组。
		SigningKey               string `yaml:"signing_key"`                 // JWT 签名密钥。
		AccessExpire             int    `yaml:"access_expire"`               // AccessToken 过期秒数。
		RefreshExpire            int    `yaml:"refresh_expire"`              // RefreshToken 过期秒数。
		SessionIdleTimeout       int    `yaml:"session_idle_timeout"`        // 会话空闲超时秒数（空闲超时后 refresh 也失效）。
		AllowInsecureMemoryStore bool   `yaml:"allow_insecure_memory_store"` // 是否允许内存 token 存储。
	} `yaml:"jwt"` // JWT 鉴权配置。

	CORS struct { // 跨域分组。
		AllowedOrigins   []string `yaml:"allowed_origins"`   // 允许跨域来源。
		AllowCredentials bool     `yaml:"allow_credentials"` // 是否允许携带凭据。
	} `yaml:"cors"` // CORS 配置。

	Auth struct { // 认证分组。
		AllowPublicRegister bool `yaml:"allow_public_register"` // 是否允许公开注册。
		Cookie              struct {
			AccessTokenName  string `yaml:"access_token_name"`  // AccessToken Cookie 名称。
			RefreshTokenName string `yaml:"refresh_token_name"` // RefreshToken Cookie 名称。
			Domain           string `yaml:"domain"`             // Cookie 域名。
			Path             string `yaml:"path"`               // Cookie 路径。
			Secure           bool   `yaml:"secure"`             // 是否仅 HTTPS 发送。
			HTTPOnly         bool   `yaml:"http_only"`          // 是否开启 HttpOnly。
			SameSite         string `yaml:"same_site"`          // SameSite 策略（lax/strict/none）。
		} `yaml:"cookie"` // 认证 Cookie 配置。
	} `yaml:"auth"` // 认证配置。

	Log struct { // 日志分组。
		Level    string `yaml:"level"`     // 日志级别。
		FilePath string `yaml:"file_path"` // 日志文件路径。
	} `yaml:"log"` // 日志配置。

	Upload struct { // 上传分组。
		SavePath string `yaml:"save_path"` // 上传文件保存目录。
		BaseURL  string `yaml:"base_url"`  // 上传资源访问前缀。
	} `yaml:"upload"` // 上传配置。
}

var cfg Config // 全局配置实例。

// Init 初始化并加载配置。
func Init(path string) {
	viper.SetConfigFile(path)                              // 指定配置文件路径。
	viper.SetConfigType("yaml")                            // 指定配置格式。
	viper.AddConfigPath(".")                               // 增加当前目录作为配置搜索路径。
	viper.SetEnvPrefix("godash")                           // 指定环境变量前缀。
	viper.AutomaticEnv()                                   // 启用环境变量覆盖。
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_")) // 支持配置键与环境变量键映射。

	viper.SetDefault("server.host", "0.0.0.0")                                                           // 默认监听地址。
	viper.SetDefault("server.port", 8080)                                                                // 默认监听端口。
	viper.SetDefault("jwt.signing_key", "")                                                              // 默认签名密钥（需外部注入）。
	viper.SetDefault("jwt.access_expire", 3600)                                                          // 默认 access token 1 小时。
	viper.SetDefault("jwt.refresh_expire", 604800)                                                       // 默认 refresh token 7 天。
	viper.SetDefault("jwt.session_idle_timeout", 3600)                                                   // 默认会话空闲超时 1 小时。
	viper.SetDefault("jwt.allow_insecure_memory_store", false)                                           // 默认禁止不安全内存存储。
	viper.SetDefault("database.dsn", "")                                                                 // 默认数据库连接串为空。
	viper.SetDefault("database.log_level", "warn")                                                       // 默认数据库日志级别。
	viper.SetDefault("redis.addr", "")                                                                   // 默认 Redis 地址为空。
	viper.SetDefault("redis.password", "")                                                               // 默认 Redis 密码为空。
	viper.SetDefault("redis.db", 0)                                                                      // 默认 Redis DB 0。
	viper.SetDefault("cors.allowed_origins", []string{"http://localhost:5173", "http://localhost:3000"}) // 默认允许本地前端来源。
	viper.SetDefault("cors.allow_credentials", true)                                                     // 默认允许跨域凭据（用于 Cookie 鉴权）。
	viper.SetDefault("auth.allow_public_register", false)                                                // 默认关闭公开注册。
	viper.SetDefault("auth.cookie.access_token_name", "tk_admin_access_token")                           // 默认 access token cookie 名称。
	viper.SetDefault("auth.cookie.refresh_token_name", "tk_admin_refresh_token")                         // 默认 refresh token cookie 名称。
	viper.SetDefault("auth.cookie.domain", "")                                                           // 默认不限制域名。
	viper.SetDefault("auth.cookie.path", "/")                                                            // 默认全路径可用。
	viper.SetDefault("auth.cookie.secure", false)                                                        // 默认开发环境允许 HTTP。
	viper.SetDefault("auth.cookie.http_only", true)                                                      // 默认启用 HttpOnly。
	viper.SetDefault("auth.cookie.same_site", "lax")                                                     // 默认 SameSite=Lax。
	viper.SetDefault("log.level", "info")                                                                // 默认日志级别 info。
	viper.SetDefault("log.file_path", "")                                                                // 默认输出到标准输出。
	viper.SetDefault("upload.save_path", "uploads")                                                      // 默认上传目录。
	viper.SetDefault("upload.base_url", "/uploads")                                                      // 默认上传访问前缀。

	if err := viper.ReadInConfig(); err != nil { // 尝试读取配置文件。
		fmt.Printf("未找到配置文件，使用默认配置: %v\n", err) // 未读取到配置时打印提示。
	} else { // 配置文件读取成功。
		fmt.Println("加载配置文件:", viper.ConfigFileUsed()) // 输出实际加载的配置文件路径。
	}

	cfg.Server.Host = viper.GetString("server.host")        // 读取监听地址。
	cfg.Server.Port = viper.GetInt("server.port")           // 读取监听端口。
	cfg.JWT.SigningKey = viper.GetString("jwt.signing_key") // 读取签名密钥。
	if cfg.JWT.SigningKey == "" {                           // 兼容旧配置键。
		cfg.JWT.SigningKey = viper.GetString("jwt.secret") // 回退读取旧签名键。
	}

	cfg.JWT.AccessExpire = viper.GetInt("jwt.access_expire") // 读取 access 过期时间。
	if cfg.JWT.AccessExpire <= 0 {                           // 兼容旧小时制配置。
		oldHour := viper.GetInt("jwt.expireHour") // 读取旧配置小时值。
		if oldHour > 0 {                          // 旧值有效才转换。
			cfg.JWT.AccessExpire = oldHour * 3600 // 小时转换为秒。
		}
	}

	cfg.JWT.RefreshExpire = viper.GetInt("jwt.refresh_expire") // 读取 refresh 过期时间。
	if cfg.JWT.RefreshExpire <= 0 {                            // 兼容旧小时制配置。
		oldHour := viper.GetInt("jwt.refreshHour") // 读取旧 refresh 小时值。
		if oldHour > 0 {                           // 旧值有效才转换。
			cfg.JWT.RefreshExpire = oldHour * 3600 // 小时转换为秒。
		}
	}
	cfg.JWT.SessionIdleTimeout = viper.GetInt("jwt.session_idle_timeout") // 读取会话空闲超时。
	if cfg.JWT.SessionIdleTimeout <= 0 {                                  // 兼容旧配置缺失时兜底 1 小时。
		cfg.JWT.SessionIdleTimeout = 3600 // 默认 1 小时空闲超时。
	}

	cfg.JWT.AllowInsecureMemoryStore = viper.GetBool("jwt.allow_insecure_memory_store")  // 读取内存存储开关。
	cfg.Database.DSN = viper.GetString("database.dsn")                                   // 读取数据库连接串。
	cfg.Database.LogLevel = viper.GetString("database.log_level")                        // 读取数据库日志级别。
	cfg.Redis.Addr = viper.GetString("redis.addr")                                       // 读取 Redis 地址。
	cfg.Redis.Password = viper.GetString("redis.password")                               // 读取 Redis 密码。
	cfg.Redis.DB = viper.GetInt("redis.db")                                              // 读取 Redis 分库。
	cfg.CORS.AllowedOrigins = viper.GetStringSlice("cors.allowed_origins")               // 读取跨域来源列表。
	cfg.CORS.AllowCredentials = viper.GetBool("cors.allow_credentials")                  // 读取跨域凭据开关。
	cfg.Auth.AllowPublicRegister = viper.GetBool("auth.allow_public_register")           // 读取公开注册开关。
	cfg.Auth.Cookie.AccessTokenName = viper.GetString("auth.cookie.access_token_name")   // 读取 access token cookie 名称。
	cfg.Auth.Cookie.RefreshTokenName = viper.GetString("auth.cookie.refresh_token_name") // 读取 refresh token cookie 名称。
	cfg.Auth.Cookie.Domain = viper.GetString("auth.cookie.domain")                       // 读取 cookie 域名。
	cfg.Auth.Cookie.Path = viper.GetString("auth.cookie.path")                           // 读取 cookie 路径。
	cfg.Auth.Cookie.Secure = viper.GetBool("auth.cookie.secure")                         // 读取 secure 开关。
	cfg.Auth.Cookie.HTTPOnly = viper.GetBool("auth.cookie.http_only")                    // 读取 httpOnly 开关。
	cfg.Auth.Cookie.SameSite = viper.GetString("auth.cookie.same_site")                  // 读取 SameSite 策略。
	cfg.Log.Level = viper.GetString("log.level")                                         // 读取日志级别。
	cfg.Log.FilePath = viper.GetString("log.file_path")                                  // 读取日志文件路径。
	cfg.Upload.SavePath = viper.GetString("upload.save_path")                            // 读取上传保存目录。
	cfg.Upload.BaseURL = viper.GetString("upload.base_url")                              // 读取上传访问前缀。
}

// GetConfig 返回已初始化配置。
func GetConfig() Config {
	return cfg // 返回全局配置快照。
}
