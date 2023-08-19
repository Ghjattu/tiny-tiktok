package services

import (
	"strconv"
	"testing"
	"time"

	"github.com/Ghjattu/tiny-tiktok/middleware/redis"
	"github.com/Ghjattu/tiny-tiktok/models"
	"github.com/stretchr/testify/assert"
)

var (
	videoService = &VideoService{}
)

func TestCreateNewVideo(t *testing.T) {
	models.Flush()

	status_code, statue_msg := videoService.CreateNewVideo("test", "test", 1, time.Now())

	assert.Equal(t, int32(0), status_code)
	assert.Equal(t, "create new video successfully", statue_msg)
}

// func TestCreateNewVideoWithRedis(t *testing.T) {
// 	models.Flush()

// 	// Insert a test user to redis.
// 	testUser := &models.UserDetail{
// 		ID:        1,
// 		Name:      "test",
// 		WorkCount: 0,
// 	}
// 	userKey := redis.UserKey + "1"
// 	redis.Rdb.HSet(redis.Ctx, userKey, testUser)

// 	statusCode, statusMsg := videoService.CreateNewVideo("test", "test", 1, time.Now())
// 	workCount := redis.Rdb.HGet(redis.Ctx, userKey, "work_count").Val()

// 	assert.Equal(t, int32(0), statusCode)
// 	assert.Equal(t, "create new video successfully", statusMsg)
// 	assert.Equal(t, "1", workCount)
// }

func TestGetVideoListByAuthorID(t *testing.T) {
	models.Flush()

	// Create a new test user.
	testUser, _ := models.CreateTestUser("test", "123456")
	// Create a new test video.
	testVideo, _ := models.CreateTestVideo(testUser.ID, time.Now(), "test")
	// Create a test favorite relation.
	models.CreateTestFavoriteRel(testUser.ID+1, testVideo.ID)

	status_code, _, videoList := videoService.GetVideoListByAuthorID(testUser.ID, testUser.ID+1)

	assert.Equal(t, int32(0), status_code)
	assert.Equal(t, 1, len(videoList))
	assert.Equal(t, int64(1), videoList[0].FavoriteCount)
	assert.True(t, videoList[0].IsFavorite)
}

func TestGetVideoListByAuthorIDWithRedis(t *testing.T) {
	models.Flush()

	// Create a new test user.
	testUser, _ := models.CreateTestUser("test", "123456")
	// Create a test video.
	testVideo, _ := models.CreateTestVideo(testUser.ID, time.Now(), "test")
	// Insert video id list to redis.
	videoAuthorKey := redis.VideosByAuthorKey + strconv.FormatInt(testUser.ID, 10)
	redis.Rdb.RPush(redis.Ctx, videoAuthorKey, testVideo.ID)

	statusCode, _, videoList := videoService.GetVideoListByAuthorID(testUser.ID, testUser.ID+1)

	assert.Equal(t, int32(0), statusCode)
	assert.Equal(t, 1, len(videoList))
	assert.Equal(t, testVideo.Title, videoList[0].Title)
}

func TestGetMost30Videos(t *testing.T) {
	models.Flush()

	status_code, statue_msg, _, videoList := videoService.GetMost30Videos(time.Now(), 0)

	assert.Equal(t, int32(0), status_code)
	assert.Equal(t, "get most 30 videos successfully", statue_msg)
	assert.Equal(t, 0, len(videoList))
}

func TestGetVideoListByVideoIDList(t *testing.T) {
	models.Flush()

	// Create a new test user.
	testUser, _ := models.CreateTestUser("test", "123456")
	// Create a new test video.
	testVideo, _ := models.CreateTestVideo(testUser.ID, time.Now(), "test")

	status_code, statue_msg, videoList := videoService.GetVideoListByVideoIDList([]int64{testVideo.ID}, 1)

	assert.Equal(t, int32(0), status_code)
	assert.Equal(t, "get video list successfully", statue_msg)
	assert.Equal(t, 1, len(videoList))
}

func TestGetVideoListByVideoIDListWithRedis(t *testing.T) {
	models.Flush()

	// Create a test video.
	video, _ := models.CreateTestVideo(1, time.Now(), "test")
	testVideo := &models.VideoDetail{ID: video.ID, Title: video.Title}
	// Insert video to redis.
	videoKey := redis.VideoKey + strconv.FormatInt(testVideo.ID, 10)
	err := redis.Rdb.HSet(redis.Ctx, videoKey, testVideo).Err()
	if err != nil {
		t.Fatalf("Error when insert video to redis: %v", err)
	}

	statusCode, _, videoList := videoService.GetVideoListByVideoIDList([]int64{testVideo.ID}, 1)

	assert.Equal(t, int32(0), statusCode)
	assert.Equal(t, 1, len(videoList))
	assert.Equal(t, testVideo.Title, videoList[0].Title)
}
