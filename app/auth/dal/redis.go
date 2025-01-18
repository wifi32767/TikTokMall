package dal

import (
	"context"

	"github.com/go-redis/redis/v8"
	"github.com/wifi32767/TikTokMall/app/auth/conf"
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
