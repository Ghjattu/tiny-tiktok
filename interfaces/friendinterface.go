package interfaces

import "github.com/Ghjattu/tiny-tiktok/models"

type FriendInterface interface {
	// GetFriendListByUserID get the friend list of a user.
	// Return status code, status message, friend list.
	GetFriendListByUserID(currentUserID, queryUserID int64) (int32, string, []models.UserDetail)
}
