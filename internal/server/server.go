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

	"github.com/go-redis/redis/v8"
	gormx "tk-common/utils/dbx/gormx"
	commonlogx "tk-common/utils/logx"
	redisx "tk-common/utils/redisx/v8"
)

// Run 处理Run相关逻辑。
func Run(cfg config.Config) error {
	// 1) 初始化全局日志记录器。
	// 说明：logger 在整个进程生命周期内复用，便于统一输出格式与等级控制。
	logCfg := commonlogx.DefaultLogConfig()
	// 更新当前变量或字段值。
	logCfg.Level = commonlogx.LogLevelFromString(cfg.Log.Level)
	// 更新当前变量或字段值。
	logCfg.FilePath = cfg.Log.FilePath

	// 判断条件并进入对应分支逻辑。
	if err := commonlogx.InitGlobalLogger(logCfg); err != nil {
		// 返回当前处理结果。
		return fmt.Errorf("failed to init logger: %w", err)
	}
	// 注册延迟执行逻辑。
	defer func(logger *commonlogx.Logger) {
		// 尝试关闭日志文件句柄，失败时打印到标准错误避免静默丢失错误。
		if err := logger.Close(); err != nil {
			fmt.Printf("close logger failed: %v\n", err)
		}
	}(commonlogx.GetLogger())

	// 定义并初始化当前变量。
	logger := commonlogx.GetLogger()
	// 调用logger.Info完成当前处理。
	logger.Info("Starting application...")

	// 2) JWT 配置安全校验（fail fast）。
	// 说明：启动阶段直接拦截弱密钥/无效过期时间，避免服务“带病运行”。
	signingKey := strings.TrimSpace(cfg.JWT.SigningKey)
	// 判断条件并进入对应分支逻辑。
	if len(signingKey) < 32 || signingKey == "change-this-secret" || strings.Contains(signingKey, "replace-with") {
		// 返回当前处理结果。
		return fmt.Errorf("invalid jwt.signing_key: require at least 32 chars and non-placeholder value")
	}

	// 判断条件并进入对应分支逻辑。
	if cfg.JWT.AccessExpire <= 0 || cfg.JWT.RefreshExpire <= 0 || cfg.JWT.SessionIdleTimeout <= 0 {
		// 返回当前处理结果。
		return fmt.Errorf("invalid jwt expiration config: access_expire, refresh_expire and session_idle_timeout must be positive seconds")
	}

	// 3) 初始化数据库连接。
	dbCfg := gormx.DefaultDBConfig()
	// 更新当前变量或字段值。
	dbCfg.DSN = cfg.Database.DSN
	// 更新当前变量或字段值。
	dbCfg.LogLevel = gormx.GormLogLevelFromString(cfg.Database.LogLevel)

	// 定义并初始化当前变量。
	db, err := gormx.NewMySQLDB(dbCfg)
	// 判断条件并进入对应分支逻辑。
	if err != nil {
		// 调用logger.Fatal完成当前处理。
		logger.Fatal("Failed to connect to database: %v", err)
	}

	// 4) 初始化 Redis 连接（可选）。
	// 说明：
	// - 配置了 redis.addr 时优先使用 Redis 作为 token 存储；
	// - Redis 不可用时，是否允许降级到内存存储由 jwt.allow_insecure_memory_store 决定。
	var redisClient *redis.Client
	// 判断条件并进入对应分支逻辑。
	if cfg.Redis.Addr != "" {
		// 定义并初始化当前变量。
		redisCfg := redisx.DefaultConfig()
		// 更新当前变量或字段值。
		redisCfg.Addr = cfg.Redis.Addr
		// 更新当前变量或字段值。
		redisCfg.Password = cfg.Redis.Password
		// 更新当前变量或字段值。
		redisCfg.DB = cfg.Redis.DB

		// 更新当前变量或字段值。
		redisClient, err = redisx.NewClient(context.Background(), redisCfg)
		// 判断条件并进入对应分支逻辑。
		if err != nil {
			// 判断条件并进入对应分支逻辑。
			if !cfg.JWT.AllowInsecureMemoryStore {
				// 返回当前处理结果。
				return fmt.Errorf("failed to connect to redis and insecure memory token store is disabled: %w", err)
			}
			// 调用logger.Warn完成当前处理。
			logger.Warn("Failed to connect to Redis: %v, falling back to memory token store", err)
			// 进入新的代码块进行处理。
		} else {
			// 调用logger.Info完成当前处理。
			logger.Info("Redis connected successfully")
		}
	}

	// 5) 数据库结构迁移（AutoMigrate）。
	// 关键说明：
	// - 这里只做“结构对齐”，不会执行 database/mysql/001_rbac_schema_and_seed.sql；
	// - 也不会自动写入角色/权限/菜单/admin 等初始化数据；
	// - 生产环境建议配合版本化 SQL 迁移工具，而不是仅依赖 AutoMigrate。
	ctx := context.Background()
	// 判断条件并进入对应分支逻辑。
	if err := db.WithContext(ctx).AutoMigrate(
		// 处理当前语句逻辑。
		&models.User{}, &models.Role{}, &models.Permission{}, &models.Menu{},
		// 处理当前语句逻辑。
		&models.RolePermission{},
		// 处理当前语句逻辑。
		&models.LoginLog{}, &models.RefreshTokenRecord{},
		// 处理当前语句逻辑。
		&models.SystemMessage{},
		// 处理当前语句逻辑。
		&models.WUser{},
		// 处理当前语句逻辑。
		&models.WBanner{},
		// 处理当前语句逻辑。
		&models.WBroadcast{},
		// 处理当前语句逻辑。
		&models.WSpecialLottery{},
		// 处理当前语句逻辑。
		&models.WLotteryCategory{},
		// 处理当前语句逻辑。
		&models.WLotteryInfo{},
		// 处理当前语句逻辑。
		&models.WLotteryOption{},
		// 处理当前语句逻辑。
		&models.WPostArticle{},
		// 处理当前语句逻辑。
		&models.WComment{},
		// 处理当前语句逻辑。
		&models.WExternalLink{},
		// 处理当前语句逻辑。
		&models.WHomePopup{},
		// 处理当前语句逻辑。
		&models.WSMSChannel{},
		// 进入新的代码块进行处理。
	); err != nil {
		// 调用logger.Fatal完成当前处理。
		logger.Fatal("Failed to migrate database: %v", err)
	}

	// 6) 创建 token 存储实现（Store）。
	// - RedisStore：支持共享状态，适合多实例部署；
	// - MemoryStore：仅本进程可见，仅建议本地开发使用。
	var store tokenstore.Store
	// 判断条件并进入对应分支逻辑。
	if redisClient != nil {
		// 创建Redis存储适配器
		store = tokenstore.NewRedisStoreWithClient(redisClient)
		// 进入新的代码块进行处理。
	} else {
		// 判断条件并进入对应分支逻辑。
		if !cfg.JWT.AllowInsecureMemoryStore {
			// 返回当前处理结果。
			return fmt.Errorf("未配置redis，且jwt.allow_insecure_memory_store=false")
		}
		// 使用内存存储
		store = tokenstore.NewMemoryStore()
		// 调用logger.Warn完成当前处理。
		logger.Warn("Using in-memory token store (not suitable for production)")
	}

	// 7) 配置 Token 管理器（签发 / 校验 / 刷新 / 撤销都在这里完成）。
	jwtCfg := tokenjwt.DefaultConfig()
	// 更新当前变量或字段值。
	jwtCfg.SigningKey = signingKey
	// 更新当前变量或字段值。
	jwtCfg.AccessExpire = time.Duration(cfg.JWT.AccessExpire) * time.Second
	// 更新当前变量或字段值。
	jwtCfg.RefreshExpire = time.Duration(cfg.JWT.RefreshExpire) * time.Second
	// 更新当前变量或字段值。
	jwtCfg.SessionIdleTimeout = time.Duration(cfg.JWT.SessionIdleTimeout) * time.Second

	// 定义并初始化当前变量。
	mgr := tokenjwt.NewManager(jwtCfg, store)

	// 8) 装配路由与中间件。
	r := SetupRouter(mgr, db, redisClient, logger, cfg)

	// 定义并初始化当前变量。
	addr := fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port)
	// 调用logger.Info完成当前处理。
	logger.Info("Server starting on %s", addr)

	// 返回当前处理结果。
	return r.Run(addr)
}
