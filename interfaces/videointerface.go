package interfaces

import (
	"time"

	"github.com/Ghjattu/tiny-tiktok/models"
)

type VideoInterface interface {
	// CreateNewVideo creates a new video.
	// Return status_code, status_msg.
	CreateNewVideo(playUrl string, title string, authorID int64, publishTime time.Time) (int32, string)

	// GetPublishListByAuthorID returns a list of videos published by the author.
	// Return status_code, status_msg, video_list.
	GetPublishListByAuthorID(authorID int64) (int32, string, []models.VideoDetail)

	// GetMost30Videos returns the most 30 videos published before latest_time.
	// Return status_code, status_msg, next_time, video_list.
	GetMost30Videos(latestTime time.Time) (int32, string, int64, []models.VideoDetail)
}
