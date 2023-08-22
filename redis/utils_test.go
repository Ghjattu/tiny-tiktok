package redis

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestRandomDay(t *testing.T) {
	day := RandomDay()

	assert.GreaterOrEqual(t, day, 24*time.Hour)
	assert.LessOrEqual(t, day, 48*time.Hour)
}

func TestHashIncrBy(t *testing.T) {
	t.Run("key does not exist", func(t *testing.T) {
		Rdb.FlushDB(Ctx)

		val, statusMsg, err := HashIncrBy("test", "test", 1)

		assert.Equal(t, int64(0), val)
		assert.Equal(t, "key does not exist", statusMsg)
		assert.Nil(t, err)
	})

	t.Run("increase", func(t *testing.T) {
		Rdb.FlushDB(Ctx)

		Rdb.HSet(Ctx, "test_hash", "test", 1)
		val, _, err := HashIncrBy("test_hash", "test", 1)

		assert.Equal(t, int64(2), val)
		assert.Nil(t, err)
	})

	t.Run("decrease", func(t *testing.T) {
		Rdb.FlushDB(Ctx)

		Rdb.HSet(Ctx, "test_hash", "test", 1)
		val, _, err := HashIncrBy("test_hash", "test", -1)

		assert.Equal(t, int64(0), val)
		assert.Nil(t, err)
	})
}

func TestHashGetAll(t *testing.T) {
	t.Run("key does not exist", func(t *testing.T) {
		Rdb.FlushDB(Ctx)

		result, err := HashGetAll("test")

		assert.Nil(t, result)
		assert.NotNil(t, err)
		assert.Equal(t, "key does not exist", err.Error())
	})

	t.Run("key exists", func(t *testing.T) {
		Rdb.FlushDB(Ctx)

		Rdb.HSet(Ctx, "test_hash", "test", 1)
		result, err := HashGetAll("test_hash")

		assert.Nil(t, err)
		assert.Equal(t, "1", result.Val()["test"])
	})
}
