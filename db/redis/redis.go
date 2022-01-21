package redis

import (
	"github.com/go-redis/redis/v8"
	"github.com/joho/godotenv"
	"github.com/junminhong/member-services-center/pkg/logger"
	"os"
)

var sugar = logger.Setup()

func init() {
	err := godotenv.Load()
	if err != nil {
		sugar.Info(err.Error())
	}
}

func Setup() *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_HOST") + ":" + os.Getenv("REDIS_PORT"),
		Password: os.Getenv("REDIS_PASSWORD"), // no password set
		DB:       0,                           // use default DB
	})
	return client
}
