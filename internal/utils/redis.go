package utils

import (
	"context"
	"time"

	"github.com/go-redis/redis/v8"
	redisx "tk-common/utils/redisx/v8"
)

// RedisConfig Redis配置
type RedisConfig struct {
	Addr         string        // 地址
	Password     string        // 密码
	DB           int           // 数据库
	PoolSize     int           // 连接池大小
	MinIdleConns int           // 最小空闲连接数
	DialTimeout  time.Duration // 连接超时
	ReadTimeout  time.Duration // 读取超时
	WriteTimeout time.Duration // 写入超时
}

// DefaultRedisConfig 默认Redis配置
func DefaultRedisConfig() RedisConfig {
	cfg := redisx.DefaultConfig()
	return RedisConfig{
		PoolSize:     cfg.PoolSize,
		MinIdleConns: cfg.MinIdleConns,
		DialTimeout:  cfg.DialTimeout,
		ReadTimeout:  cfg.ReadTimeout,
		WriteTimeout: cfg.WriteTimeout,
	}
}

// NewRedisClient 初始化Redis客户端
func NewRedisClient(cfg RedisConfig) (*redis.Client, error) {
	return redisx.NewClient(context.Background(), redisx.Config{
		Addr:         cfg.Addr,
		Password:     cfg.Password,
		DB:           cfg.DB,
		PoolSize:     cfg.PoolSize,
		MinIdleConns: cfg.MinIdleConns,
		DialTimeout:  cfg.DialTimeout,
		ReadTimeout:  cfg.ReadTimeout,
		WriteTimeout: cfg.WriteTimeout,
	})
}

// NewRedisClientWithContext 创建带上下文的Redis客户端
func NewRedisClientWithContext(ctx context.Context, cfg RedisConfig) (*redis.Client, error) {
	return redisx.NewClient(ctx, redisx.Config{
		Addr:         cfg.Addr,
		Password:     cfg.Password,
		DB:           cfg.DB,
		PoolSize:     cfg.PoolSize,
		MinIdleConns: cfg.MinIdleConns,
		DialTimeout:  cfg.DialTimeout,
		ReadTimeout:  cfg.ReadTimeout,
		WriteTimeout: cfg.WriteTimeout,
	})
}

// RedisFromContext 从上下文中获取Redis客户端（需提前设置）
func RedisFromContext(ctx context.Context) *redis.Client {
	return redisx.RedisFromContext(ctx, "redis")
}
