package services

import (
	"strconv"
	"testing"
	"time"

	"github.com/Ghjattu/tiny-tiktok/redis"
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
		redis.Rdb.HSet(redis.Ctx, videoKey, testVideoOneCache)
		// Create a comment id list.
		commentVideoKey := redis.CommentsByVideoKey + strconv.FormatInt(testVideoOne.ID, 10)
		redis.Rdb.RPush(redis.Ctx, commentVideoKey, "")

		statusCode, _, _ :=
			commentService.CreateNewComment(0, testVideoOne.ID, "test", time.Now())
		waitForConsumer()
		commentCount := redis.Rdb.HGet(redis.Ctx, videoKey, "comment_count").Val()
		commentListLength := redis.Rdb.LLen(redis.Ctx, commentVideoKey).Val()

		assert.Equal(t, int32(0), statusCode)
		assert.Equal(t, "1", commentCount)
		assert.Equal(t, int64(2), commentListLength)
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
		testVideoOneCache.CommentCount = 1
		redis.Rdb.HSet(redis.Ctx, videoKey, testVideoOneCache)
		// Create a comment id list.
		commentVideoKey := redis.CommentsByVideoKey + strconv.FormatInt(testVideoOne.ID, 10)
		redis.Rdb.RPush(redis.Ctx, commentVideoKey, testCommentOne.ID)

		statusCode, _, _ :=
			commentService.DeleteCommentByCommentID(testUserOne.ID, testCommentOne.ID)
		waitForConsumer()
		commentCount := redis.Rdb.HGet(redis.Ctx, videoKey, "comment_count").Val()
		commentIDListLength := redis.Rdb.LLen(redis.Ctx, commentVideoKey).Val()

		assert.Equal(t, int32(0), statusCode)
		assert.Equal(t, "0", commentCount)
		assert.Equal(t, int64(0), commentIDListLength)
	})
}

func TestGetCommentListByVideoID(t *testing.T) {
	setup()

	t.Run("video does not exist", func(t *testing.T) {
		statusCode, _, _ := commentService.GetCommentListByVideoID(0, 0)

		assert.Equal(t, int32(1), statusCode)
	})

	t.Run("get comment list successfully with cache miss", func(t *testing.T) {
		redis.Rdb.FlushDB(redis.Ctx)

		statusCode, _, commentList :=
			commentService.GetCommentListByVideoID(0, testVideoOne.ID)
		waitForConsumer()
		commentVideoKey := redis.CommentsByVideoKey + strconv.FormatInt(testVideoOne.ID, 10)
		commentIDListLength := redis.Rdb.LLen(redis.Ctx, commentVideoKey).Val()

		assert.Equal(t, int32(0), statusCode)
		assert.Equal(t, 1, len(commentList))
		assert.Equal(t, int64(1), commentIDListLength)
	})

	t.Run("get comment list successfully with cache hit", func(t *testing.T) {
		redis.Rdb.FlushDB(redis.Ctx)

		// Insert test video into cache.
		// videoKey := redis.VideoKey + strconv.FormatInt(testVideoOne.ID, 10)
		// testVideoOneCache.CommentCount = 1
		// redis.Rdb.HSet(redis.Ctx, videoKey, testVideoOneCache)
		// Create a comment id list.
		commentVideoKey := redis.CommentsByVideoKey + strconv.FormatInt(testVideoOne.ID, 10)
		redis.Rdb.RPush(redis.Ctx, commentVideoKey, testCommentOne.ID)

		statusCode, _, commentList :=
			commentService.GetCommentListByVideoID(0, testVideoOne.ID)

		assert.Equal(t, int32(0), statusCode)
		assert.Equal(t, 1, len(commentList))
	})
}

func TestGetCommentListByCommentIDList(t *testing.T) {
	setup()

	t.Run("get comment list successfully with cache miss", func(t *testing.T) {
		redis.Rdb.FlushDB(redis.Ctx)

		statusCode, _, commentList :=
			commentService.GetCommentListByCommentIDList(0, []int64{testCommentOne.ID})

		assert.Equal(t, int32(0), statusCode)
		assert.Equal(t, 1, len(commentList))
		assert.Equal(t, testCommentOne.Content, commentList[0].Content)
	})

	t.Run("get comment list successfully with cache hit", func(t *testing.T) {
		redis.Rdb.FlushDB(redis.Ctx)

		// Insert test comment into cache.
		commentKey := redis.CommentKey + strconv.FormatInt(testCommentOne.ID, 10)
		redis.Rdb.HSet(redis.Ctx, commentKey, testCommentOneCache)

		statusCode, _, commentList :=
			commentService.GetCommentListByCommentIDList(0, []int64{testCommentOne.ID})

		assert.Equal(t, int32(0), statusCode)
		assert.Equal(t, 1, len(commentList))
		assert.Equal(t, testCommentOne.Content, commentList[0].Content)
	})
}

func TestGetCommentDetailByCommentID(t *testing.T) {
	setup()

	t.Run("comment does not exist", func(t *testing.T) {
		commentDetail, err :=
			commentService.GetCommentDetailByCommentID(0, testCommentOne.ID+100)

		assert.Nil(t, commentDetail)
		assert.NotNil(t, err)
	})

	t.Run("get comment detail successfully", func(t *testing.T) {
		commentDetail, err :=
			commentService.GetCommentDetailByCommentID(0, testCommentOne.ID)
		waitForConsumer()
		commentKey := redis.CommentKey + strconv.FormatInt(testCommentOne.ID, 10)
		commentContent := redis.Rdb.HGet(redis.Ctx, commentKey, "content").Val()

		assert.Nil(t, err)
		assert.Equal(t, testCommentOne.Content, commentDetail.Content)
		assert.Equal(t, testCommentOne.Content, commentContent)
	})
}

func BenchmarkGetCommentListByVideoID(b *testing.B) {
	setup()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		commentService.GetCommentListByVideoID(0, testVideoOne.ID)
	}
}

func BenchmarkGetCommentListByCommentIDList(b *testing.B) {
	setup()

	commentIDList := []int64{testCommentOne.ID}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		commentService.GetCommentListByCommentIDList(0, commentIDList)
	}
}

func BenchmarkGetCommentDetailByCommentID(b *testing.B) {
	setup()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		commentService.GetCommentDetailByCommentID(0, testCommentOne.ID)
	}
}
