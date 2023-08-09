package controllers

import (
	"fmt"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/Ghjattu/tiny-tiktok/models"
	"github.com/stretchr/testify/assert"
)

var (
	validUserID int64
	validUser   *models.User
	validToken  string

	validVideoIDStr string
)

func init() {
	validUserID, validUser, validToken = registerTestUser("test", "123456")

	testVideo, _ := createTestVideo(validUserID, time.Now(), "test")
	validVideoIDStr = fmt.Sprintf("%d", testVideo.ID)
}

func TestFavoriteActionWithInvalidVideoID(t *testing.T) {
	url := "http://127.0.0.1/douyin/favorite/action/?video_id=abc&action_type=1&token=" + validToken
	req := httptest.NewRequest("POST", url, nil)

	w, r := sendRequest(req)
	res := r.(*Response)

	assert.Equal(t, 400, w.Code)
	assert.Equal(t, int32(1), res.StatusCode)
	assert.Equal(t, "invalid syntax", res.StatusMsg)
}

func TestFavoriteActionWithInvalidActionType(t *testing.T) {
	url := "http://127.0.0.1/douyin/favorite/action/?video_id=1&action_type=abc&token=" + validToken
	req := httptest.NewRequest("POST", url, nil)

	w, r := sendRequest(req)
	res := r.(*Response)

	assert.Equal(t, 400, w.Code)
	assert.Equal(t, int32(1), res.StatusCode)
	assert.Equal(t, "invalid syntax", res.StatusMsg)
}

func TestFavoriteActionWithValidVideoIDAndType(t *testing.T) {
	url := "http://127.0.0.1/douyin/favorite/action/?video_id=" + validVideoIDStr +
		"&action_type=1&token=" + validToken
	req := httptest.NewRequest("POST", url, nil)

	w, r := sendRequest(req)
	res := r.(*Response)

	assert.Equal(t, 200, w.Code)
	assert.Equal(t, int32(0), res.StatusCode)
	assert.Equal(t, "favorite action success", res.StatusMsg)
}
