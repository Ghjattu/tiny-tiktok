package services

import (
	"strconv"

	"github.com/Ghjattu/tiny-tiktok/models"
	"github.com/Ghjattu/tiny-tiktok/rabbitmq"
	"github.com/Ghjattu/tiny-tiktok/redis"
	"github.com/Ghjattu/tiny-tiktok/utils"
	"gorm.io/gorm"
)

// FavoriteService implements FavoriteInterface.
type FavoriteService struct{}

func (fs *FavoriteService) CreateNewFavoriteRel(userID, videoID int64) (int32, string) {
	// Check if the video exist.
	video, err := models.GetVideoByID(videoID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return 1, "the video is not exist"
		}
		return 1, "get video by id failed"
	}

	// Check if the favorite relation exist.
	exist, err := models.CheckFavoriteRelExist(userID, videoID)
	if err != nil {
		return 1, "check favorite rel exist failed"
	}
	if exist {
		return 1, "you have already favorited this video"
	}

	// Create a new favorite relation.
	fr := &models.FavoriteRel{
		UserID:  userID,
		VideoID: videoID,
	}

	_, err = models.CreateNewFavoriteRel(fr)
	if err != nil {
		return 1, "favorite action failed"
	}

	// Update the TotalFavorited of the video's author in cache.
	authorKey := redis.UserKey + strconv.FormatInt(video.AuthorID, 10)
	// redis.HashIncrBy(authorKey, "total_favorited", 1)
	rabbitmq.ProduceMessage("Hash", "Incr", "", authorKey, "total_favorited", 1)

	// Update the FavoriteCount of the user in cache.
	userKey := redis.UserKey + strconv.FormatInt(userID, 10)
	// redis.HashIncrBy(userKey, "favorite_count", 1)
	rabbitmq.ProduceMessage("Hash", "Incr", "", userKey, "favorite_count", 1)

	// Update the FavoriteCount of the video in cache.
	videoKey := redis.VideoKey + strconv.FormatInt(videoID, 10)
	// redis.HashIncrBy(videoKey, "favorite_count", 1)
	rabbitmq.ProduceMessage("Hash", "Incr", "", videoKey, "favorite_count", 1)

	// Update the favorite videos id list of the user in cache.
	favoriteVideosKey := redis.FavoriteVideosKey + strconv.FormatInt(userID, 10)
	// redis.Rdb.RPush(redis.Ctx, favoriteVideosKey, videoID)
	videoIDList := []int64{videoID}
	rabbitmq.ProduceMessage("List", "RPushX", "", favoriteVideosKey, "", videoIDList)

	return 0, "favorite action success"
}

func (fs *FavoriteService) DeleteFavoriteRel(userID, videoID int64) (int32, string) {
	// Check if the video exist.
	video, err := models.GetVideoByID(videoID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return 1, "the video is not exist"
		}
		return 1, "failed to check if the video exist"
	}

	// Check if the favorite relation exist.
	exist, err := models.CheckFavoriteRelExist(userID, videoID)
	if err != nil {
		return 1, "failed to check if the favorite relation exist"
	}
	if !exist {
		return 1, "you have not favorited this video"
	}

	_, err = models.DeleteFavoriteRel(userID, videoID)
	if err != nil {
		return 1, "unfavorite action failed"
	}

	// Update the TotalFavorited of the video's author in cache.
	authorKey := redis.UserKey + strconv.FormatInt(video.AuthorID, 10)
	rabbitmq.ProduceMessage("Hash", "Incr", "", authorKey, "total_favorited", -1)

	// Update the FavoriteCount of the user in cache.
	userKey := redis.UserKey + strconv.FormatInt(userID, 10)
	rabbitmq.ProduceMessage("Hash", "Incr", "", userKey, "favorite_count", -1)

	// Update the FavoriteCount of the video in cache.
	videoKey := redis.VideoKey + strconv.FormatInt(videoID, 10)
	rabbitmq.ProduceMessage("Hash", "Incr", "", videoKey, "favorite_count", -1)

	// Update the favorite videos id list of the user in cache.
	favoriteVideosKey := redis.FavoriteVideosKey + strconv.FormatInt(userID, 10)
	rabbitmq.ProduceMessage("List", "LRem", "", favoriteVideosKey, "", videoID)

	return 0, "unfavorite action success"
}

func (fs *FavoriteService) GetFavoriteVideoListByUserID(currentUserID, queryUserID int64) (int32, string, []models.VideoDetail) {
	vs := &VideoService{}

	// Try to get favorite video id list from redis.
	favoriteVideosKey := redis.FavoriteVideosKey + strconv.FormatInt(queryUserID, 10)
	if redis.Rdb.Exists(redis.Ctx, favoriteVideosKey).Val() == 1 {
		// Cache hit.
		IDStrList, err := redis.Rdb.LRange(redis.Ctx, favoriteVideosKey, 0, -1).Result()
		if err == nil {
			videoIDList, _ := utils.ConvertStringToInt64(IDStrList)

			// Update the expire time.
			redis.Rdb.Expire(redis.Ctx, favoriteVideosKey, redis.RandomDay())

			return vs.GetVideoListByVideoIDList(videoIDList, currentUserID)
		}
	}

	// Cache miss or some error occurs.
	// Get favorite video id list by user id.
	favoriteVideoIDList, err := models.GetFavoriteVideoIDListByUserID(queryUserID)
	if err != nil {
		return 1, "failed to get favorite video id list", nil
	}

	// Save favorite video id list to redis.
	rabbitmq.ProduceMessage("List", "RPush", "", favoriteVideosKey, "", favoriteVideoIDList)

	// Get favorite video list by video id list.
	return vs.GetVideoListByVideoIDList(favoriteVideoIDList, currentUserID)
}

// GetTotalFavoritedByUserID returns the total number of favorited by user id.
//
//	@receiver fs FavoriteService
//	@param userID int64
//	@return int64
func (fs *FavoriteService) GetTotalFavoritedByUserID(userID int64) int64 {
	// Try to get total favorited from redis.
	userKey := redis.UserKey + strconv.FormatInt(userID, 10)
	if redis.Rdb.Exists(redis.Ctx, userKey).Val() == 1 {
		// Cache hit.
		totalFavorited, err := redis.Rdb.HGet(redis.Ctx, userKey, "total_favorited").Int64()
		if err == nil {
			// Update the expire time.
			redis.Rdb.Expire(redis.Ctx, userKey, redis.RandomDay())

			return totalFavorited
		}
	}

	// Cache miss or some error occurs.
	totalFavorited := int64(0)

	videoIDList, _ := models.GetVideoIDListByAuthorID(userID)
	for _, videoID := range videoIDList {
		count, _ := models.GetFavoriteCountByVideoID(videoID)
		totalFavorited += count
	}

	return totalFavorited
}

// GetFavoriteCountByUserID returns the number of favorite by user id.
//
//	@receiver fs FavoriteService
//	@param userID int64
//	@return int64 "favorite count"
//	@return error
func (fs *FavoriteService) GetFavoriteCountByUserID(userID int64) (int64, error) {
	favoriteVideoKey := redis.FavoriteVideosKey + strconv.FormatInt(userID, 10)
	if redis.Rdb.Exists(redis.Ctx, favoriteVideoKey).Val() == 1 {
		// Cache hit.
		return redis.Rdb.LLen(redis.Ctx, favoriteVideoKey).Result()
	}

	return models.GetFavoriteCountByUserID(userID)
}
