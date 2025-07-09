package redis_cache

import (
	"context"
	"encoding/json"
	"time"

	"github.com/bits-and-blooms/bloom"
	"github.com/redis/go-redis/v9"
	"golang.org/x/exp/rand"
	"gorm.io/gorm"
)

var CacheMiss error = redis.Nil

type RedisCache struct {
	client *redis.Client
	filter *bloom.BloomFilter
}

func NewRedisCache(client *redis.Client, filter *bloom.BloomFilter) *RedisCache {
	return &RedisCache{
		client: client,
		filter: filter,
	}
}

func (c *RedisCache) Set(ctx context.Context, key string, value any) {
	c.filter.AddString(key)
	encoded, err := json.Marshal(value)
	if err != nil {
		return
	}
	_ = c.client.Set(ctx, key, encoded, randomDuration())
}

func (c *RedisCache) Del(ctx context.Context, key string) error {
	return c.client.Del(ctx, key).Err()
}

func (c *RedisCache) Get(ctx context.Context, key string, result any) error {
	// 布隆过滤器判断
	if !c.filter.TestString(key) {
		return gorm.ErrRecordNotFound
	}
	// 读取redis缓存
	cachedResult := c.client.Get(ctx, key)
	err := cachedResult.Err()
	if err != nil {
		return err
	}
	cachedResultByte, err := cachedResult.Bytes()
	if err != nil {
		return err
	}
	// 续期
	_ = c.client.Expire(ctx, key, randomDuration())
	err = json.Unmarshal(cachedResultByte, &result)
	if err != nil {
		return err
	}
	return nil
}

// 在失效时间内加入随机数，可以防止缓存雪崩的情况
func randomDuration() time.Duration {
	return time.Duration(1+rand.Intn(2)) * time.Hour
}
