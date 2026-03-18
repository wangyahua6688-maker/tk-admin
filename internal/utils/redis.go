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
	// 定义并初始化当前变量。
	cfg := redisx.DefaultConfig()
	// 返回当前处理结果。
	return RedisConfig{
		// 处理当前语句逻辑。
		PoolSize: cfg.PoolSize,
		// 处理当前语句逻辑。
		MinIdleConns: cfg.MinIdleConns,
		// 处理当前语句逻辑。
		DialTimeout: cfg.DialTimeout,
		// 处理当前语句逻辑。
		ReadTimeout: cfg.ReadTimeout,
		// 处理当前语句逻辑。
		WriteTimeout: cfg.WriteTimeout,
	}
}

// NewRedisClient 初始化Redis客户端
func NewRedisClient(cfg RedisConfig) (*redis.Client, error) {
	// 返回当前处理结果。
	return redisx.NewClient(context.Background(), redisx.Config{
		// 处理当前语句逻辑。
		Addr: cfg.Addr,
		// 处理当前语句逻辑。
		Password: cfg.Password,
		// 处理当前语句逻辑。
		DB: cfg.DB,
		// 处理当前语句逻辑。
		PoolSize: cfg.PoolSize,
		// 处理当前语句逻辑。
		MinIdleConns: cfg.MinIdleConns,
		// 处理当前语句逻辑。
		DialTimeout: cfg.DialTimeout,
		// 处理当前语句逻辑。
		ReadTimeout: cfg.ReadTimeout,
		// 处理当前语句逻辑。
		WriteTimeout: cfg.WriteTimeout,
	})
}
