package config

import (
	"context"
	"log/slog"

	"github.com/redis/go-redis/v9"
)

// InitRedis 创建 Redis 客户端并验证连接可用
func InitRedis(cfg *Config) (*redis.Client, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     cfg.Redis.Addr,
		Password: cfg.Redis.Password,
		DB:       cfg.Redis.DB,
	})
	_, err := client.Ping(context.Background()).Result()
	if err != nil {
		return nil, err
	}
	slog.Info("redis connected", "addr", cfg.Redis.Addr)
	return client, nil
}
