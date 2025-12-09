package config

import (
	"fmt"
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
		SigningKey    string `yaml:"signing_key"`
		AccessExpire  int    `yaml:"access_expire"`
		RefreshExpire int    `yaml:"refresh_expire"`
	} `yaml:"jwt"`

	Log struct {
		Level    string `yaml:"level"`
		FilePath string `yaml:"file_path"`
	} `yaml:"log"`
}

var cfg Config

func Init(path string) {
	viper.SetConfigFile(path)
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.SetEnvPrefix("godash")
	viper.AutomaticEnv()

	viper.SetDefault("server.host", "0.0.0.0")
	viper.SetDefault("server.port", 8080)
	viper.SetDefault("jwt.secret", "replace-with-strong-secret")
	viper.SetDefault("jwt.expireHour", 1)
	viper.SetDefault("jwt.refreshHour", 168)
	viper.SetDefault("database.dialect", "sqlite")
	viper.SetDefault("database.dsn", "data.db")
	viper.SetDefault("redis.addr", "")
	viper.SetDefault("redis.password", "")
	viper.SetDefault("redis.db", 0)

	if err := viper.ReadInConfig(); err != nil {
		fmt.Printf("未找到配置文件，使用默认配置: %v\n", err)
	} else {
		fmt.Println("加载配置文件:", viper.ConfigFileUsed())
	}

	cfg.Server.Host = viper.GetString("server.host")
	cfg.Server.Port = viper.GetInt("server.port")
	cfg.JWT.SigningKey = viper.GetString("jwt.signing_key")
	cfg.JWT.AccessExpire = viper.GetInt("jwt.access_expire")
	cfg.JWT.RefreshExpire = viper.GetInt("jwt.refresh_expire")
	cfg.Database.DSN = viper.GetString("database.dsn")
	cfg.Redis.Addr = viper.GetString("redis.addr")
	cfg.Redis.Password = viper.GetString("redis.password")
	cfg.Redis.DB = viper.GetInt("redis.db")
}

func GetConfig() Config {
	return cfg
}
