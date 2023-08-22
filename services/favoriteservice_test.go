package services

import (
	"strconv"
	"testing"

	"github.com/Ghjattu/tiny-tiktok/models"
	"github.com/Ghjattu/tiny-tiktok/redis"
	"github.com/stretchr/testify/assert"
)

var (
	favoriteService = &FavoriteService{}
)

func TestCreateNewFavoriteRel(t *testing.T) {
	setup()

	t.Run("video does not exist", func(t *testing.T) {
		statusCode, statusMsg :=
			favoriteService.CreateNewFavoriteRel(testUserOne.ID, testVideoOne.ID+100)

		assert.Equal(t, int32(1), statusCode)
		assert.Equal(t, "the video is not exist", statusMsg)
	})

	t.Run("favorite relation already exist", func(t *testing.T) {
		// Create a test favorite relation.
		models.CreateTestFavoriteRel(testUserOne.ID, testVideoOne.ID)

		statusCode, statusMsg :=
			favoriteService.CreateNewFavoriteRel(testUserOne.ID, testVideoOne.ID)

		assert.Equal(t, int32(1), statusCode)
		assert.Equal(t, "you have already favorited this video", statusMsg)
	})

	t.Run("favorite action success", func(t *testing.T) {
		redis.Rdb.FlushDB(redis.Ctx)

		// Insert test user and video to redis.
		userKey := redis.UserKey + strconv.FormatInt(testUserOne.ID, 10)
		redis.Rdb.HSet(redis.Ctx, userKey, testUserOneDetail)
		videoKey := redis.VideoKey + strconv.FormatInt(testVideoTwo.ID, 10)
		redis.Rdb.HSet(redis.Ctx, videoKey, testVideoTwoCache)

		statusCode, statusMsg :=
			favoriteService.CreateNewFavoriteRel(testUserOne.ID, testVideoTwo.ID)

		totalFavorited := redis.Rdb.HGet(redis.Ctx, userKey, "total_favorited").Val()
		userFavoriteCount := redis.Rdb.HGet(redis.Ctx, userKey, "favorite_count").Val()
		videoFavoriteCount := redis.Rdb.HGet(redis.Ctx, videoKey, "favorite_count").Val()
		favoriteVideoKey := redis.FavoriteVideosKey + strconv.FormatInt(testUserOne.ID, 10)
		favoriteVideoIDList := redis.Rdb.LRange(redis.Ctx, favoriteVideoKey, 0, -1).Val()

		assert.Equal(t, int32(0), statusCode)
		assert.Equal(t, "favorite action success", statusMsg)
		assert.Equal(t, "1", totalFavorited)
		assert.Equal(t, "1", userFavoriteCount)
		assert.Equal(t, "1", videoFavoriteCount)
		assert.Equal(t, 1, len(favoriteVideoIDList))
		assert.Equal(t, strconv.FormatInt(testVideoTwo.ID, 10), favoriteVideoIDList[0])
	})
}

