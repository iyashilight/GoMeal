package cache

import (
	"context"
	"fmt"

	"github.com/redis/go-redis/v9"
)

// StockCache 基于 Redis + Lua 脚本实现的库存缓存，保证扣减的原子性
// 用于订单创建时预扣库存，避免超卖
type StockCache struct {
	client     *redis.Client
	decrScript *redis.Script // 扣减库存 Lua 脚本：返回 -1（key 不存在），0（库存不足），1（成功）
	incrScript *redis.Script // 恢复库存 Lua 脚本：直接增加
}

// NewStockCache 创建库存缓存，注册 Lua 脚本
func NewStockCache(client *redis.Client) *StockCache {
	return &StockCache{client: client,
		decrScript: redis.NewScript(`
			local key = KEYS[1]
			local qty = tonumber(ARGV[1])
			local stock = redis.call('GET',key)
			if not stock then
				return -1
			end
			stock = tonumber(stock)
			if stock < qty then
				return 0
			end
			redis.call('DECRBY',key,qty)
			return 1
		`),
		incrScript: redis.NewScript(`
			local key = KEYS[1]
			redis.call('INCRBY',key,ARGV[1])
			return 1
		`)}
}

// stockKey 生成库存 Redis 键：food:stock:{foodID}
func (s *StockCache) stockKey(foodID uint) string {
	return fmt.Sprintf("food:stock:%d", foodID)
}

// Decrease 尝试扣减库存，返回：-1（key 不存在，需初始化），0（库存不足），1（扣减成功）
func (s *StockCache) Decrease(foodID uint, quantity int) (int, error) {
	result := s.decrScript.Run(context.Background(), s.client,
		[]string{s.stockKey(foodID)}, quantity)
	res, err := result.Int()
	if err != nil {
		return 0, err
	}
	return res, nil
}

// Increase 恢复/增加库存，用于订单取消或回滚时释放 Redis 预扣库存
func (s *StockCache) Increase(foodID uint, quantity int) error {
	result := s.incrScript.Run(context.Background(), s.client,
		[]string{s.stockKey(foodID)}, quantity)
	_, err := result.Int()
	if err != nil {
		return err
	}
	return nil
}

// Init 初始化 Redis 中的库存值，通常在查询商品详情后调用
func (s *StockCache) Init(foodID uint, stock int) error {
	return s.client.Set(context.Background(), s.stockKey(foodID), stock, 0).Err()
}
