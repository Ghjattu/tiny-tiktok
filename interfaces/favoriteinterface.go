package interfaces

import "github.com/Ghjattu/tiny-tiktok/models"

type FavoriteInterface interface {
	// FavoriteAction favorite or unfavorite a video by action type.
	// Return status_code, status_msg.
	FavoriteAction(userID, videoID, actionType int64) (int32, string)

	// GetFavoriteVideoListByUserID get user's favorite video list by user id.
	// Return status_code, status_msg, video_list.
	GetFavoriteVideoListByUserID(currentUserID, queryUserID int64) (int32, string, []models.VideoDetail)

	// GetTotalFavoritedByUserID get total favorited count by user id.
	// Return count.
	GetTotalFavoritedByUserID(userID int64) int64
}
