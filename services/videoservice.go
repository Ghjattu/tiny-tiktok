package services

import (
	"strconv"
	"time"

	"github.com/Ghjattu/tiny-tiktok/middleware/redis"
	"github.com/Ghjattu/tiny-tiktok/models"
)

// VideoService implements VideoInterface.
type VideoService struct{}

// CreateNewVideo creates a new video.
//
//	@receiver vs *VideoService
//	@param playUrl string
//	@param title string
//	@param authorID int64
//	@param publishTime  time.Time
//	@return int32 "status_code"
//	@return string "status_msg"
func (vs *VideoService) CreateNewVideo(playUrl, title string, authorID int64, publishTime time.Time) (int32, string) {
	video := &models.Video{
		AuthorID:    authorID,
		PublishTime: publishTime,
		PlayUrl:     playUrl,
		Title:       title,
	}

	// Update the WorkCount of the user in cache.
	userKey := redis.UserKey + strconv.FormatInt(authorID, 10)
	redis.HashIncrBy(userKey, "work_count", 1)

	// Insert new video to database.
	_, err := models.CreateNewVideo(video)
	if err != nil {
		return 1, "failed to create new video"
	}

	// If the video is created successfully, insert the video id to redis.
	videoAuthorKey := redis.VideosByAuthorKey + strconv.FormatInt(authorID, 10)
	redis.Rdb.RPush(redis.Ctx, videoAuthorKey, video.ID)

	return 0, "create new video successfully"
}

// GetVideoListByAuthorID returns a list of videos published by the author.
//
//	@receiver vs *VideoService
//	@param authorID int64
//	@param currentUserID int64
//	@return int32 "status_code"
//	@return string "status_msg"
//	@return []models.VideoDetail
func (vs *VideoService) GetVideoListByAuthorID(authorID, currentUserID int64) (int32, string, []models.VideoDetail) {
	// Try to get video id list from redis.
	videoAuthorKey := redis.VideosByAuthorKey + strconv.FormatInt(authorID, 10)
	if redis.Rdb.Exists(redis.Ctx, videoAuthorKey).Val() == 1 {
		// Cache hit.
		IDStrList, err := redis.Rdb.LRange(redis.Ctx, videoAuthorKey, 0, -1).Result()
		if err == nil {
			videoIDList, _ := redis.ConvertStringToInt64(IDStrList)

			// Update the expire time.
			redis.Rdb.Expire(redis.Ctx, videoAuthorKey, redis.RandomDay())

			return vs.GetVideoListByVideoIDList(videoIDList, currentUserID)
		}
	}

	// Cache miss or some error occurs.
	// Get video id list by author id.
	videoIDList, err := models.GetVideoIDListByAuthorID(authorID)
	if err != nil {
		return 1, "failed to get video list", nil
	}

	redis.Rdb.RPush(redis.Ctx, videoAuthorKey, videoIDList)
	redis.Rdb.Expire(redis.Ctx, videoAuthorKey, redis.RandomDay())

	return vs.GetVideoListByVideoIDList(videoIDList, currentUserID)
}

// GetMost30Videos returns the most 30 videos published before latestTime.
//
//	@receiver vs *VideoService
//	@param latestTime time.Time
//	@param currentUserID int64
//	@return int32 "status code"
//	@return string "status message"
//	@return int64 "the seconds of the earliest publish time of the returned video list"
//	@return []models.VideoDetail
func (vs *VideoService) GetMost30Videos(latestTime time.Time, currentUserID int64) (int32, string, int64, []models.VideoDetail) {
	videoIDList, earliestTime, err := models.GetMost30Videos(latestTime)
	if err != nil {
		return 1, "failed to get most 30 videos", -1, nil
	}

	_, _, videoDetailList := vs.GetVideoListByVideoIDList(videoIDList, currentUserID)

	return 0, "get most 30 videos successfully", earliestTime.Unix(), videoDetailList
}

// GetVideoListByVideoIDList returns a list of videos by video id list.
//
//	@receiver vs *VideoService
//	@param videoIDList []int64
//	@param currentUserID int64
//	@return int32 "status code"
//	@return string "status message"
//	@return []models.VideoDetail
func (vs *VideoService) GetVideoListByVideoIDList(videoIDList []int64, currentUserID int64) (int32, string, []models.VideoDetail) {
	videoDetailList := make([]models.VideoDetail, 0, len(videoIDList))

	for _, videoID := range videoIDList {
		// Try to get video from redis.
		videoKey := redis.VideoKey + strconv.FormatInt(videoID, 10)
		if redis.Rdb.Exists(redis.Ctx, videoKey).Val() == 1 {
			// Cache hit.
			videoDetail := &models.VideoDetail{}
			videoCache := redis.Rdb.HGetAll(redis.Ctx, videoKey)
			if videoCache.Err() == nil {
				if err := videoCache.Scan(videoDetail); err == nil {
					// Get the video's author.
					authorID, _ := models.GetAuthorIDByVideoID(videoID)
					us := &UserService{}
					_, _, videoDetail.Author = us.GetUserDetailByUserID(currentUserID, authorID)
					// Update the video's IsFavorite field.
					videoDetail.IsFavorite, _ = models.CheckFavoriteRelExist(currentUserID, videoID)

					videoDetailList = append(videoDetailList, *videoDetail)

					redis.Rdb.Expire(redis.Ctx, videoKey, redis.RandomDay())

					continue
				}
			}
		}

		// Cache miss or some error occurs.
		videoDetail, err := vs.GetVideoDetailByVideoID(videoID, currentUserID)
		if err == nil {
			videoDetailList = append(videoDetailList, *videoDetail)
		}
	}

	return 0, "get video list successfully", videoDetailList
}

// GetVideoDetailByVideoID returns the video detail by video id.
//
//	@receiver vs *VideoService
//	@param videoID int64
//	@param currentUserID int64
//	@return *models.VideoDetail
//	@return error
func (vs *VideoService) GetVideoDetailByVideoID(videoID, currentUserID int64) (*models.VideoDetail, error) {
	video, err := models.GetVideoByID(videoID)
	if err != nil {
		return nil, err
	}

	videoDetail := &models.VideoDetail{}
	videoDetail.ID = video.ID
	videoDetail.PlayUrl = video.PlayUrl
	videoDetail.CoverUrl = video.CoverUrl
	videoDetail.Title = video.Title

	// Get the video's author.
	us := &UserService{}
	_, _, videoDetail.Author = us.GetUserDetailByUserID(currentUserID, video.AuthorID)
	// Get the video's favorite count.
	videoDetail.FavoriteCount, _ = models.GetFavoriteCountByVideoID(video.ID)
	// Get the video's comment count.
	videoDetail.CommentCount, _ = models.GetCommentCountByVideoID(video.ID)
	// Update the video's IsFavorite field.
	videoDetail.IsFavorite, _ = models.CheckFavoriteRelExist(currentUserID, video.ID)

	// Insert the video to redis.
	videoKey := redis.VideoKey + strconv.FormatInt(videoID, 10)
	redis.Rdb.HSet(redis.Ctx, videoKey, videoDetail)
	redis.Rdb.Expire(redis.Ctx, videoKey, redis.RandomDay())

	return videoDetail, nil
}
