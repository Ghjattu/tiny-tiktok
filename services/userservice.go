package services

import (
	"github.com/Ghjattu/tiny-tiktok/models"
	"gorm.io/gorm"
)

type UserService struct{}

// GetUserByUserID gets a user by its user id.
//
//	@receiver us *UserService
//	@param userID int64
//	@return int32 status_code
//	@return string status_msg
//	@return *models.User user
func (us *UserService) GetUserByUserID(userID int64) (int32, string, *models.User) {
	user, err := models.GetUserByUserID(userID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return 1, "user not found", nil
		}
		return 1, "failed to get user", nil
	}

	user.Password = ""
	return 0, "get user successfully", user
}
