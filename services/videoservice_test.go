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
		redis.Rdb.HSet(redis.Ctx, userKey, testUserOneCache)
		// Insert a test video id list to redis.
		videoAuthorKey := redis.VideosByAuthorKey + strconv.FormatInt(testUserOne.ID, 10)
		redis.Rdb.RPush(redis.Ctx, videoAuthorKey, "")

		statusCode, statusMsg :=
			videoService.CreateNewVideo("test", "test", testUserOne.ID, time.Now())
		waitForConsumer()
		workCount := redis.Rdb.HGet(redis.Ctx, userKey, "work_count").Val()
		videoIDListLength := redis.Rdb.LLen(redis.Ctx, videoAuthorKey).Val()

		assert.Equal(t, int32(0), statusCode)
		assert.Equal(t, "create new video successfully", statusMsg)
		assert.Equal(t, "1", workCount)
		assert.Equal(t, int64(2), videoIDListLength)
	})
}

func TestGetVideoListByAuthorID(t *testing.T) {
	setup()

	t.Run("get video list successfully with cache miss", func(t *testing.T) {
		redis.Rdb.FlushDB(redis.Ctx)

		statusCode, _, videoList :=
			videoService.GetVideoListByAuthorID(testUserOne.ID, 0)
		waitForConsumer()
		videoAuthorKey := redis.VideosByAuthorKey + strconv.FormatInt(testUserOne.ID, 10)
		videoIDListLength := redis.Rdb.LLen(redis.Ctx, videoAuthorKey).Val()

		assert.Equal(t, int32(0), statusCode)
		assert.Equal(t, 2, len(videoList))
		assert.Equal(t, int64(2), videoIDListLength)
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
		testVideoOne := &redis.VideoCache{ID: testVideoOne.ID, Title: testVideoOne.Title}
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

func TestGetVideoDetailByAuthorID(t *testing.T) {
	setup()

	t.Run("failed to get video", func(t *testing.T) {
		_, err := videoService.GetVideoDetailByVideoID(0, 0)

		assert.NotNil(t, err)
	})

	t.Run("get video detail successfully", func(t *testing.T) {
		redis.Rdb.FlushDB(redis.Ctx)

		video, err := videoService.GetVideoDetailByVideoID(testVideoOne.ID, 0)
		waitForConsumer()
		videoKey := redis.VideoKey + strconv.FormatInt(testVideoOne.ID, 10)
		videoTitle := redis.Rdb.HGet(redis.Ctx, videoKey, "title").Val()

		assert.Nil(t, err)
		assert.Equal(t, testVideoOne.Title, video.Title)
		assert.Equal(t, testVideoOne.Title, videoTitle)
	})
}

func TestGetVideoCountByAuthorID(t *testing.T) {
	setup()

	t.Run("get video count with cache miss", func(t *testing.T) {
		redis.Rdb.FlushDB(redis.Ctx)

		videoCount, _ := videoService.GetVideoCountByAuthorID(testUserOne.ID)

		assert.Equal(t, int64(2), videoCount)
	})

	t.Run("get video count with cache hit", func(t *testing.T) {
		redis.Rdb.FlushDB(redis.Ctx)

		videoAuthorKey := redis.VideosByAuthorKey + strconv.FormatInt(testUserOne.ID, 10)
		redis.Rdb.RPush(redis.Ctx, videoAuthorKey, testVideoOne.ID)
		redis.Rdb.RPush(redis.Ctx, videoAuthorKey, testVideoTwo.ID)

		videoCount, _ := videoService.GetVideoCountByAuthorID(testUserOne.ID)

		assert.Equal(t, int64(2), videoCount)
	})
}

func BenchmarkGetVideoListByAuthorID(b *testing.B) {
	setup()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		videoService.GetVideoListByAuthorID(testUserOne.ID, 0)
	}
}

func BenchmarkGetVideoListByVideoIDList(b *testing.B) {
	setup()

	videoIDList := []int64{testVideoOne.ID}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		videoService.GetVideoListByVideoIDList(videoIDList, 0)
	}
}

func BenchmarkGetVideoCountByAuthorID(b *testing.B) {
	setup()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		videoService.GetVideoCountByAuthorID(testUserOne.ID)
	}
}
