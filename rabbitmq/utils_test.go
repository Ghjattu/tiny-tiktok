package rabbitmq

import (
	"testing"

	"github.com/Ghjattu/tiny-tiktok/redis"
	"github.com/stretchr/testify/assert"
)

func TestConsumeMessage(t *testing.T) {
	t.Run("type is Hash and sub type is set", func(t *testing.T) {
		redis.Rdb.FlushDB(redis.Ctx)

		message := &Message{
			Type:       "Hash",
			SubType:    "set",
			Key:        "test",
			StructName: "video",
			Value:      &redis.VideoCache{ID: 1},
		}
		ConsumeMessage(message)

		id, err := redis.Rdb.HGet(redis.Ctx, message.Key, "id").Result()

		assert.Nil(t, err)
		assert.Equal(t, "1", id)
	})

	t.Run("type is Hash and sub type is Incr", func(t *testing.T) {
		redis.Rdb.FlushDB(redis.Ctx)
		videoCache := &redis.VideoCache{ID: 1}
		redis.Rdb.HSet(redis.Ctx, "test", videoCache)

		message := &Message{
			Type:    "Hash",
			SubType: "Incr",
			Key:     "test",
			Field:   "id",
			Value:   int64(1),
		}
		ConsumeMessage(message)

		id, err := redis.Rdb.HGet(redis.Ctx, message.Key, message.Field).Result()

		assert.Nil(t, err)
		assert.Equal(t, "2", id)
	})

	t.Run("type is list and sub type is RPush", func(t *testing.T) {
		redis.Rdb.FlushDB(redis.Ctx)

		message := &Message{
			Type:    "list",
			SubType: "RPush",
			Key:     "test",
			Value:   []int64{1, 2, 3},
		}
		ConsumeMessage(message)

		len, err := redis.Rdb.LLen(redis.Ctx, message.Key).Result()

		assert.Nil(t, err)
		assert.Equal(t, int64(3), len)
	})
}

func TestCacheStructSelector(t *testing.T) {
	t.Run("name is video", func(t *testing.T) {
		redis.Rdb.FlushDB(redis.Ctx)

		videoCache := &redis.VideoCache{ID: 1}

		got := CacheStructSelector("video", videoCache)
		returnedValue, ok := got.(*redis.VideoCache)

		assert.True(t, ok)
		assert.Equal(t, videoCache.ID, returnedValue.ID)
	})

	t.Run("name is comment", func(t *testing.T) {
		redis.Rdb.FlushDB(redis.Ctx)

		commentCache := &redis.CommentCache{ID: 1}

		got := CacheStructSelector("comment", commentCache)
		returnedValue, ok := got.(*redis.CommentCache)

		assert.True(t, ok)
		assert.Equal(t, commentCache.ID, returnedValue.ID)
	})
}
