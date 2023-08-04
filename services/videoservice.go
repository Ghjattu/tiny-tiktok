package services

import "github.com/Ghjattu/tiny-tiktok/models"

// VideoService implements VideoInterface.
type VideoService struct{}

func (vs *VideoService) CreateNewVideo(playUrl string, title string, authorID int64, authorName string) (int32, string) {
	// Check title is empty or not.
	if title == "" {
		return 1, "video title is empty"
	}

	video := &models.Video{
		AuthorID:   authorID,
		AuthorName: authorName,
		PlayUrl:    playUrl,
		Title:      title,
	}

	// Insert new video to database.
	_, err := models.CreateNewVideo(video)
	if err != nil {
		return 1, "failed to create new video"
	}

	return 0, "create new video successfully"
}
