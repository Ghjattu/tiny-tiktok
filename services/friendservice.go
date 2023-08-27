package services

import "github.com/Ghjattu/tiny-tiktok/models"

// FriendService implements FriendInterface.
type FriendService struct{}

// GetFriendListByUserID get the friend list of a user.
//
//	@receiver fs *FriendService
//	@param currentUserID int64
//	@param queryUserID int64
//	@return int32 "status code"
//	@return string "status message"
//	@return []models.UserDetail "friend list"
func (fs *FriendService) GetFriendListByUserID(currentUserID, queryUserID int64) (int32, string, []models.UserDetail) {
	// Get the following id list of the query user.
	followingIDList, err := models.GetFollowingListByUserID(queryUserID)
	if err != nil {
		return 1, "failed to get following list of query user", nil
	}

	// Get the friend list of the query user.
	friendList := make([]models.UserDetail, 0)
	us := &UserService{}
	for _, followingID := range followingIDList {
		// Check if they follow each other.
		isFriend, _ := models.CheckFollowRelExist(followingID, queryUserID)
		if isFriend {
			statusCode, _, friend := us.GetUserDetailByUserID(currentUserID, followingID)
			if statusCode == 0 {
				friendList = append(friendList, *friend)
			}
		}
	}

	return 0, "get friend list successfully", friendList
}
