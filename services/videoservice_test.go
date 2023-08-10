package services

import (
	"testing"
	"time"

	"github.com/Ghjattu/tiny-tiktok/models"
	"github.com/stretchr/testify/assert"
)

func TestCreateNewVideoWithCorrectVideo(t *testing.T) {
	models.InitDatabase(true)

	vs := &VideoService{}

	status_code, statue_msg := vs.CreateNewVideo("test", "test", 1, time.Now())

	assert.Equal(t, int32(0), status_code)
	assert.Equal(t, "create new video successfully", statue_msg)
}

func TestGetPublishListByAuthorIDWithNonExistID(t *testing.T) {
	models.InitDatabase(true)

	vs := &VideoService{}

	status_code, statue_msg, videoList := vs.GetPublishListByAuthorID(1, -1)

	assert.Equal(t, int32(1), status_code)
	assert.Equal(t, "failed to get publish list", statue_msg)
	assert.Equal(t, 0, len(videoList))
}

func TestGetPublishListByAuthorIDWithCorrectID(t *testing.T) {
	models.InitDatabase(true)

	// Create a new test user.
	testUser, _ := createTestUser("test", "123456")

	// Create a new test video.
	testVideo, _ := createTestVideo(testUser.ID, time.Now(), "test")

	// Create a test favorite relation.
	fr := &models.FavoriteRel{
		UserID:  testUser.ID + 1,
		VideoID: testVideo.ID,
	}
	models.CreateNewFavoriteRel(fr)

	vs := &VideoService{}

	status_code, statue_msg, videoList := vs.GetPublishListByAuthorID(testUser.ID, testUser.ID+1)

	assert.Equal(t, int32(0), status_code)
	assert.Equal(t, "get publish list successfully", statue_msg)
	assert.Equal(t, 1, len(videoList))
	assert.Equal(t, testVideo.AuthorID, videoList[0].Author.ID)
	assert.Equal(t, testVideo.PlayUrl, videoList[0].PlayUrl)
	assert.Equal(t, testVideo.Title, videoList[0].Title)
	assert.Equal(t, int64(1), videoList[0].FavoriteCount)
	assert.True(t, videoList[0].IsFavorite)
}

func TestGetMost30Videos(t *testing.T) {
	models.InitDatabase(true)

	vs := &VideoService{}

	status_code, statue_msg, _, videoList := vs.GetMost30Videos(time.Now())

	assert.Equal(t, int32(0), status_code)
	assert.Equal(t, "get most 30 videos successfully", statue_msg)
	assert.Equal(t, 0, len(videoList))
}