func TestDeleteFavoriteRel(t *testing.T) {
	setup()

	t.Run("video does not exist", func(t *testing.T) {
		statusCode, statusMsg :=
			favoriteService.DeleteFavoriteRel(testUserOne.ID, testVideoOne.ID+100)

		assert.Equal(t, int32(1), statusCode)
		assert.Equal(t, "the video is not exist", statusMsg)
	})

	t.Run("favorite relation does not exist", func(t *testing.T) {
		statusCode, statusMsg :=
			favoriteService.DeleteFavoriteRel(testUserOne.ID, testVideoOne.ID)

		assert.Equal(t, int32(1), statusCode)
		assert.Equal(t, "you have not favorited this video", statusMsg)
	})

	t.Run("unfavorite action success", func(t *testing.T) {
		redis.Rdb.FlushDB(redis.Ctx)

		// Insert test user and video to redis.
		userKey := redis.UserKey + strconv.FormatInt(testUserOne.ID, 10)
		redis.Rdb.HSet(redis.Ctx, userKey, testUserOneDetail)
		videoKey := redis.VideoKey + strconv.FormatInt(testVideoTwo.ID, 10)
		redis.Rdb.HSet(redis.Ctx, videoKey, testVideoTwoCache)

		// Create a test favorite relationship.
		favoriteService.CreateNewFavoriteRel(testUserOne.ID, testVideoTwo.ID)

		statusCode, statusMsg :=
			favoriteService.DeleteFavoriteRel(testUserOne.ID, testVideoTwo.ID)
		totalFavorited := redis.Rdb.HGet(redis.Ctx, userKey, "total_favorited").Val()
		userFavoriteCount := redis.Rdb.HGet(redis.Ctx, userKey, "favorite_count").Val()
		videoFavoriteCount := redis.Rdb.HGet(redis.Ctx, videoKey, "favorite_count").Val()
		favoriteVideoKey := redis.FavoriteVideosKey + strconv.FormatInt(testUserOne.ID, 10)
		favoriteVideoIDList := redis.Rdb.LRange(redis.Ctx, favoriteVideoKey, 0, -1).Val()

		assert.Equal(t, int32(0), statusCode)
		assert.Equal(t, "unfavorite action success", statusMsg)
		assert.Equal(t, "0", totalFavorited)
		assert.Equal(t, "0", userFavoriteCount)
		assert.Equal(t, "0", videoFavoriteCount)
		assert.Equal(t, 0, len(favoriteVideoIDList))
	})
}

func TestGetFavoriteVideoListByUserID(t *testing.T) {
	setup()

	t.Run("get video list successfully with cache miss", func(t *testing.T) {
		redis.Rdb.FlushDB(redis.Ctx)

		// Create a test favorite relationship.
		models.CreateTestFavoriteRel(testUserOne.ID, testVideoOne.ID)

		statusCode, _, favoriteVideoList :=
			favoriteService.GetFavoriteVideoListByUserID(0, testUserOne.ID)

		assert.Equal(t, int32(0), statusCode)
		assert.Equal(t, 1, len(favoriteVideoList))
		assert.Equal(t, testVideoOne.Title, favoriteVideoList[0].Title)
	})

	t.Run("get video list successfully with cache hit", func(t *testing.T) {
		redis.Rdb.FlushDB(redis.Ctx)

		// Create a test favorite relationship.
		favoriteService.CreateNewFavoriteRel(testUserOne.ID, testVideoTwo.ID)

		statusCode, _, favoriteVideoList :=
			favoriteService.GetFavoriteVideoListByUserID(0, testUserOne.ID)

		assert.Equal(t, int32(0), statusCode)
		assert.Equal(t, 1, len(favoriteVideoList))
		assert.Equal(t, testVideoTwo.Title, favoriteVideoList[0].Title)
	})
}

func TestGetTotalFavoritedByUserID(t *testing.T) {
	setup()

	t.Run("get total favorited successfully with cache miss", func(t *testing.T) {
		redis.Rdb.FlushDB(redis.Ctx)

		// Create a test favorite relationship.
		models.CreateTestFavoriteRel(testUserOne.ID, testVideoOne.ID)

		totalFavorited := favoriteService.GetTotalFavoritedByUserID(testUserOne.ID)

		assert.Equal(t, int64(1), totalFavorited)
	})

	t.Run("get total favorited successfully with cache hit", func(t *testing.T) {
		redis.Rdb.FlushDB(redis.Ctx)

		// Insert test user to redis.
		userKey := redis.UserKey + strconv.FormatInt(testUserOne.ID, 10)
		redis.Rdb.HSet(redis.Ctx, userKey, testUserOneDetail)

		// Create a test favorite relationship.
		favoriteService.CreateNewFavoriteRel(testUserOne.ID, testVideoTwo.ID)

		totalFavorited := favoriteService.GetTotalFavoritedByUserID(testUserOne.ID)

		assert.Equal(t, int64(1), totalFavorited)
	})
}
