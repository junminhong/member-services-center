package redis

import (
	"github.com/go-redis/redis/v8"
)

func InitRedis() *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr:     "127.0.0.1:6379",
		Password: "@test#test", // no password set
		DB:       0,            // use default DB
	})
	return client
}
