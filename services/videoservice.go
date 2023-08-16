package services

import (
	"time"

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

	// Insert new video to database.
	_, err := models.CreateNewVideo(video)
	if err != nil {
		return 1, "failed to create new video"
	}

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
	// Get video list by author id.
	videoList, err := models.GetVideoListByAuthorID(authorID)
	if err != nil {
		return 1, "failed to get publish list", nil
	}

	// Convert video list to video detail list.
	statusCode, videoDetailList := convertVideoToVideoDetail(videoList, currentUserID)
	if statusCode == 1 {
		return 1, "failed to get publish list", nil
	}

	return 0, "get publish list successfully", videoDetailList
}

// GetMost30Videos returns the most 30 videos published before latestTime.
//
//	@receiver vs *VideoService
//	@param latestTime time.Time
//	@return int32 "status_code"
//	@return string "status_msg"
//	@return int64 "the seconds of the earliest publish time of the returned video list"
//	@return []models.VideoDetail
func (vs *VideoService) GetMost30Videos(latestTime time.Time) (int32, string, int64, []models.VideoDetail) {
	videoList, earliestTime, err := models.GetMost30Videos(latestTime)
	if err != nil {
		return 1, "failed to get most 30 videos", -1, nil
	}

	// Convert video list to video detail list.
	statusCode, videoDetailList := convertVideoToVideoDetail(videoList, 0)
	if statusCode == 1 {
		return 1, "failed to get most 30 videos", -1, nil
	}

	return 0, "get most 30 videos successfully", earliestTime.Unix(), videoDetailList
}

// GetVideoListByVideoIDList returns a list of videos by video id list.
//
//	@receiver vs *VideoService
//	@param videoIDList []int64
//	@param currentUserID int64
//	@return int32 "status_code"
//	@return string "status_msg"
//	@return []models.VideoDetail
func (vs *VideoService) GetVideoListByVideoIDList(videoIDList []int64, currentUserID int64) (int32, string, []models.VideoDetail) {
	videoList := make([]models.Video, 0, len(videoIDList))

	for _, videoID := range videoIDList {
		video, err := models.GetVideoByID(videoID)
		if err == nil {
			videoList = append(videoList, *video)
		}
	}

	// Convert video list to video detail list.
	statusCode, videoDetailList := convertVideoToVideoDetail(videoList, currentUserID)
	if statusCode == 1 {
		return 1, "failed to get video list", nil
	}

	return 0, "get video list successfully", videoDetailList
}

// convertVideoToVideoDetail converts video list to video detail list.
//
//	@param videoList []models.Video
//	@param currentUserID int64
//	@return int32 "status_code"
//	@return []models.VideoDetail
func convertVideoToVideoDetail(videoList []models.Video, currentUserID int64) (int32, []models.VideoDetail) {
	// Initialize.
	us := &UserService{}
	videoDetail := &models.VideoDetail{}
	videoDetailList := make([]models.VideoDetail, 0, len(videoList))

	for _, video := range videoList {
		videoDetail.ID = video.ID
		videoDetail.PlayUrl = video.PlayUrl
		videoDetail.CoverUrl = video.CoverUrl
		videoDetail.Title = video.Title

		// Get the video's author.
		statusCode, _, author := us.GetUserDetailByUserID(currentUserID, video.AuthorID)
		if statusCode == 1 {
			return 1, nil
		}
		videoDetail.Author = author

		// Get the video's favorite count.
		count, err := models.GetFavoriteCountByVideoID(video.ID)
		if err != nil {
			return 1, nil
		}
		videoDetail.FavoriteCount = count

		// Get the video's comment count.
		count, err = models.GetCommentCountByVideoID(video.ID)
		if err != nil {
			return 1, nil
		}
		videoDetail.CommentCount = count

		// Update the video's IsFavorite field.
		isFavorite, err := models.CheckFavoriteRelExist(currentUserID, video.ID)
		if err != nil {
			return 1, nil
		}
		videoDetail.IsFavorite = isFavorite

		videoDetailList = append(videoDetailList, *videoDetail)
	}

	return 0, videoDetailList
}
