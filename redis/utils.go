package redis

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/redis/go-redis/v9"
)

// RandomDay returns a random time.Duration between 24h and 48h.
//
//	@return time.Duration
func RandomDay() time.Duration {
	r := rand.New(rand.NewSource(time.Now().Unix()))
	return time.Hour * time.Duration(24+r.Int31n(25))
}

// HashIncrBy calls HIncrBy command of redis.
//
//	@param key string
//	@param field string
//	@param incr int64
//	@return int64 "the value of the field after the increment"
//	@return error
func HashIncrBy(key, field string, incr int64) (int64, error) {
	if Rdb.Exists(Ctx, key).Val() == 1 {
		res := Rdb.HIncrBy(Ctx, key, field, incr)

		return res.Val(), res.Err()
	}

	return 0, fmt.Errorf("failed to hash incr by: key %s does not exist", key)
}

// HashGetAll calls HGetAll command of redis.
//
//	@param key string
//	@return *redis.MapStringStringCmd
//	@return error
func HashGetAll(key string) (*redis.MapStringStringCmd, error) {
	exist, err := Rdb.Exists(Ctx, key).Result()
	if err != nil {
		return nil, err
	}

	if exist == 1 {
		result := Rdb.HGetAll(Ctx, key)

		return result, result.Err()
	}

	return nil, fmt.Errorf("key does not exist")
}
