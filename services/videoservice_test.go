package services

import (
	"strconv"
	"testing"
	"time"

	"github.com/Ghjattu/tiny-tiktok/models"
	"github.com/Ghjattu/tiny-tiktok/redis"
	"github.com/stretchr/testify/assert"
)

var (
	videoService = &VideoService{}
)

func TestCreateNewVideo(t *testing.T) {
	setup()

	t.Run("create video successfully", func(t *testing.T) {
		redis.Rdb.FlushDB(redis.Ctx)

		// Insert a test user to redis.
		userKey := redis.UserKey + strconv.FormatInt(testUserOne.ID, 10)
		redis.Rdb.HSet(redis.Ctx, userKey, testUserOneDetail)

		statusCode, statusMsg :=
			videoService.CreateNewVideo("test", "test", testUserOne.ID, time.Now())
		workCount := redis.Rdb.HGet(redis.Ctx, userKey, "work_count").Val()

		assert.Equal(t, int32(0), statusCode)
		assert.Equal(t, "create new video successfully", statusMsg)
		assert.Equal(t, "1", workCount)
	})
}

func TestGetVideoListByAuthorID(t *testing.T) {
	setup()

	t.Run("get video list successfully with cache miss", func(t *testing.T) {
		redis.Rdb.FlushDB(redis.Ctx)

		statusCode, _, videoList :=
			videoService.GetVideoListByAuthorID(testUserOne.ID, 0)

		assert.Equal(t, int32(0), statusCode)
		assert.Equal(t, 2, len(videoList))
	})

	t.Run("get video list successfully with cache hit", func(t *testing.T) {
		redis.Rdb.FlushDB(redis.Ctx)

		// Insert video id list to redis.
		videoAuthorKey := redis.VideosByAuthorKey + strconv.FormatInt(testUserOne.ID, 10)
		redis.Rdb.RPush(redis.Ctx, videoAuthorKey, testVideoOne.ID)

		statusCode, _, videoList :=
			videoService.GetVideoListByAuthorID(testUserOne.ID, 0)

		assert.Equal(t, int32(0), statusCode)
		assert.Equal(t, 1, len(videoList))
	})
}

func TestGetMost30Videos(t *testing.T) {
	models.InitDatabase(true)

	statusCode, statusMsg, _, videoList := videoService.GetMost30Videos(time.Now(), 0)

	assert.Equal(t, int32(0), statusCode)
	assert.Equal(t, "get most 30 videos successfully", statusMsg)
	assert.Equal(t, 0, len(videoList))
}

func TestGetVideoListByVideoIDList(t *testing.T) {
	setup()

	t.Run("get video list successfully with cache miss", func(t *testing.T) {
		redis.Rdb.FlushDB(redis.Ctx)

		statusCode, statusMsg, videoList :=
			videoService.GetVideoListByVideoIDList([]int64{testVideoOne.ID}, 0)

		assert.Equal(t, int32(0), statusCode)
		assert.Equal(t, "get video list successfully", statusMsg)
		assert.Equal(t, 1, len(videoList))
	})

	t.Run("get video list successfully with cache hit", func(t *testing.T) {
		redis.Rdb.FlushDB(redis.Ctx)

		// Insert video to redis.
		testVideoOne := &models.VideoCache{ID: testVideoOne.ID, Title: testVideoOne.Title}
		videoKey := redis.VideoKey + strconv.FormatInt(testVideoOne.ID, 10)
		err := redis.Rdb.HSet(redis.Ctx, videoKey, testVideoOne).Err()
		if err != nil {
			t.Fatalf("Error when insert video to redis: %v", err)
		}

		statusCode, _, videoList :=
			videoService.GetVideoListByVideoIDList([]int64{testVideoOne.ID}, 0)

		assert.Equal(t, int32(0), statusCode)
		assert.Equal(t, 1, len(videoList))
	})
}
