package dal

import (
	"github.com/go-redis/redis/v8"
)

var RedisDB *redis.Client

func RedisInit() {
	RedisDB = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})
}
