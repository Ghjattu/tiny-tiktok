package services

import (
	"github.com/Ghjattu/tiny-tiktok/models"
	"gorm.io/gorm"
)

// FollowService implements the FollowInterface.
type FollowService struct{}

// CreateNewFollowRel creates a new follow relationship.
//
//	@receiver fs *FollowService
//	@param followerID int64
//	@param followingID int64
//	@return int32 "status code"
//	@return string "status message"
func (fs *FollowService) CreateNewFollowRel(followerID, followingID int64) (int32, string) {
	// Check if the user exists.
	_, err := models.GetUserByUserID(followingID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return 1, "the user you want to follow does not exist"
		}
		return 1, "failed to check user existence"
	}

	// Check if the follow relationship exists.
	exist, _ := models.CheckFollowRelExist(followerID, followingID)
	if exist {
		return 1, "you have already followed this user"
	}

	// Create the follow relationship.
	fr := &models.FollowRel{
		FollowerID:  followerID,
		FollowingID: followingID,
	}
	_, err = models.CreateNewFollowRel(fr)
	if err != nil {
		return 1, "failed to create follow relationship"
	}

	return 0, "follow success"
}

// DeleteFollowRel delete a follow relationship by follower id and following id.
//
//	@receiver fs *FollowService
//	@param followerID int64
//	@param followingID int64
//	@return int32 "status code"
//	@return string "status message"
func (fs *FollowService) DeleteFollowRel(followerID, followingID int64) (int32, string) {
	// Delete the follow relationship.
	_, err := models.DeleteFollowRel(followerID, followingID)
	if err != nil {
		return 1, "failed to delete follow relationship"
	}

	return 0, "unfollow success"
}
