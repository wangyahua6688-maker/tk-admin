package store

import (
	"context"
	"fmt"
	"go-admin-full/internal/constants"
	"sync"
	"time"

	"github.com/go-redis/redis/v8"
)

// Store 定义 token 存储抽象。
// 业务约束：
// - Set/Get/Delete 需满足并发安全；
// - expire 语义需尽量与 Redis 行为一致；
// - Get 未命中统一返回 constants.ErrTokenNotFound。
type Store interface {
	Set(key string, value string, expire time.Duration) error
	Get(key string) (string, error)
	Delete(key string) error
	Ping() error
}

// 定义类型常量，用于识别存储类型
type StoreType string

const (
	StoreTypeRedis  StoreType = "redis"
	StoreTypeMemory StoreType = "memory"
)

// StoreWithType 扩展Store接口，添加类型信息
type StoreWithType interface {
	Store
	GetType() StoreType
}

// RedisStore 基于 Redis 的存储实现（生产推荐）。
type RedisStore struct {
	client *redis.Client
	ctx    context.Context
}

// NewRedisStore 创建Redis存储
func NewRedisStore(addr, password string, db int) Store {
	cli := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
	})
	return &RedisStore{client: cli, ctx: context.Background()}
}

// NewRedisStoreWithClient 使用现有Redis客户端创建存储
func NewRedisStoreWithClient(client *redis.Client) Store {
	return &RedisStore{client: client, ctx: context.Background()}
}

// GetType 返回存储类型
func (r *RedisStore) GetType() StoreType {
	return StoreTypeRedis
}

// Set 实现
func (r *RedisStore) Set(key string, value string, expire time.Duration) error {
	return r.client.Set(r.ctx, key, value, expire).Err()
}

// Get 实现
func (r *RedisStore) Get(key string) (string, error) {
	v, err := r.client.Get(r.ctx, key).Result()
	if err != nil {
		if err == redis.Nil {
			return "", constants.ErrTokenNotFound
		}
		return "", err
	}
	return v, nil
}

// Delete 实现
func (r *RedisStore) Delete(key string) error {
	return r.client.Del(r.ctx, key).Err()
}

// Ping 实现
func (r *RedisStore) Ping() error {
	_, err := r.client.Ping(r.ctx).Result()
	if err != nil {
		return fmt.Errorf("redis ping failed: %w", err)
	}
	return nil
}

// MemoryStore 内存存储实现（仅建议开发环境使用）。
// 注意：
// - 数据仅当前进程可见，重启即丢失；
// - 已实现基础过期语义与并发保护，便于本地调试行为接近 Redis。
type MemoryStore struct {
	mu   sync.RWMutex
	data map[string]memoryEntry
}

// memoryEntry 内存存储条目，包含过期时间（零值表示不过期）。
type memoryEntry struct {
	value    string
	expireAt time.Time
}

// NewMemoryStore 创建内存存储
func NewMemoryStore() Store {
	return &MemoryStore{
		data: make(map[string]memoryEntry),
	}
}

// GetType 返回存储类型
func (m *MemoryStore) GetType() StoreType {
	return StoreTypeMemory
}

// Set 实现
func (m *MemoryStore) Set(key string, value string, expire time.Duration) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	entry := memoryEntry{value: value}
	if expire > 0 {
		entry.expireAt = time.Now().Add(expire)
	}
	m.data[key] = entry
	return nil
}

// Get 实现
func (m *MemoryStore) Get(key string) (string, error) {
	m.mu.RLock()
	entry, exists := m.data[key]
	m.mu.RUnlock()

	if !exists {
		return "", constants.ErrTokenNotFound
	}

	// 兼容过期能力：读取时惰性清理，保证本地开发与 Redis 行为一致。
	if !entry.expireAt.IsZero() && time.Now().After(entry.expireAt) {
		m.mu.Lock()
		delete(m.data, key)
		m.mu.Unlock()
		return "", constants.ErrTokenNotFound
	}

	return entry.value, nil
}

// Delete 实现
func (m *MemoryStore) Delete(key string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	delete(m.data, key)
	return nil
}

// Ping 实现（内存存储总是成功的）
func (m *MemoryStore) Ping() error {
	return nil
}
