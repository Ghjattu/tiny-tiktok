package interfaces

import "github.com/Ghjattu/tiny-tiktok/models"

type UserInterface interface {
	// GetUserByUserID gets a user by its user_id.
	// Return status_code, status_msg, user.
	GetUserByUserID(userID int64) (int32, string, *models.User)
}
