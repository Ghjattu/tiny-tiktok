package services

import (
	"strconv"
	"testing"
	"time"

	"github.com/Ghjattu/tiny-tiktok/middleware/redis"
	"github.com/stretchr/testify/assert"
)

var (
	commentService = &CommentService{}
)

func TestCreateNewComment(t *testing.T) {
	setup()

	t.Run("empty content", func(t *testing.T) {
		statusCode, _, _ := commentService.CreateNewComment(0, 0, "", time.Now())

		assert.Equal(t, int32(1), statusCode)
	})

	t.Run("video does not exist", func(t *testing.T) {
		statusCode, _, _ := commentService.CreateNewComment(0, 0, "test", time.Now())

		assert.Equal(t, int32(1), statusCode)
	})

	t.Run("create comment successfully", func(t *testing.T) {
		redis.Rdb.FlushDB(redis.Ctx)

		// Insert test video into cache.
		videoKey := redis.VideoKey + strconv.FormatInt(testVideoOne.ID, 10)
		redis.Rdb.HSet(redis.Ctx, videoKey, testVideoOneDetail)

		statusCode, _, _ :=
			commentService.CreateNewComment(0, testVideoOne.ID, "test", time.Now())
		commentCount := redis.Rdb.HGet(redis.Ctx, videoKey, "comment_count").Val()

		assert.Equal(t, int32(0), statusCode)
		assert.Equal(t, "1", commentCount)
	})
}

func TestDeleteCommentByCommentID(t *testing.T) {
	setup()

	t.Run("comment does not exist", func(t *testing.T) {
		statusCode, _, _ := commentService.DeleteCommentByCommentID(0, 0)

		assert.Equal(t, int32(1), statusCode)
	})

	t.Run("user does not have permission", func(t *testing.T) {
		statusCode, _, _ := commentService.DeleteCommentByCommentID(0, testCommentOne.ID)

		assert.Equal(t, int32(1), statusCode)
	})

	t.Run("delete comment successfully", func(t *testing.T) {
		redis.Rdb.FlushDB(redis.Ctx)

		// Insert test video into cache.
		videoKey := redis.VideoKey + strconv.FormatInt(testVideoOne.ID, 10)
		redis.Rdb.HSet(redis.Ctx, videoKey, testVideoOneDetail)

		commentService.CreateNewComment(0, testVideoOne.ID, "test", time.Now())

		statusCode, _, _ :=
			commentService.DeleteCommentByCommentID(testUserOne.ID, testCommentOne.ID)
		commentCount := redis.Rdb.HGet(redis.Ctx, videoKey, "comment_count").Val()

		assert.Equal(t, int32(0), statusCode)
		assert.Equal(t, "0", commentCount)
	})
}

func TestGetCommentListByVideoID(t *testing.T) {
	setup()

	t.Run("video does not exist", func(t *testing.T) {
		statusCode, _, _ := commentService.GetCommentListByVideoID(0, 0)

		assert.Equal(t, int32(1), statusCode)
	})

	t.Run("get comment list successfully", func(t *testing.T) {
		statusCode, _, commentList :=
			commentService.GetCommentListByVideoID(0, testVideoOne.ID)

		assert.Equal(t, int32(0), statusCode)
		assert.Equal(t, 1, len(commentList))
	})
}
