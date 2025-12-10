package server

import (
	"context"
	"fmt"
	"strings"
	"time"

	"go-admin-full/config"
	"go-admin-full/internal/models"
	"go-admin-full/internal/tokenpkg"
	"go-admin-full/internal/utils"

	"github.com/go-redis/redis/v8"
)

func Run(cfg config.Config) error {
	// 1. 初始化全局日志记录器
	logCfg := utils.DefaultLogConfig()
	logCfg.Level = utils.LogLevelFromString(cfg.Log.Level)
	logCfg.FilePath = cfg.Log.FilePath

	if err := utils.InitGlobalLogger(logCfg); err != nil {
		return fmt.Errorf("failed to init logger: %w", err)
	}
	defer utils.GetLogger().Close()

	logger := utils.GetLogger()
	logger.Info("Starting application...")

	// 2. JWT配置安全校验（fail fast）
	signingKey := strings.TrimSpace(cfg.JWT.SigningKey)
	if len(signingKey) < 32 || signingKey == "change-this-secret" || strings.Contains(signingKey, "replace-with") {
		return fmt.Errorf("invalid jwt.signing_key: require at least 32 chars and non-placeholder value")
	}

	if cfg.JWT.AccessExpire <= 0 || cfg.JWT.RefreshExpire <= 0 {
		return fmt.Errorf("invalid jwt expiration config: access_expire and refresh_expire must be positive seconds")
	}

	// 3. 初始化数据库连接
	dbCfg := utils.DefaultDBConfig()
	dbCfg.DSN = cfg.Database.DSN
	dbCfg.LogLevel = utils.GormLogLevelFromString(cfg.Database.LogLevel)

	db, err := utils.NewMySQLDB(dbCfg)
	if err != nil {
		logger.Fatal("Failed to connect to database: %v", err)
	}

	// 4. 初始化Redis连接（可选）
	var redisClient *redis.Client
	if cfg.Redis.Addr != "" {
		redisCfg := utils.DefaultRedisConfig()
		redisCfg.Addr = cfg.Redis.Addr
		redisCfg.Password = cfg.Redis.Password
		redisCfg.DB = cfg.Redis.DB

		redisClient, err = utils.NewRedisClient(redisCfg)
		if err != nil {
			if !cfg.JWT.AllowInsecureMemoryStore {
				return fmt.Errorf("failed to connect to redis and insecure memory token store is disabled: %w", err)
			}
			logger.Warn("Failed to connect to Redis: %v, falling back to memory token store", err)
		} else {
			logger.Info("Redis connected successfully")
		}
	}

	// 5. 数据库迁移
	ctx := context.Background()
	if err := db.WithContext(ctx).AutoMigrate(
		&models.User{}, &models.Role{}, &models.Permission{}, &models.Menu{},
		&models.RolePermission{},
		&models.LoginLog{}, &models.RefreshTokenRecord{},
	); err != nil {
		logger.Fatal("Failed to migrate database: %v", err)
	}

	// 6. 创建token存储
	var store tokenpkg.Store
	if redisClient != nil {
		// 创建Redis存储适配器
		store = tokenpkg.NewRedisStoreWithClient(redisClient)
	} else {
		if !cfg.JWT.AllowInsecureMemoryStore {
			return fmt.Errorf("未配置redis，且jwt.allow_insecure_memory_store=false")
		}
		// 使用内存存储
		store = tokenpkg.NewMemoryStore()
		logger.Warn("Using in-memory token store (not suitable for production)")
	}

	// 7. 配置token管理器
	jwtCfg := tokenpkg.DefaultConfig()
	jwtCfg.SigningKey = signingKey
	jwtCfg.AccessExpire = time.Duration(cfg.JWT.AccessExpire) * time.Second
	jwtCfg.RefreshExpire = time.Duration(cfg.JWT.RefreshExpire) * time.Second

	mgr := tokenpkg.NewManager(jwtCfg, store)

	// 8. 设置路由
	r := SetupRouter(mgr, db, redisClient, logger, cfg)

	addr := fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port)
	logger.Info("Server starting on %s", addr)

	return r.Run(addr)
}
