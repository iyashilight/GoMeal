package middleware

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

// RateLimit 基于 Redis 有序集合（ZSET）实现的滑动窗口限流中间件
// 使用 Lua 脚本保证原子性，记录每个 IP 在窗口内的请求次数
// 参数：limit — 窗口内允许的最大请求数，window — 时间窗口
func RateLimit(rclient *redis.Client, limit int, window time.Duration) gin.HandlerFunc {
	script := redis.NewScript(`
	local key = KEYS[1]
	local now = tonumber(ARGV[1])
	local window = tonumber(ARGV[2])
	local limit = tonumber(ARGV[3])

	-- 删除窗口外的记录
	redis.call('ZREMRANGEBYSCORE', key, 0, now - window)

	-- 统计当前窗口内请求数
	local count = redis.call('ZCARD', key)

	if count >= limit then
		return 0
	end

	-- 添加当前请求
	redis.call('ZADD', key, now, now)
	redis.call('EXPIRE', key, math.ceil(window / 1000))
	return 1
	`)

	return func(c *gin.Context) {
		key := fmt.Sprintf("ratelimit:%s", c.ClientIP())
		now := time.Now().UnixMilli()
		allowed, err := script.Run(context.Background(), rclient, []string{key},
			now, window.Milliseconds(), limit).Int()
		if err != nil || allowed == 0 {
			c.JSON(http.StatusTooManyRequests, gin.H{
				"code":    429,
				"message": "请求太频繁，请稍后再试",
			})
			c.Abort()
			return
		}
		c.Next()
	}
}
