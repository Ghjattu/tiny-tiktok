package services

import (
	"testing"
	"time"

	"github.com/Ghjattu/tiny-tiktok/models"
	"github.com/stretchr/testify/assert"
)

var (
	favoriteService = &FavoriteService{}
)

func TestCreateNewFavoriteRelWithNonExistVideo(t *testing.T) {
	models.InitDatabase(true)

	statusCode, statusMsg := favoriteService.CreateNewFavoriteRel(1, 1)

	assert.Equal(t, int32(1), statusCode)
	assert.Equal(t, "the video is not exist", statusMsg)
}

func TestCreateNewFavoriteRel(t *testing.T) {
	models.InitDatabase(true)

	// Create a new test video.
	models.CreateTestVideo(1, time.Now(), "test")

	statusCode, statusMsg := favoriteService.CreateNewFavoriteRel(1, 1)

	assert.Equal(t, int32(0), statusCode)
	assert.Equal(t, "favorite action success", statusMsg)
}

func TestCreateNewFavoriteRelWithRepetition(t *testing.T) {
	models.InitDatabase(true)

	// Create a new test video.
	models.CreateTestVideo(1, time.Now(), "test")

	statusCode, statusMsg := favoriteService.CreateNewFavoriteRel(1, 1)

	assert.Equal(t, int32(0), statusCode)
	assert.Equal(t, "favorite action success", statusMsg)

	statusCode, statusMsg = favoriteService.CreateNewFavoriteRel(1, 1)

	assert.Equal(t, int32(1), statusCode)
	assert.Equal(t, "you have already favorited this video", statusMsg)
}

func TestDeleteFavoriteRel(t *testing.T) {
	models.InitDatabase(true)

	// Create a new test video.
	models.CreateTestVideo(1, time.Now(), "test")

	statusCode, statusMsg := favoriteService.DeleteFavoriteRel(1, 1)

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

	statusCode, statusMsg, favoriteVideoList := favoriteService.GetFavoriteVideoListByUserID(testUser.ID, testUser.ID)

	assert.Equal(t, int32(0), statusCode)
	assert.Equal(t, "get favorite video list successfully", statusMsg)
	assert.Equal(t, 1, len(favoriteVideoList))
}

func TestGetTotalFavoritedByUserID(t *testing.T) {
	models.InitDatabase(true)

	// Create a test video.
	testVideo, _ := models.CreateTestVideo(1, time.Now(), "test")
	// Create a test favorite relationship.
	models.CreateTestFavoriteRel(1, testVideo.ID)

	count := favoriteService.GetTotalFavoritedByUserID(1)

	assert.Equal(t, int64(1), count)
}
