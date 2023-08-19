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
	godotenv.Load("../../.env")

	redis_ip := os.Getenv("REDIS_IP")
	redis_port := os.Getenv("REDIS_PORT")
	redis_password := os.Getenv("REDIS_PASSWORD")

	Rdb = redis.NewClient(&redis.Options{
		Addr:     redis_ip + ":" + redis_port,
		Password: redis_password,
		DB:       0,
	})

	_, err := Rdb.Ping(Ctx).Result()
	if err != nil {
		panic("failed to connect redis")
	}
}
