package rabbitmq

import (
	"testing"
	"time"

	"github.com/Ghjattu/tiny-tiktok/redis"
	"github.com/stretchr/testify/assert"
)

func TestProducer(t *testing.T) {
	redis.Rdb.FlushDB(redis.Ctx)

	RedisMQ.Producer(&Message{
		Type:       "Hash",
		SubType:    "Set",
		StructName: "VideoCache",
		Key:        "test",
		Value:      &redis.VideoCache{ID: 1},
	})

	time.Sleep(100 * time.Millisecond)

	id, err := redis.Rdb.HGet(redis.Ctx, "test", "id").Result()

	assert.Nil(t, err)
	assert.Equal(t, "1", id)
}
