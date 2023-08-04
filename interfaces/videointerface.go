package interfaces

type VideoInterface interface {
	// CreateNewVideo creates a new video.
	// Return status_code, status_msg.
	CreateNewVideo(playUrl string, title string, authorID int64, authorName string) (int32, string)
}
