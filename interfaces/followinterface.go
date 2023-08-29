package interfaces

import "github.com/Ghjattu/tiny-tiktok/models"

type FollowInterface interface {
	// CreateNewFollowRel creates a new follow relationship.
	// Return status code, status message.
	CreateNewFollowRel(followerID, followingID int64) (int32, string)

	// DeleteFollowRel delete a follow relationship by follower id and following id.
	// Return status code, status message.
	DeleteFollowRel(followerID, followingID int64) (int32, string)

	// GetFollowingListByUserID get the list of users that a user is following.
	// Return status code, status message, user detail list.
	GetFollowingListByUserID(currentUserID, queryUserID int64) (int32, string, []models.UserDetail)

	// GetFollowerListByUserID get the list of followers of a user.
	// Return status code, status message, user detail list.
	GetFollowerListByUserID(currentUserID, queryUserID int64) (int32, string, []models.UserDetail)

	// GetFollowingCountByUserID get the number of users that a user is following.
	// Return following count, error.
	GetFollowingCountByUserID(userID int64) (int64, error)

	// GetFollowerCountByUserID get the number of followers of a user.
	// Return follower count, error.
	GetFollowerCountByUserID(userID int64) (int64, error)
}
