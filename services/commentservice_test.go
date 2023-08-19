package services

import (
	"testing"
	"time"

	"github.com/Ghjattu/tiny-tiktok/middleware/redis"
	"github.com/Ghjattu/tiny-tiktok/models"
	"github.com/stretchr/testify/assert"
)

var (
	commentService = &CommentService{}
)

func TestCreateNewCommentWithNonExistVideo(t *testing.T) {
	models.Flush()

	statusCode, statusMsg, _ := commentService.CreateNewComment(1, 1, "test", time.Now())

	assert.Equal(t, int32(1), statusCode)
	assert.Equal(t, "the video is not exist", statusMsg)
}

func TestCreateNewCommentWithEmptyContent(t *testing.T) {
	models.Flush()

	statusCode, statusMsg, _ := commentService.CreateNewComment(1, 1, "", time.Now())

	assert.Equal(t, int32(1), statusCode)
	assert.Equal(t, "comment text cannot be empty", statusMsg)
}

func TestCreateNewComment(t *testing.T) {
	models.Flush()

	// Create a test user.
	testUser, _ := models.CreateTestUser("test", "123456")
	// Create a new video.
	testVideo, _ := models.CreateTestVideo(1, time.Now(), "test")

	timestamp := time.Now()
	statusCode, statusMsg, commentDetail :=
		commentService.CreateNewComment(testUser.ID, testVideo.ID, "test", timestamp)

	assert.Equal(t, int32(0), statusCode)
	assert.Equal(t, "create new comment successfully", statusMsg)
	assert.Equal(t, "test", commentDetail.Content)
	assert.Equal(t, timestamp.Format("01-02"), commentDetail.CreateDate)
}

func TestCreateNewCommentWithRedis(t *testing.T) {
	models.Flush()

	// Create a test video.
	video, _ := models.CreateTestVideo(1, time.Now(), "test")
	testVideo := &models.VideoDetail{ID: video.ID, CommentCount: 0}
	// Insert the video into cache.
	redis.Rdb.HSet(redis.Ctx, redis.VideoKey+"1", testVideo)

	statusCode, _, _ := commentService.CreateNewComment(1, 1, "test", time.Now())
	commentCount := redis.Rdb.HGet(redis.Ctx, redis.VideoKey+"1", "comment_count").Val()

	assert.Equal(t, int32(0), statusCode)
	assert.Equal(t, "1", commentCount)
}

func TestDeleteCommentByCommentIDWithNonExistComment(t *testing.T) {
	models.Flush()

	statusCode, statusMsg, _ := commentService.DeleteCommentByCommentID(1, 1)

	assert.Equal(t, int32(1), statusCode)
	assert.Equal(t, "the comment is not exist", statusMsg)
}

func TestDeleteCommentByCommentIDWithNonAuthor(t *testing.T) {
	models.Flush()

	// Create a test video.
	testVideo, _ := models.CreateTestVideo(1, time.Now(), "test")
	// Create a test comment.
	testComment, _ := models.CreateTestComment(1, testVideo.ID)

	statusCode, statusMsg, _ := commentService.DeleteCommentByCommentID(2, testComment.ID)

	assert.Equal(t, int32(1), statusCode)
	assert.Equal(t, "you do not have permission to delete this comment", statusMsg)
}

func TestDeleteCommentByCommentID(t *testing.T) {
	models.Flush()

	// Create a test video.
	testVideo, _ := models.CreateTestVideo(1, time.Now(), "test")
	// Create a test comment.
	testComment, _ := models.CreateTestComment(1, testVideo.ID)

	statusCode, statusMsg, commentDetail :=
		commentService.DeleteCommentByCommentID(1, testComment.ID)

	assert.Equal(t, int32(0), statusCode)
	assert.Equal(t, "delete comment successfully", statusMsg)
	assert.Equal(t, int64(1), commentDetail.ID)
}

func TestDeleteCommentByCommentIDWithRedis(t *testing.T) {
	models.Flush()

	// Create a test video.
	video, _ := models.CreateTestVideo(1, time.Now(), "test")
	testVideo := &models.VideoDetail{ID: video.ID, CommentCount: 1}
	// Create a test comment.
	testComment, _ := models.CreateTestComment(1, video.ID)
	// Insert the video into cache.
	redis.Rdb.HSet(redis.Ctx, redis.VideoKey+"1", testVideo)

	statusCode, _, _ := commentService.DeleteCommentByCommentID(1, testComment.ID)
	commentCount := redis.Rdb.HGet(redis.Ctx, redis.VideoKey+"1", "comment_count").Val()

	assert.Equal(t, int32(0), statusCode)
	assert.Equal(t, "0", commentCount)
}

func TestGetCommentListByVideoIDWithNonExistVideo(t *testing.T) {
	models.Flush()

	statusCode, statusMsg, _ := commentService.GetCommentListByVideoID(1, 1)

	assert.Equal(t, int32(1), statusCode)
	assert.Equal(t, "the video is not exist", statusMsg)
}

func TestGetCommentListByVideoID(t *testing.T) {
	models.Flush()

	// Create a test user.
	testUser, _ := models.CreateTestUser("test", "123456")
	// Create a test video.
	testVideo, _ := models.CreateTestVideo(1, time.Now(), "test")
	// Create a test comment.
	models.CreateTestComment(1, testVideo.ID)

	statusCode, statusMsg, commentList :=
		commentService.GetCommentListByVideoID(testUser.ID, testVideo.ID)

	assert.Equal(t, int32(0), statusCode)
	assert.Equal(t, "get comment list successfully", statusMsg)
	assert.Equal(t, 1, len(commentList))
	assert.Equal(t, int64(1), commentList[0].ID)
}
