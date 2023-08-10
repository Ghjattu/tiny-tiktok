package controllers

import (
	"fmt"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/Ghjattu/tiny-tiktok/models"
	"github.com/stretchr/testify/assert"
)

func TestFavoriteActionWithInvalidVideoID(t *testing.T) {
	models.InitDatabase(true)

	// Register a new test user.
	_, _, token := registerTestUser("test", "123456")

	url := "http://127.0.0.1/douyin/favorite/action/?video_id=abc&action_type=1&token=" + token
	req := httptest.NewRequest("POST", url, nil)

	w, r := sendRequest(req)
	res := r.(*Response)

	assert.Equal(t, 400, w.Code)
	assert.Equal(t, int32(1), res.StatusCode)
	assert.Equal(t, "invalid syntax", res.StatusMsg)
}

func TestFavoriteActionWithInvalidActionType(t *testing.T) {
	models.InitDatabase(true)

	// Register a new test user.
	_, _, token := registerTestUser("test", "123456")

	url := "http://127.0.0.1/douyin/favorite/action/?video_id=1&action_type=abc&token=" + token
	req := httptest.NewRequest("POST", url, nil)

	w, r := sendRequest(req)
	res := r.(*Response)

	assert.Equal(t, 400, w.Code)
	assert.Equal(t, int32(1), res.StatusCode)
	assert.Equal(t, "invalid syntax", res.StatusMsg)
}

func TestFavoriteActionWithValidVideoIDAndType(t *testing.T) {
	models.InitDatabase(true)

	// Register a new test user.
	userID, _, token := registerTestUser("test", "123456")

	// Create a new test video.
	testVideo, _ := createTestVideo(userID, time.Now(), "test")
	videoIDStr := fmt.Sprintf("%d", testVideo.ID)

	url := "http://127.0.0.1/douyin/favorite/action/?video_id=" + videoIDStr +
		"&action_type=1&token=" + token
	req := httptest.NewRequest("POST", url, nil)

	w, r := sendRequest(req)
	res := r.(*Response)

	assert.Equal(t, 200, w.Code)
	assert.Equal(t, int32(0), res.StatusCode)
	assert.Equal(t, "favorite action success", res.StatusMsg)
}
