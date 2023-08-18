package interfaces

import "github.com/Ghjattu/tiny-tiktok/models"

type FavoriteInterface interface {
	// CreateNewFavoriteRel create a new favorite rel.
	// Return status code, status message.
	CreateNewFavoriteRel(userID, videoID int64) (int32, string)

	// DeleteFavoriteRel delete a favorite rel by user id and video id.
	// Return status code, status message.
	DeleteFavoriteRel(userID, videoID int64) (int32, string)

	// GetFavoriteVideoListByUserID get user's favorite video list by user id.
	// Return status_code, status_msg, video_list.
	GetFavoriteVideoListByUserID(currentUserID, queryUserID int64) (int32, string, []models.VideoDetail)

	// GetTotalFavoritedByUserID get total favorited count by user id.
	// Return count.
	GetTotalFavoritedByUserID(userID int64) int64
}
