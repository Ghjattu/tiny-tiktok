package services

import (
	"testing"
	"time"

	"github.com/Ghjattu/tiny-tiktok/models"
	"github.com/stretchr/testify/assert"
)

var (
	videoService = &VideoService{}
)

func TestCreateNewVideo(t *testing.T) {
	models.InitDatabase(true)

	status_code, statue_msg := videoService.CreateNewVideo("test", "test", 1, time.Now())

	assert.Equal(t, int32(0), status_code)
	assert.Equal(t, "create new video successfully", statue_msg)
}

func TestGetPublishListByAuthorID(t *testing.T) {
	models.InitDatabase(true)

	// Create a new test user.
	testUser, _ := models.CreateTestUser("test", "123456")
	// Create a new test video.
	testVideo, _ := models.CreateTestVideo(testUser.ID, time.Now(), "test")
	// Create a test favorite relation.
	models.CreateTestFavoriteRel(testUser.ID+1, testVideo.ID)

	status_code, statue_msg, videoList := videoService.GetVideoListByAuthorID(testUser.ID, testUser.ID+1)

	assert.Equal(t, int32(0), status_code)
	assert.Equal(t, "get publish list successfully", statue_msg)
	assert.Equal(t, 1, len(videoList))
	assert.Equal(t, int64(1), videoList[0].FavoriteCount)
	assert.True(t, videoList[0].IsFavorite)
}

func TestGetMost30Videos(t *testing.T) {
	models.InitDatabase(true)

	status_code, statue_msg, _, videoList := videoService.GetMost30Videos(time.Now())

	assert.Equal(t, int32(0), status_code)
	assert.Equal(t, "get most 30 videos successfully", statue_msg)
	assert.Equal(t, 0, len(videoList))
}

func TestGetVideoListByVideoIDList(t *testing.T) {
	models.InitDatabase(true)

	// Create a new test user.
	testUser, _ := models.CreateTestUser("test", "123456")
	// Create a new test video.
	testVideo, _ := models.CreateTestVideo(testUser.ID, time.Now(), "test")

	status_code, statue_msg, videoList := videoService.GetVideoListByVideoIDList([]int64{testVideo.ID}, 1)

	assert.Equal(t, int32(0), status_code)
	assert.Equal(t, "get video list successfully", statue_msg)
	assert.Equal(t, 1, len(videoList))
}
