package utils

import (
	"context"
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
	"log"
	"time"
)

// DBConfig 数据库配置
type DBConfig struct {
	DSN             string          // 连接字符串
	MaxIdleConns    int             // 最大空闲连接数
	MaxOpenConns    int             // 最大打开连接数
	ConnMaxLifetime time.Duration   // 连接最大生命周期
	LogLevel        logger.LogLevel // 日志级别
}

// DefaultDBConfig 默认数据库配置
func DefaultDBConfig() DBConfig {
	return DBConfig{
		MaxIdleConns:    10,
		MaxOpenConns:    100,
		ConnMaxLifetime: time.Hour,
		LogLevel:        logger.Warn,
	}
}

// NewMySQLDB 初始化MySQL数据库连接
func NewMySQLDB(cfg DBConfig) (*gorm.DB, error) {
	if cfg.DSN == "" {
		return nil, fmt.Errorf("database DSN is empty")
	}

	gormConfig := &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true, // 使用单数表名
		},
		Logger: logger.Default.LogMode(cfg.LogLevel),
		NowFunc: func() time.Time {
			return time.Now().UTC() // 使用UTC时间
		},
	}

	db, err := gorm.Open(mysql.Open(cfg.DSN), gormConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get database instance: %w", err)
	}

	// 设置连接池
	sqlDB.SetMaxIdleConns(cfg.MaxIdleConns)
	sqlDB.SetMaxOpenConns(cfg.MaxOpenConns)
	sqlDB.SetConnMaxLifetime(cfg.ConnMaxLifetime)

	// 测试连接
	if err := sqlDB.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	log.Println("MySQL database connected successfully")
	return db, nil
}

// NewMySQLDBWithContext 创建带上下文的数据库连接
func NewMySQLDBWithContext(ctx context.Context, cfg DBConfig) (*gorm.DB, error) {
	db, err := NewMySQLDB(cfg)
	if err != nil {
		return nil, err
	}

	// 设置数据库操作时的默认上下文
	return db.WithContext(ctx), nil
}

// DBFromContext 从上下文中获取数据库实例（需提前设置）
func DBFromContext(ctx context.Context) *gorm.DB {
	if db, ok := ctx.Value("db").(*gorm.DB); ok {
		return db
	}
	return nil
}

// GormLogLevelFromString 将字符串转换为gorm日志级别
func GormLogLevelFromString(level string) logger.LogLevel {
	switch level {
	case "silent":
		return logger.Silent
	case "error":
		return logger.Error
	case "warn":
		return logger.Warn
	case "info":
		return logger.Info
	default:
		return logger.Warn
	}
}
