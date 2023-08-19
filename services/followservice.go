package services

import (
	"strconv"

	"github.com/Ghjattu/tiny-tiktok/middleware/redis"
	"github.com/Ghjattu/tiny-tiktok/models"
	"gorm.io/gorm"
)

// FollowService implements the FollowInterface.
type FollowService struct{}

// TODO: retrieve following and follower from redis.
// key = following:user_id
// key = follower:user_id

// CreateNewFollowRel creates a new follow relationship.
//
//	@receiver fs *FollowService
//	@param followerID int64
//	@param followingID int64
//	@return int32 "status code"
//	@return string "status message"
func (fs *FollowService) CreateNewFollowRel(followerID, followingID int64) (int32, string) {
	if followerID == followingID {
		return 1, "you can not follow yourself"
	}

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

	// Update the FollowCount and FollowerCount of the user in cache.
	userKey := redis.UserKey + strconv.FormatInt(followerID, 10)
	if redis.Rdb.Exists(redis.Ctx, userKey).Val() == 1 {
		redis.Rdb.HIncrBy(redis.Ctx, userKey, "follow_count", 1)
		redis.Rdb.Expire(redis.Ctx, userKey, redis.RandomDay())
	}
	userKey = redis.UserKey + strconv.FormatInt(followingID, 10)
	if redis.Rdb.Exists(redis.Ctx, userKey).Val() == 1 {
		redis.Rdb.HIncrBy(redis.Ctx, userKey, "follower_count", 1)
		redis.Rdb.Expire(redis.Ctx, userKey, redis.RandomDay())
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
	// Check if the follow relationship exists.
	exist, err := models.CheckFollowRelExist(followerID, followingID)
	if err != nil {
		return 1, "failed to check follow relationship existence"
	}
	if !exist {
		return 1, "you have not followed this user"
	}

	// Update the FollowCount and FollowerCount of the user in cache.
	userKey := redis.UserKey + strconv.FormatInt(followerID, 10)
	if redis.Rdb.Exists(redis.Ctx, userKey).Val() == 1 {
		redis.Rdb.HIncrBy(redis.Ctx, userKey, "follow_count", -1)
		redis.Rdb.Expire(redis.Ctx, userKey, redis.RandomDay())
	}
	userKey = redis.UserKey + strconv.FormatInt(followingID, 10)
	if redis.Rdb.Exists(redis.Ctx, userKey).Val() == 1 {
		redis.Rdb.HIncrBy(redis.Ctx, userKey, "follower_count", -1)
		redis.Rdb.Expire(redis.Ctx, userKey, redis.RandomDay())
	}

	// Delete the follow relationship.
	_, err = models.DeleteFollowRel(followerID, followingID)
	if err != nil {
		return 1, "failed to delete follow relationship"
	}

	return 0, "unfollow success"
}

// GetFollowingListByUserID get the list of users that a user is following.
//
//	@receiver fs *FollowService
//	@param currentUserID int64
//	@param queryUserID int64
//	@return int32 "status code"
//	@return string "status message"
//	@return []models.UserDetail
func (fs *FollowService) GetFollowingListByUserID(currentUserID, queryUserID int64) (int32, string, []models.UserDetail) {
	// Check if the user exists.
	_, err := models.GetUserByUserID(queryUserID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return 1, "the user you want to query does not exist", nil
		}
		return 1, "failed to check user existence", nil
	}

	// Get the following list.
	followingList, err := models.GetFollowingListByUserID(queryUserID)
	if err != nil {
		return 1, "failed to get following list", nil
	}

	// Get the user detail list.
	us := &UserService{}
	userList := make([]models.UserDetail, 0, len(followingList))
	for _, followingID := range followingList {
		statusCode, _, user := us.GetUserDetailByUserID(currentUserID, followingID)
		if statusCode == 0 {
			userList = append(userList, *user)
		}
	}

	return 0, "get following list success", userList
}

// GetFollowerListByUserID get the list of followers of a user.
//
//	@receiver fs *FollowService
//	@param currentUserID int64
//	@param queryUserID int64
//	@return int32 "status code"
//	@return string "status message"
//	@return []models.UserDetail
func (fs *FollowService) GetFollowerListByUserID(currentUserID, queryUserID int64) (int32, string, []models.UserDetail) {
	// Check if the user exists.
	_, err := models.GetUserByUserID(queryUserID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return 1, "the user you want to query does not exist", nil
		}
		return 1, "failed to check user existence", nil
	}

	// Get the follower list.
	followerList, err := models.GetFollowerListByUserID(queryUserID)
	if err != nil {
		return 1, "failed to get follower list", nil
	}

	// Get the user detail list.
	us := &UserService{}
	userList := make([]models.UserDetail, 0, len(followerList))
	for _, followerID := range followerList {
		statusCode, _, user := us.GetUserDetailByUserID(currentUserID, followerID)
		if statusCode == 0 {
			userList = append(userList, *user)
		}
	}

	return 0, "get follower list success", userList
}
