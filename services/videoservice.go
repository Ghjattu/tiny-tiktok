package services

import (
	"time"

	"github.com/Ghjattu/tiny-tiktok/models"
)

// VideoService implements VideoInterface.
type VideoService struct{}

func (vs *VideoService) CreateNewVideo(playUrl string, title string, authorID int64, publishTime time.Time) (int32, string) {
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

func (vs *VideoService) GetPublishListByAuthorID(authorID int64, currentUserID int64) (int32, string, []models.VideoDetail) {
	videoList, err := models.GetVideoListByUserID(authorID)
	if err != nil {
		return 1, "failed to get publish list", nil
	}

	// Update IsFavorite field and FavoriteCount field for each video.
	for i := 0; i < len(videoList); i++ {
		video := &videoList[i]
		video.IsFavorite, _ = models.CheckFavoriteRelExist(currentUserID, video.ID)
		video.FavoriteCount, _ = models.GetFavoriteCountByVideoID(video.ID)
	}

	return 0, "get publish list successfully", videoList
}

func (vs *VideoService) GetMost30Videos(latestTime time.Time) (int32, string, int64, []models.VideoDetail) {
	videoList, earliestTime, err := models.GetMost30Videos(latestTime)
	if err != nil {
		return 1, "failed to get most 30 videos", -1, nil
	}

	return 0, "get most 30 videos successfully", earliestTime.Unix(), videoList
}
