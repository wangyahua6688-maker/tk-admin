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
	// 返回当前处理结果。
	return DBConfig{
		// 处理当前语句逻辑。
		MaxIdleConns: 10,
		// 处理当前语句逻辑。
		MaxOpenConns: 100,
		// 处理当前语句逻辑。
		ConnMaxLifetime: time.Hour,
		// 处理当前语句逻辑。
		LogLevel: logger.Warn,
	}
}

// NewMySQLDB 初始化MySQL数据库连接
func NewMySQLDB(cfg DBConfig) (*gorm.DB, error) {
	// 判断条件并进入对应分支逻辑。
	if cfg.DSN == "" {
		// 返回当前处理结果。
		return nil, fmt.Errorf("database DSN is empty")
	}

	// 定义并初始化当前变量。
	gormConfig := &gorm.Config{
		// 进入新的代码块进行处理。
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true, // 使用单数表名
		},
		// 调用logger.Default.LogMode完成当前处理。
		Logger: logger.Default.LogMode(cfg.LogLevel),
		// 调用func完成当前处理。
		NowFunc: func() time.Time {
			return time.Now().UTC() // 使用UTC时间
		},
	}

	// 定义并初始化当前变量。
	db, err := gorm.Open(mysql.Open(cfg.DSN), gormConfig)
	// 判断条件并进入对应分支逻辑。
	if err != nil {
		// 返回当前处理结果。
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	// 定义并初始化当前变量。
	sqlDB, err := db.DB()
	// 判断条件并进入对应分支逻辑。
	if err != nil {
		// 返回当前处理结果。
		return nil, fmt.Errorf("failed to get database instance: %w", err)
	}

	// 设置连接池
	sqlDB.SetMaxIdleConns(cfg.MaxIdleConns)
	// 调用sqlDB.SetMaxOpenConns完成当前处理。
	sqlDB.SetMaxOpenConns(cfg.MaxOpenConns)
	// 调用sqlDB.SetConnMaxLifetime完成当前处理。
	sqlDB.SetConnMaxLifetime(cfg.ConnMaxLifetime)

	// 测试连接
	if err := sqlDB.Ping(); err != nil {
		// 返回当前处理结果。
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	// 调用log.Println完成当前处理。
	log.Println("MySQL database connected successfully")
	// 返回当前处理结果。
	return db, nil
}

// DBFromContext 从上下文中获取数据库实例（需提前设置）
func DBFromContext(ctx context.Context) *gorm.DB {
	// 判断条件并进入对应分支逻辑。
	if db, ok := ctx.Value("db").(*gorm.DB); ok {
		// 返回当前处理结果。
		return db
	}
	// 返回当前处理结果。
	return nil
}

// GormLogLevelFromString 将字符串转换为gorm日志级别
func GormLogLevelFromString(level string) logger.LogLevel {
	// 根据表达式进入多分支处理。
	switch level {
	case "silent":
		// 返回当前处理结果。
		return logger.Silent
	case "error":
		// 返回当前处理结果。
		return logger.Error
	case "warn":
		// 返回当前处理结果。
		return logger.Warn
	case "info":
		// 返回当前处理结果。
		return logger.Info
	default:
		// 返回当前处理结果。
		return logger.Warn
	}
}
