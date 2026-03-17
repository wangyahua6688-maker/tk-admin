package config

import (
	"fmt"
	"strings"

	"github.com/spf13/viper"
)

// config/config.go
type Config struct {
	Server struct {
		Host string `yaml:"host"`
		Port int    `yaml:"port"`
	} `yaml:"server"`

	Database struct {
		DSN      string `yaml:"dsn"`
		LogLevel string `yaml:"log_level"`
	} `yaml:"database"`

	Redis struct {
		Addr     string `yaml:"addr"`
		Password string `yaml:"password"`
		DB       int    `yaml:"db"`
	} `yaml:"redis"`

	JWT struct {
		SigningKey               string `yaml:"signing_key"`
		AccessExpire             int    `yaml:"access_expire"`
		RefreshExpire            int    `yaml:"refresh_expire"`
		AllowInsecureMemoryStore bool   `yaml:"allow_insecure_memory_store"`
	} `yaml:"jwt"`

	CORS struct {
		AllowedOrigins   []string `yaml:"allowed_origins"`
		AllowCredentials bool     `yaml:"allow_credentials"`
	} `yaml:"cors"`

	Auth struct {
		AllowPublicRegister bool `yaml:"allow_public_register"`
	} `yaml:"auth"`

	Log struct {
		Level    string `yaml:"level"`
		FilePath string `yaml:"file_path"`
	} `yaml:"log"`

	Upload struct {
		SavePath string `yaml:"save_path"`
		BaseURL  string `yaml:"base_url"`
	} `yaml:"upload"`
}

var cfg Config

func Init(path string) {
	viper.SetConfigFile(path)
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.SetEnvPrefix("godash")
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	viper.SetDefault("server.host", "0.0.0.0")
	viper.SetDefault("server.port", 8080)
	viper.SetDefault("jwt.signing_key", "")
	viper.SetDefault("jwt.access_expire", 900)
	viper.SetDefault("jwt.refresh_expire", 604800)
	viper.SetDefault("jwt.allow_insecure_memory_store", false)
	viper.SetDefault("database.dsn", "")
	viper.SetDefault("database.log_level", "warn")
	viper.SetDefault("redis.addr", "")
	viper.SetDefault("redis.password", "")
	viper.SetDefault("redis.db", 0)
	viper.SetDefault("cors.allowed_origins", []string{"http://localhost:3000"})
	viper.SetDefault("cors.allow_credentials", false)
	viper.SetDefault("auth.allow_public_register", false)
	viper.SetDefault("log.level", "info")
	viper.SetDefault("log.file_path", "")
	viper.SetDefault("upload.save_path", "uploads")
	viper.SetDefault("upload.base_url", "/uploads")

	if err := viper.ReadInConfig(); err != nil {
		fmt.Printf("未找到配置文件，使用默认配置: %v\n", err)
	} else {
		fmt.Println("加载配置文件:", viper.ConfigFileUsed())
	}

	cfg.Server.Host = viper.GetString("server.host")
	cfg.Server.Port = viper.GetInt("server.port")

	// Backward compatibility: old keys jwt.secret / jwt.expireHour / jwt.refreshHour
	cfg.JWT.SigningKey = viper.GetString("jwt.signing_key")
	if cfg.JWT.SigningKey == "" {
		cfg.JWT.SigningKey = viper.GetString("jwt.secret")
	}

	cfg.JWT.AccessExpire = viper.GetInt("jwt.access_expire")
	if cfg.JWT.AccessExpire <= 0 {
		oldHour := viper.GetInt("jwt.expireHour")
		if oldHour > 0 {
			cfg.JWT.AccessExpire = oldHour * 3600
		}
	}

	cfg.JWT.RefreshExpire = viper.GetInt("jwt.refresh_expire")
	if cfg.JWT.RefreshExpire <= 0 {
		oldHour := viper.GetInt("jwt.refreshHour")
		if oldHour > 0 {
			cfg.JWT.RefreshExpire = oldHour * 3600
		}
	}
	cfg.JWT.AllowInsecureMemoryStore = viper.GetBool("jwt.allow_insecure_memory_store")

	cfg.Database.DSN = viper.GetString("database.dsn")
	cfg.Database.LogLevel = viper.GetString("database.log_level")
	cfg.Redis.Addr = viper.GetString("redis.addr")
	cfg.Redis.Password = viper.GetString("redis.password")
	cfg.Redis.DB = viper.GetInt("redis.db")
	cfg.CORS.AllowedOrigins = viper.GetStringSlice("cors.allowed_origins")
	cfg.CORS.AllowCredentials = viper.GetBool("cors.allow_credentials")
	cfg.Auth.AllowPublicRegister = viper.GetBool("auth.allow_public_register")
	cfg.Log.Level = viper.GetString("log.level")
	cfg.Log.FilePath = viper.GetString("log.file_path")
	cfg.Upload.SavePath = viper.GetString("upload.save_path")
	cfg.Upload.BaseURL = viper.GetString("upload.base_url")
}

func GetConfig() Config {
	return cfg
}
