package store

import (
	"context"
	"fmt"
	"go-admin/internal/constants"
	"sync"
	"time"

	"github.com/go-redis/redis/v8"
	redisx "github.com/wangyahua6688-maker/tk-common/utils/redisx/v8"
)

// Store 定义 token 存储抽象。
// 业务约束：
// - Set/Get/Delete 需满足并发安全；
// - expire 语义需尽量与 Redis 行为一致；
// - Get 未命中统一返回 constants.ErrTokenNotFound。
type Store interface {
	// 调用Set完成当前处理。
	Set(key string, value string, expire time.Duration) error
	// 调用Get完成当前处理。
	Get(key string) (string, error)
	// 调用Delete完成当前处理。
	Delete(key string) error
	// 调用Ping完成当前处理。
	Ping() error
}

// 定义类型常量，用于识别存储类型
type StoreType string

// 声明当前常量。
const (
	// 更新当前变量或字段值。
	StoreTypeRedis StoreType = "redis"
	// 更新当前变量或字段值。
	StoreTypeMemory StoreType = "memory"
)

// StoreWithType 扩展Store接口，添加类型信息
type StoreWithType interface {
	// 处理当前语句逻辑。
	Store
	// 调用GetType完成当前处理。
	GetType() StoreType
}

// RedisStore 基于 Redis 的存储实现（生产推荐）。
type RedisStore struct {
	// 处理当前语句逻辑。
	client *redis.Client
	// 处理当前语句逻辑。
	ctx context.Context
}

// NewRedisStoreWithClient 使用现有Redis客户端创建存储
func NewRedisStoreWithClient(client *redis.Client) Store {
	// 返回当前处理结果。
	return &RedisStore{client: client, ctx: context.Background()}
}

// GetType 返回存储类型
func (r *RedisStore) GetType() StoreType {
	// 返回当前处理结果。
	return StoreTypeRedis
}

// Set 实现
func (r *RedisStore) Set(key string, value string, expire time.Duration) error {
	// 返回当前处理结果。
	return redisx.SetString(r.ctx, r.client, key, value, expire)
}

// Get 实现
func (r *RedisStore) Get(key string) (string, error) {
	// 定义并初始化当前变量。
	v, hit, err := redisx.GetString(r.ctx, r.client, key)
	// 判断条件并进入对应分支逻辑。
	if err != nil {
		// 返回当前处理结果。
		return "", err
	}
	// 判断条件并进入对应分支逻辑。
	if !hit {
		// 返回当前处理结果。
		return "", constants.ErrTokenNotFound
	}
	// 返回当前处理结果。
	return v, nil
}

// Delete 实现
func (r *RedisStore) Delete(key string) error {
	// 返回当前处理结果。
	return redisx.Del(r.ctx, r.client, key)
}

// Ping 实现
func (r *RedisStore) Ping() error {
	// 定义并初始化当前变量。
	_, err := r.client.Ping(r.ctx).Result()
	// 判断条件并进入对应分支逻辑。
	if err != nil {
		// 返回当前处理结果。
		return fmt.Errorf("redis ping failed: %w", err)
	}
	// 返回当前处理结果。
	return nil
}

// MemoryStore 内存存储实现（仅建议开发环境使用）。
// 注意：
// - 数据仅当前进程可见，重启即丢失；
// - 已实现基础过期语义与并发保护，便于本地调试行为接近 Redis。
type MemoryStore struct {
	// 处理当前语句逻辑。
	mu sync.RWMutex
	// 处理当前语句逻辑。
	data map[string]memoryEntry
}

// memoryEntry 内存存储条目，包含过期时间（零值表示不过期）。
type memoryEntry struct {
	// 处理当前语句逻辑。
	value string
	// 处理当前语句逻辑。
	expireAt time.Time
}

// NewMemoryStore 创建内存存储
func NewMemoryStore() Store {
	// 返回当前处理结果。
	return &MemoryStore{
		// 调用make完成当前处理。
		data: make(map[string]memoryEntry),
	}
}

// GetType 返回存储类型
func (m *MemoryStore) GetType() StoreType {
	// 返回当前处理结果。
	return StoreTypeMemory
}

// Set 实现
func (m *MemoryStore) Set(key string, value string, expire time.Duration) error {
	// 调用m.mu.Lock完成当前处理。
	m.mu.Lock()
	// 注册延迟执行逻辑。
	defer m.mu.Unlock()

	// 定义并初始化当前变量。
	entry := memoryEntry{value: value}
	// 判断条件并进入对应分支逻辑。
	if expire > 0 {
		// 更新当前变量或字段值。
		entry.expireAt = time.Now().Add(expire)
	}
	// 更新当前变量或字段值。
	m.data[key] = entry
	// 返回当前处理结果。
	return nil
}

// Get 实现
func (m *MemoryStore) Get(key string) (string, error) {
	// 调用m.mu.RLock完成当前处理。
	m.mu.RLock()
	// 定义并初始化当前变量。
	entry, exists := m.data[key]
	// 调用m.mu.RUnlock完成当前处理。
	m.mu.RUnlock()

	// 判断条件并进入对应分支逻辑。
	if !exists {
		// 返回当前处理结果。
		return "", constants.ErrTokenNotFound
	}

	// 兼容过期能力：读取时惰性清理，保证本地开发与 Redis 行为一致。
	if !entry.expireAt.IsZero() && time.Now().After(entry.expireAt) {
		// 调用m.mu.Lock完成当前处理。
		m.mu.Lock()
		// 调用delete完成当前处理。
		delete(m.data, key)
		// 调用m.mu.Unlock完成当前处理。
		m.mu.Unlock()
		// 返回当前处理结果。
		return "", constants.ErrTokenNotFound
	}

	// 返回当前处理结果。
	return entry.value, nil
}

// Delete 实现
func (m *MemoryStore) Delete(key string) error {
	// 调用m.mu.Lock完成当前处理。
	m.mu.Lock()
	// 注册延迟执行逻辑。
	defer m.mu.Unlock()
	// 调用delete完成当前处理。
	delete(m.data, key)
	// 返回当前处理结果。
	return nil
}

// Ping 实现（内存存储总是成功的）
func (m *MemoryStore) Ping() error {
	// 返回当前处理结果。
	return nil
}
