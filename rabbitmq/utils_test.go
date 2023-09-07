package rabbitmq

import (
	"testing"
	"time"

	"github.com/Ghjattu/tiny-tiktok/redis"
	"github.com/stretchr/testify/assert"
)

func TestProduceMessage(t *testing.T) {
	ProduceMessage("Hash", "Set", "VideoCache", "test", "", &redis.VideoCache{ID: 1})

	time.Sleep(100 * time.Millisecond)

	id, err := redis.Rdb.HGet(redis.Ctx, "test", "id").Result()

	assert.Nil(t, err)
	assert.Equal(t, "1", id)
}

func TestConsumeMessage(t *testing.T) {
	t.Run("type is Hash and sub type is set", func(t *testing.T) {
		redis.Rdb.FlushDB(redis.Ctx)

		message := &Message{
			Type:       "Hash",
			SubType:    "Set",
			Key:        "test",
			StructName: "VideoCache",
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

	t.Run("type is List and sub type is RPush", func(t *testing.T) {
		redis.Rdb.FlushDB(redis.Ctx)

		message := &Message{
			Type:    "List",
			SubType: "RPush",
			Key:     "test",
			Value:   []int64{1, 2, 3},
		}
		ConsumeMessage(message)

		len, err := redis.Rdb.LLen(redis.Ctx, message.Key).Result()

		assert.Nil(t, err)
		assert.Equal(t, int64(3), len)
	})

	t.Run("type is List and sub type is RPushX", func(t *testing.T) {
		redis.Rdb.FlushDB(redis.Ctx)
		key := "test"
		redis.Rdb.RPush(redis.Ctx, key, []string{"1", "2", "3"})

		message := &Message{
			Type:    "List",
			SubType: "RPushX",
			Key:     key,
			Value:   []int64{4, 5, 6},
		}
		ConsumeMessage(message)

		len, err := redis.Rdb.LLen(redis.Ctx, message.Key).Result()

		assert.Nil(t, err)
		assert.Equal(t, int64(6), len)
	})

	t.Run("type is List and sub type is LRem", func(t *testing.T) {
		redis.Rdb.FlushDB(redis.Ctx)
		key := "test"
		redis.Rdb.RPush(redis.Ctx, key, []string{"1", "2", "3"})

		message := &Message{
			Type:    "List",
			SubType: "LRem",
			Key:     "test",
			Value:   int64(2),
		}
		ConsumeMessage(message)

		len, err := redis.Rdb.LLen(redis.Ctx, message.Key).Result()

		assert.Nil(t, err)
		assert.Equal(t, int64(2), len)
	})
}

func TestCacheStructSelector(t *testing.T) {
	t.Run("name is VideoCache", func(t *testing.T) {
		redis.Rdb.FlushDB(redis.Ctx)

		videoCache := &redis.VideoCache{ID: 1}

		got := CacheStructSelector("VideoCache", videoCache)
		returnedValue, ok := got.(*redis.VideoCache)

		assert.True(t, ok)
		assert.Equal(t, videoCache.ID, returnedValue.ID)
	})

	t.Run("name is CommentCache", func(t *testing.T) {
		redis.Rdb.FlushDB(redis.Ctx)

		commentCache := &redis.CommentCache{ID: 1}

		got := CacheStructSelector("CommentCache", commentCache)
		returnedValue, ok := got.(*redis.CommentCache)

		assert.True(t, ok)
		assert.Equal(t, commentCache.ID, returnedValue.ID)
	})

	t.Run("name is UserCache", func(t *testing.T) {
		redis.Rdb.FlushDB(redis.Ctx)

		userCache := &redis.UserCache{ID: 1}

		got := CacheStructSelector("UserCache", userCache)
		returnedValue, ok := got.(*redis.UserCache)

		assert.True(t, ok)
		assert.Equal(t, userCache.ID, returnedValue.ID)
	})
}
