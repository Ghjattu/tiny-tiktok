package interfaces

type FavoriteInterface interface {
	// FavoriteAction favorite or unfavorite a video by action type.
	// Return status_code, status_msg.
	FavoriteAction(userID int64, videoID int64, actionType int64) (int32, string)
}
