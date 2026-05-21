// Package cache 封装 Redis 缓存操作，提供 JSON 序列化的快捷存取
package cache

import (
	"context"
	"encoding/json"
	"time"

	"github.com/redis/go-redis/v9"
)

// Cache 通用 JSON 缓存，用于商家列表、商品详情等读多写少的数据
type Cache struct {
	client *redis.Client
}

// NewCache 创建缓存实例
func NewCache(client *redis.Client) *Cache {
	return &Cache{client: client}
}

// GetJSON 从 Redis 获取 JSON 键并反序列化到 dest，返回是否命中缓存
func (c *Cache) GetJSON(key string, dest interface{}) (bool, error) {
	data, err := c.client.Get(context.Background(), key).Bytes()
	if err == redis.Nil {
		return false, nil // 缓存未命中
	}
	if err != nil {
		return false, err // 其他 Redis 错误
	}
	err = json.Unmarshal(data, dest)
	return true, err
}

// SetJSON 将值序列化为 JSON 存入 Redis，并指定 TTL
func (c *Cache) SetJSON(key string, value interface{}, ttl time.Duration) error {
	data, err := json.Marshal(value)
	if err != nil {
		return err
	}
	return c.client.Set(context.Background(), key, data, ttl).Err()
}

// Del 删除指定的缓存键
func (c *Cache) Del(key string) error {
	return c.client.Del(context.Background(), key).Err()
}
