package services

import (
	"testing"
	"time"

	"github.com/Ghjattu/tiny-tiktok/middleware/redis"
	"github.com/Ghjattu/tiny-tiktok/models"
	"github.com/stretchr/testify/assert"
)

var (
	favoriteService = &FavoriteService{}
)

func TestCreateNewFavoriteRelWithNonExistVideo(t *testing.T) {
	models.Flush()

	statusCode, statusMsg := favoriteService.CreateNewFavoriteRel(1, 1)

	assert.Equal(t, int32(1), statusCode)
	assert.Equal(t, "the video is not exist", statusMsg)
}

func TestCreateNewFavoriteRel(t *testing.T) {
	models.Flush()

	// Create a new test video.
	models.CreateTestVideo(1, time.Now(), "test")

	statusCode, statusMsg := favoriteService.CreateNewFavoriteRel(1, 1)

	assert.Equal(t, int32(0), statusCode)
	assert.Equal(t, "favorite action success", statusMsg)
}

func TestCreateNewFavoriteRelWithRedis(t *testing.T) {
	models.Flush()

	// Create a test video.
	video, _ := models.CreateTestVideo(1, time.Now(), "test")
	testVideo := &models.VideoDetail{ID: video.ID, Title: "test"}
	// Insert two test users and video to redis.
	testUser1 := &models.UserDetail{ID: 1, Name: "test"}
	testUser2 := &models.UserDetail{ID: 2, Name: "test"}
	redis.Rdb.HSet(redis.Ctx, redis.UserKey+"1", testUser1)
	redis.Rdb.HSet(redis.Ctx, redis.UserKey+"2", testUser2)
	redis.Rdb.HSet(redis.Ctx, redis.VideoKey+"1", testVideo)

	statusCode, statusMsg := favoriteService.CreateNewFavoriteRel(2, 1)
	favoriteCount := redis.Rdb.HGet(redis.Ctx, redis.UserKey+"2", "favorite_count").Val()
	totalFavorited := redis.Rdb.HGet(redis.Ctx, redis.UserKey+"1", "total_favorited").Val()
	videoFavoriteCount := redis.Rdb.HGet(redis.Ctx, redis.VideoKey+"1", "favorite_count").Val()
	favoriteVideosList := redis.Rdb.LRange(redis.Ctx, redis.FavoriteVideosKey+"2", 0, -1).Val()

	assert.Equal(t, int32(0), statusCode)
	assert.Equal(t, "favorite action success", statusMsg)
	assert.Equal(t, "1", favoriteCount)
	assert.Equal(t, "1", totalFavorited)
	assert.Equal(t, "1", videoFavoriteCount)
	assert.Equal(t, 1, len(favoriteVideosList))
}

func TestCreateNewFavoriteRelWithRepetition(t *testing.T) {
	models.Flush()

	// Create a new test video.
	models.CreateTestVideo(1, time.Now(), "test")

	statusCode, statusMsg := favoriteService.CreateNewFavoriteRel(1, 1)

	assert.Equal(t, int32(0), statusCode)
	assert.Equal(t, "favorite action success", statusMsg)

	statusCode, statusMsg = favoriteService.CreateNewFavoriteRel(1, 1)

	assert.Equal(t, int32(1), statusCode)
	assert.Equal(t, "you have already favorited this video", statusMsg)
}

func TestDeleteFavoriteRelWithNonExistVideo(t *testing.T) {
	models.Flush()

	statusCode, statusMsg := favoriteService.DeleteFavoriteRel(1, 1)

	assert.Equal(t, int32(1), statusCode)
	assert.Equal(t, "the video is not exist", statusMsg)
}

func TestDeleteFavoriteRel(t *testing.T) {
	models.Flush()

	// Create a new test video.
	models.CreateTestVideo(1, time.Now(), "test")
	// Create a test favorite relationship.
	models.CreateTestFavoriteRel(1, 1)

	statusCode, statusMsg := favoriteService.DeleteFavoriteRel(1, 1)

	assert.Equal(t, int32(0), statusCode)
	assert.Equal(t, "unfavorite action success", statusMsg)
}

func TestDeleteFavoriteRelWithRedis(t *testing.T) {
	models.Flush()

	// Create a test video.
	video, _ := models.CreateTestVideo(1, time.Now(), "test")
	// Create a test favorite relationship.
	models.CreateTestFavoriteRel(2, 1)
	// Insert video to redis.
	testVideo := &models.VideoDetail{ID: video.ID, Title: "test", FavoriteCount: 1}
	redis.Rdb.HSet(redis.Ctx, redis.VideoKey+"1", testVideo)
	// Insert two test users to redis.
	testUser1 := &models.UserDetail{ID: 1, Name: "test", TotalFavorited: 1}
	testUser2 := &models.UserDetail{ID: 2, Name: "test", FavoriteCount: 1}
	redis.Rdb.HSet(redis.Ctx, redis.UserKey+"1", testUser1)
	redis.Rdb.HSet(redis.Ctx, redis.UserKey+"2", testUser2)

	statusCode, statusMsg := favoriteService.DeleteFavoriteRel(2, 1)
	favoriteCount := redis.Rdb.HGet(redis.Ctx, redis.UserKey+"2", "favorite_count").Val()
	totalFavorited := redis.Rdb.HGet(redis.Ctx, redis.UserKey+"1", "total_favorited").Val()
	videoFavoriteCount := redis.Rdb.HGet(redis.Ctx, redis.VideoKey+"1", "favorite_count").Val()

	assert.Equal(t, int32(0), statusCode)
	assert.Equal(t, "unfavorite action success", statusMsg)
	assert.Equal(t, "0", favoriteCount)
	assert.Equal(t, "0", totalFavorited)
	assert.Equal(t, "0", videoFavoriteCount)
}

func TestGetFavoriteVideoListByUserID(t *testing.T) {
	models.Flush()

	// Create a new test user.
	testUser, _ := models.CreateTestUser("test", "123456")
	// Create a new test video.
	testVideo, _ := models.CreateTestVideo(testUser.ID, time.Now(), "test")
	// Create a test favorite relation.
	models.CreateTestFavoriteRel(testUser.ID, testVideo.ID)

	statusCode, _, favoriteVideoList := favoriteService.GetFavoriteVideoListByUserID(testUser.ID, testUser.ID)

	assert.Equal(t, int32(0), statusCode)
	assert.Equal(t, 1, len(favoriteVideoList))
}

func TestGetFavoriteVideoListByUserIDWithRedis(t *testing.T) {
	models.Flush()

	// Create a new test user.
	testUser, _ := models.CreateTestUser("test", "123456")
	// Create a new test video.
	testVideo, _ := models.CreateTestVideo(testUser.ID, time.Now(), "test")
	// Create a test favorite relation.
	models.CreateTestFavoriteRel(testUser.ID, testVideo.ID)
	// Insert video id to redis.
	redis.Rdb.RPush(redis.Ctx, redis.FavoriteVideosKey+"1", testVideo.ID)

	statusCode, _, favoriteVideoList := favoriteService.GetFavoriteVideoListByUserID(testUser.ID, testUser.ID)

	assert.Equal(t, int32(0), statusCode)
	assert.Equal(t, 1, len(favoriteVideoList))
}

func TestGetTotalFavoritedByUserID(t *testing.T) {
	models.Flush()

	// Create a test video.
	testVideo, _ := models.CreateTestVideo(1, time.Now(), "test")
	// Create a test favorite relationship.
	models.CreateTestFavoriteRel(1, testVideo.ID)

	count := favoriteService.GetTotalFavoritedByUserID(1)

	assert.Equal(t, int64(1), count)
}
