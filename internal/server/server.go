package server

import (
	"context"
	"fmt"
	"strings"
	"time"

	"go-admin-full/config"
	"go-admin-full/internal/models"
	tokenjwt "go-admin-full/internal/token/jwt"
	tokenstore "go-admin-full/internal/token/store"
	"go-admin-full/internal/utils"

	"github.com/go-redis/redis/v8"
)

func Run(cfg config.Config) error {
	// 1) 初始化全局日志记录器。
	// 说明：logger 在整个进程生命周期内复用，便于统一输出格式与等级控制。
	logCfg := utils.DefaultLogConfig()
	logCfg.Level = utils.LogLevelFromString(cfg.Log.Level)
	logCfg.FilePath = cfg.Log.FilePath

	if err := utils.InitGlobalLogger(logCfg); err != nil {
		return fmt.Errorf("failed to init logger: %w", err)
	}
	defer utils.GetLogger().Close()

	logger := utils.GetLogger()
	logger.Info("Starting application...")

	// 2) JWT 配置安全校验（fail fast）。
	// 说明：启动阶段直接拦截弱密钥/无效过期时间，避免服务“带病运行”。
	signingKey := strings.TrimSpace(cfg.JWT.SigningKey)
	if len(signingKey) < 32 || signingKey == "change-this-secret" || strings.Contains(signingKey, "replace-with") {
		return fmt.Errorf("invalid jwt.signing_key: require at least 32 chars and non-placeholder value")
	}

	if cfg.JWT.AccessExpire <= 0 || cfg.JWT.RefreshExpire <= 0 {
		return fmt.Errorf("invalid jwt expiration config: access_expire and refresh_expire must be positive seconds")
	}

	// 3) 初始化数据库连接。
	dbCfg := utils.DefaultDBConfig()
	dbCfg.DSN = cfg.Database.DSN
	dbCfg.LogLevel = utils.GormLogLevelFromString(cfg.Database.LogLevel)

	db, err := utils.NewMySQLDB(dbCfg)
	if err != nil {
		logger.Fatal("Failed to connect to database: %v", err)
	}

	// 4) 初始化 Redis 连接（可选）。
	// 说明：
	// - 配置了 redis.addr 时优先使用 Redis 作为 token 存储；
	// - Redis 不可用时，是否允许降级到内存存储由 jwt.allow_insecure_memory_store 决定。
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

	// 5) 数据库结构迁移（AutoMigrate）。
	// 关键说明：
	// - 这里只做“结构对齐”，不会执行 database/mysql/001_rbac_schema_and_seed.sql；
	// - 也不会自动写入角色/权限/菜单/admin 等初始化数据；
	// - 生产环境建议配合版本化 SQL 迁移工具，而不是仅依赖 AutoMigrate。
	ctx := context.Background()
	if err := db.WithContext(ctx).AutoMigrate(
		&models.User{}, &models.Role{}, &models.Permission{}, &models.Menu{},
		&models.RolePermission{},
		&models.LoginLog{}, &models.RefreshTokenRecord{},
		&models.SystemMessage{},
	); err != nil {
		logger.Fatal("Failed to migrate database: %v", err)
	}

	// 6) 创建 token 存储实现（Store）。
	// - RedisStore：支持共享状态，适合多实例部署；
	// - MemoryStore：仅本进程可见，仅建议本地开发使用。
	var store tokenstore.Store
	if redisClient != nil {
		// 创建Redis存储适配器
		store = tokenstore.NewRedisStoreWithClient(redisClient)
	} else {
		if !cfg.JWT.AllowInsecureMemoryStore {
			return fmt.Errorf("未配置redis，且jwt.allow_insecure_memory_store=false")
		}
		// 使用内存存储
		store = tokenstore.NewMemoryStore()
		logger.Warn("Using in-memory token store (not suitable for production)")
	}

	// 7) 配置 Token 管理器（签发 / 校验 / 刷新 / 撤销都在这里完成）。
	jwtCfg := tokenjwt.DefaultConfig()
	jwtCfg.SigningKey = signingKey
	jwtCfg.AccessExpire = time.Duration(cfg.JWT.AccessExpire) * time.Second
	jwtCfg.RefreshExpire = time.Duration(cfg.JWT.RefreshExpire) * time.Second

	mgr := tokenjwt.NewManager(jwtCfg, store)

	// 8) 装配路由与中间件。
	r := SetupRouter(mgr, db, redisClient, logger, cfg)

	addr := fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port)
	logger.Info("Server starting on %s", addr)

	return r.Run(addr)
}
