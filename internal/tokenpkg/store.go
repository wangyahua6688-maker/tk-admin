package tokenpkg

import (
	"context"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
)

// Store 接口定义
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

// RedisStore 定义（添加GetType方法）
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
			return "", ErrTokenNotFound
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

// MemoryStore 内存存储实现
type MemoryStore struct {
	data map[string]string
}

// NewMemoryStore 创建内存存储
func NewMemoryStore() Store {
	return &MemoryStore{
		data: make(map[string]string),
	}
}

// GetType 返回存储类型
func (m *MemoryStore) GetType() StoreType {
	return StoreTypeMemory
}

// Set 实现
func (m *MemoryStore) Set(key string, value string, expire time.Duration) error {
	m.data[key] = value
	// 注意：内存存储不支持自动过期，需要自己实现清理逻辑
	return nil
}

// Get 实现
func (m *MemoryStore) Get(key string) (string, error) {
	value, exists := m.data[key]
	if !exists {
		return "", ErrTokenNotFound
	}
	return value, nil
}

// Delete 实现
func (m *MemoryStore) Delete(key string) error {
	delete(m.data, key)
	return nil
}

// Ping 实现（内存存储总是成功的）
func (m *MemoryStore) Ping() error {
	return nil
}
