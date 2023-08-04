package interfaces

import "github.com/Ghjattu/tiny-tiktok/models"

type VideoInterface interface {
	// CreateNewVideo creates a new video.
	// Return status_code, status_msg.
	CreateNewVideo(playUrl string, title string, authorID int64, authorName string) (int32, string)

	// GetPublishListByAuthorID returns a list of videos published by the author.
	// Return status_code, status_msg, video_list.
	GetPublishListByAuthorID(authorID int64) (int32, string, []models.VideoDetail)
}
