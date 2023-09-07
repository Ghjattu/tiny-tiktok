package redis

import (
	"context"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/joho/godotenv/autoload"
	"github.com/redis/go-redis/v9"
)

var Rdb *redis.Client
var Ctx = context.Background()

func init() {
	godotenv.Load("../.env")

	REDIS_IP := os.Getenv("REDIS_IP")
	if REDIS_IP == "" {
		REDIS_IP = "127.0.0.1"
	}

	REDIS_PORT := os.Getenv("REDIS_PORT")
	if REDIS_PORT == "" {
		REDIS_PORT = "6379"
	}

	REDIS_PASSWORD := os.Getenv("REDIS_PASSWORD")

	Rdb = redis.NewClient(&redis.Options{
		Addr:     REDIS_IP + ":" + REDIS_PORT,
		Password: REDIS_PASSWORD,
		DB:       0,
	})

	_, err := Rdb.Ping(Ctx).Result()
	if err != nil {
		panic("failed to connect redis")
	}
}
