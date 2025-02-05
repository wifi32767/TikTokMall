package dal

import (
	"context"

	"github.com/redis/go-redis/v9"
	"github.com/wifi32767/TikTokMall/app/product/conf"
)

var RedisClient *redis.Client

func RedisInit() {
	RedisClient = redis.NewClient(&redis.Options{
		Addr:     conf.GetConf().Redis.Address,
		Password: conf.GetConf().Redis.Password,
		DB:       0,
	})
	if err := RedisClient.Ping(context.Background()).Err(); err != nil {
		panic(err)
	}
}
