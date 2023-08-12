package services

import (
	"testing"
	"time"

	"github.com/Ghjattu/tiny-tiktok/models"
	"github.com/stretchr/testify/assert"
)

func TestFavoriteActionWithInvalidAction(t *testing.T) {
	models.InitDatabase(true)

	fs := &FavoriteService{}

	statusCode, statusMsg := fs.FavoriteAction(1, 1, 0)

	assert.Equal(t, int32(1), statusCode)
	assert.Equal(t, "action type is invalid", statusMsg)
}

func TestFavoriteActionWithNonExistVideoID(t *testing.T) {
	models.InitDatabase(true)

	fs := &FavoriteService{}

	statusCode, statusMsg := fs.FavoriteAction(1, 1, 1)

	assert.Equal(t, int32(1), statusCode)
	assert.Equal(t, "the video is not exist", statusMsg)
}

func TestFavoriteActionWithActionTypeOne(t *testing.T) {
	models.InitDatabase(true)

	// Create a new test video.
	models.CreateTestVideo(1, time.Now(), "test")

	fs := &FavoriteService{}

	statusCode, statusMsg := fs.FavoriteAction(1, 1, 1)

	assert.Equal(t, int32(0), statusCode)
	assert.Equal(t, "favorite action success", statusMsg)
}

func TestFavoriteActionWithRepetitiveActionTypeOne(t *testing.T) {
	models.InitDatabase(true)

	// Create a new test video.
	models.CreateTestVideo(1, time.Now(), "test")

	fs := &FavoriteService{}

	statusCode, statusMsg := fs.FavoriteAction(1, 1, 1)

	assert.Equal(t, int32(0), statusCode)
	assert.Equal(t, "favorite action success", statusMsg)

	statusCode, statusMsg = fs.FavoriteAction(1, 1, 1)

	assert.Equal(t, int32(1), statusCode)
	assert.Equal(t, "you have already favorited this video", statusMsg)
}

func TestFavoriteActionWithActionTypeTwo(t *testing.T) {
	models.InitDatabase(true)

	// Create a new test video.
	models.CreateTestVideo(1, time.Now(), "test")

	fs := &FavoriteService{}

	statusCode, statusMsg := fs.FavoriteAction(1, 1, 2)

	assert.Equal(t, int32(0), statusCode)
	assert.Equal(t, "unfavorite action success", statusMsg)
}

func TestGetFavoriteVideoListByUserID(t *testing.T) {
	models.InitDatabase(true)

	// Create a new test user.
	testUser, _ := models.CreateTestUser("test", "123456")

	// Create a new test video.
	testVideo, _ := models.CreateTestVideo(testUser.ID, time.Now(), "test")

	// Create a test favorite relation.
	models.CreateTestFavoriteRel(testUser.ID, testVideo.ID)

	fs := &FavoriteService{}

	statusCode, statusMsg, favoriteVideoList := fs.GetFavoriteVideoListByUserID(testUser.ID, testUser.ID)

	assert.Equal(t, int32(0), statusCode)
	assert.Equal(t, "get favorite video list successfully", statusMsg)
	assert.Equal(t, 1, len(favoriteVideoList))
}
