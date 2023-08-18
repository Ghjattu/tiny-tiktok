package controllers

import (
	"fmt"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/Ghjattu/tiny-tiktok/models"
	"github.com/stretchr/testify/assert"
)

func TestCommentActionWithActionTypeOne(t *testing.T) {
	models.Flush()
	// Register a test user.
	userID, _, token := registerTestUser("test", "123456")
	// Create a test video.
	testVideo, _ := models.CreateTestVideo(userID, time.Now(), "test")
	testVideoIDStr := fmt.Sprintf("%d", testVideo.ID)

	url := "http://127.0.0.1/douyin/comment/action/?video_id=" + testVideoIDStr +
		"&action_type=1&comment_text=abc&token=" + token
	req := httptest.NewRequest("POST", url, nil)

	w, r := sendRequest(req)
	res := r.(*CommentActionResponse)

	assert.Equal(t, 200, w.Code)
	assert.Equal(t, int32(0), res.StatusCode)
}

func TestCommentActionWithActionTypeTwo(t *testing.T) {
	models.Flush()
	// Register a test user.
	userID, _, token := registerTestUser("test", "123456")
	// Create a test video.
	testVideo, _ := models.CreateTestVideo(userID, time.Now(), "test")
	// Create a test comment.
	models.CreateTestComment(userID, testVideo.ID)

	url := "http://127.0.0.1/douyin/comment/action/?video_id=1&action_type=2&comment_id=1&token=" + token
	req := httptest.NewRequest("POST", url, nil)

	w, r := sendRequest(req)
	res := r.(*CommentActionResponse)

	assert.Equal(t, 200, w.Code)
	assert.Equal(t, int32(0), res.StatusCode)
}

func TestCommentList(t *testing.T) {
	models.Flush()

	// Register a test user.
	userID, _, token := registerTestUser("test", "123456")
	// Create a test video.
	testVideo, _ := models.CreateTestVideo(userID, time.Now(), "test")
	testVideoIDStr := fmt.Sprintf("%d", testVideo.ID)
	// Create a test comment.
	models.CreateTestComment(userID, testVideo.ID)

	url := "http://127.0.0.1/douyin/comment/list/?video_id=" + testVideoIDStr +
		"&token=" + token
	req := httptest.NewRequest("GET", url, nil)

	w, r := sendRequest(req)
	res := r.(*CommentListResponse)

	assert.Equal(t, 200, w.Code)
	assert.Equal(t, int32(0), res.StatusCode)
	assert.Equal(t, 1, len(res.CommentList))
}
