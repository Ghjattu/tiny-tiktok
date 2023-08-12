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

	createTestVideo(1, time.Now(), "test")

	fs := &FavoriteService{}

	statusCode, statusMsg := fs.FavoriteAction(1, 1, 1)

	assert.Equal(t, int32(0), statusCode)
	assert.Equal(t, "favorite action success", statusMsg)
}

func TestFavoriteActionWithRepetitiveActionTypeOne(t *testing.T) {
	models.InitDatabase(true)

	createTestVideo(1, time.Now(), "test")

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

	createTestVideo(1, time.Now(), "test")

	fs := &FavoriteService{}

	statusCode, statusMsg := fs.FavoriteAction(1, 1, 2)

	assert.Equal(t, int32(0), statusCode)
	assert.Equal(t, "unfavorite action success", statusMsg)
}
