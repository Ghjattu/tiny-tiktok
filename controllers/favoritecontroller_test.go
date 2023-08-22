package controllers

import (
	"fmt"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/Ghjattu/tiny-tiktok/models"
	"github.com/stretchr/testify/assert"
)

func TestFavoriteAction(t *testing.T) {
	setup()

	t.Run("invalid action type", func(t *testing.T) {
		url := "http://127.0.0.1/douyin/favorite/action/?action_type=3&token=" + token
		req := httptest.NewRequest("POST", url, nil)

		w, r := sendRequest(req)
		res := r.(*Response)

		assert.Equal(t, 200, w.Code)
		assert.Equal(t, int32(1), res.StatusCode)
		assert.Equal(t, "action type is invalid", res.StatusMsg)
	})

	t.Run("create favorite relationship successfully", func(t *testing.T) {
		// Create a new test video.
		testVideo, _ := models.CreateTestVideo(userID, time.Now(), "test")
		videoIDStr := fmt.Sprintf("%d", testVideo.ID)

		url := "http://127.0.0.1/douyin/favorite/action/?video_id=" + videoIDStr +
			"&action_type=1&token=" + token
		req := httptest.NewRequest("POST", url, nil)

		w, r := sendRequest(req)
		res := r.(*Response)

		assert.Equal(t, 200, w.Code)
		assert.Equal(t, int32(0), res.StatusCode)

	})

	t.Run("delete favorite relationship successfully", func(t *testing.T) {
		// Create a new test video.
		testVideo, _ := models.CreateTestVideo(1, time.Now(), "test")
		videoIDStr := fmt.Sprintf("%d", testVideo.ID)
		// Create a test favorite relationship.
		models.CreateTestFavoriteRel(userID, testVideo.ID)

		url := "http://127.0.0.1/douyin/favorite/action/?video_id=" + videoIDStr +
			"&action_type=2&token=" + token
		req := httptest.NewRequest("POST", url, nil)

		w, r := sendRequest(req)
		res := r.(*Response)

		assert.Equal(t, 200, w.Code)
		assert.Equal(t, int32(0), res.StatusCode)

	})
}

func TestGetFavoriteListByUserID(t *testing.T) {
	setup()

	// Create a new test video.
	testVideo, _ := models.CreateTestVideo(userID, time.Now(), "test")
	// Create a new test favorite relation.
	models.CreateTestFavoriteRel(userID, testVideo.ID)

	url := "http://127.0.0.1/douyin/favorite/list/?user_id=" + userIDStr + "&token=" + token
	req := httptest.NewRequest("GET", url, nil)

	w, r := sendRequest(req)
	res := r.(*FavoriteListResponse)

	assert.Equal(t, 200, w.Code)
	assert.Equal(t, int32(0), res.StatusCode)
	assert.Equal(t, 1, len(res.VideoList))
}
