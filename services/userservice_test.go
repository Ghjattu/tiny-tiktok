package services

import (
	"strconv"
	"testing"

	"github.com/Ghjattu/tiny-tiktok/models"
	"github.com/Ghjattu/tiny-tiktok/redis"
	"github.com/stretchr/testify/assert"
)

var (
	userService = &UserService{}
)

func TestGetUserByUserID(t *testing.T) {
	setup()

	t.Run("user does not exist", func(t *testing.T) {
		statusCode, statusMsg, _ := userService.GetUserByUserID(testUserOne.ID + 100)

		assert.Equal(t, int32(1), statusCode)
		assert.Equal(t, "user not found", statusMsg)
	})

	t.Run("get user successfully", func(t *testing.T) {
		statusCode, statusMsg, user := userService.GetUserByUserID(testUserOne.ID)

		assert.Equal(t, int32(0), statusCode)
		assert.Equal(t, "get user successfully", statusMsg)
		assert.Equal(t, "", user.Password)
	})
}

func TestGetUserDetailByUserID(t *testing.T) {
	setup()

	t.Run("user does not exist", func(t *testing.T) {
		redis.Rdb.FlushDB(redis.Ctx)

		statusCode, statusMsg, _ := userService.GetUserDetailByUserID(0, testUserOne.ID+100)

		assert.Equal(t, int32(1), statusCode)
		assert.Equal(t, "user not found", statusMsg)
	})

	t.Run("get user successfully with cache miss", func(t *testing.T) {
		redis.Rdb.FlushDB(redis.Ctx)

		statusCode, statusMsg, user := userService.GetUserDetailByUserID(0, testUserOne.ID)
		waitForConsumer()

		userKey := redis.UserKey + strconv.FormatInt(testUserOne.ID, 10)
		id := redis.Rdb.HGet(redis.Ctx, userKey, "id").Val()
		name := redis.Rdb.HGet(redis.Ctx, userKey, "name").Val()
		workCount := redis.Rdb.HGet(redis.Ctx, userKey, "work_count").Val()

		assert.Equal(t, int32(0), statusCode)
		assert.Equal(t, "get user successfully", statusMsg)
		assert.Equal(t, testUserOne.Name, user.Name)
		assert.Equal(t, strconv.FormatInt(testUserOne.ID, 10), id)
		assert.Equal(t, testUserOne.Name, name)
		assert.Equal(t, strconv.FormatInt(user.WorkCount, 10), workCount)
	})

	t.Run("get user successfully with cache hit", func(t *testing.T) {
		redis.Rdb.FlushDB(redis.Ctx)

		// Insert test user into redis.
		testUserDetail := &models.UserDetail{ID: testUserOne.ID, Name: testUserOne.Name}
		userKey := redis.UserKey + strconv.FormatInt(testUserOne.ID, 10)
		redis.Rdb.HSet(redis.Ctx, userKey, testUserDetail)

		statusCode, statusMsg, userDetail := userService.GetUserDetailByUserID(0, testUserOne.ID)

		assert.Equal(t, int32(0), statusCode)
		assert.Equal(t, "get user successfully", statusMsg)
		assert.Equal(t, testUserDetail.Name, userDetail.Name)
	})
}

func BenchmarkGetUserDetailByUserID(b *testing.B) {
	setup()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		userService.GetUserDetailByUserID(testUserOne.ID, testUserOne.ID)
	}
}
