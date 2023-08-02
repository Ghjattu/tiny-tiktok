package redis

import (
	"context"
	"os"

	_ "github.com/joho/godotenv/autoload"
	"github.com/redis/go-redis/v9"
)

var Rdb *redis.Client
var ctx = context.Background()

func init() {
	redis_ip := os.Getenv("REDIS_IP")
	redis_port := os.Getenv("REDIS_PORT")
	redis_password := os.Getenv("REDIS_PASSWORD")

	Rdb = redis.NewClient(&redis.Options{
		Addr:     redis_ip + ":" + redis_port,
		Password: redis_password,
		DB:       0,
	})

	_, err := Rdb.Ping(ctx).Result()
	if err != nil {
		panic("failed to connect redis")
	}
}
