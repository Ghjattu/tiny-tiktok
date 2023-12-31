package services

import (
	"strconv"

	"github.com/Ghjattu/tiny-tiktok/models"
	"github.com/Ghjattu/tiny-tiktok/rabbitmq"
	"github.com/Ghjattu/tiny-tiktok/redis"
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

	// Create the follow relationship.
	fr := &models.FollowRel{
		FollowerID:  followerID,
		FollowingID: followingID,
	}
	_, err = models.CreateNewFollowRel(fr)
	if err != nil {
		return 1, "failed to create follow relationship"
	}

	// Update the FollowCount and FollowerCount of the user in cache.
	followerUserKey := redis.UserKey + strconv.FormatInt(followerID, 10)
	rabbitmq.ProduceMessage("Hash", "Incr", "", followerUserKey, "follow_count", 1)

	followingUserKey := redis.UserKey + strconv.FormatInt(followingID, 10)
	rabbitmq.ProduceMessage("Hash", "Incr", "", followingUserKey, "follower_count", 1)

	// Update the following id list of the user in cache.
	followingKey := redis.FollowingKey + strconv.FormatInt(followerID, 10)
	followingIDList := []int64{followingID}
	rabbitmq.ProduceMessage("List", "RPushX", "", followingKey, "", followingIDList)

	// Update the follower id list of the user in cache.
	followerKey := redis.FollowerKey + strconv.FormatInt(followingID, 10)
	followerIDList := []int64{followerID}
	rabbitmq.ProduceMessage("List", "RPushX", "", followerKey, "", followerIDList)

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

	// Delete the follow relationship.
	_, err = models.DeleteFollowRel(followerID, followingID)
	if err != nil {
		return 1, "failed to delete follow relationship"
	}

	// Update the FollowCount and FollowerCount of the user in cache.
	followerUserKey := redis.UserKey + strconv.FormatInt(followerID, 10)
	rabbitmq.ProduceMessage("Hash", "Incr", "", followerUserKey, "follow_count", -1)

	followingUserKey := redis.UserKey + strconv.FormatInt(followingID, 10)
	rabbitmq.ProduceMessage("Hash", "Incr", "", followingUserKey, "follower_count", -1)

	// Update the following id list of the user in cache.
	followingKey := redis.FollowingKey + strconv.FormatInt(followerID, 10)
	rabbitmq.ProduceMessage("List", "LRem", "", followingKey, "", followingID)

	// Update the follower id list of the user in cache.
	followerKey := redis.FollowerKey + strconv.FormatInt(followingID, 10)
	rabbitmq.ProduceMessage("List", "LRem", "", followerKey, "", followerID)

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

	// Try to get the following id list from cache.
	followingKey := redis.FollowingKey + strconv.FormatInt(queryUserID, 10)
	if redis.Rdb.Exists(redis.Ctx, followingKey).Val() == 1 {
		followingIDList, err := redis.Rdb.LRange(redis.Ctx, followingKey, 0, -1).Result()
		if err == nil {
			// Get the user detail list.
			us := &UserService{}
			userList := make([]models.UserDetail, 0, len(followingIDList))
			for _, followingID := range followingIDList {
				id, _ := strconv.ParseInt(followingID, 10, 64)
				statusCode, _, user := us.GetUserDetailByUserID(currentUserID, id)
				if statusCode == 0 {
					userList = append(userList, *user)
				}
			}

			redis.Rdb.Expire(redis.Ctx, followingKey, redis.RandomDay())

			return 0, "get following user list successfully", userList
		}
	}

	// Cache miss or some error occurs.
	// Get the following id list.
	followingIDList, err := models.GetFollowingListByUserID(queryUserID)
	if err != nil {
		return 1, "failed to get following list", nil
	}

	// Save the following id list to cache.
	rabbitmq.ProduceMessage("List", "RPush", "", followingKey, "", followingIDList)

	// Get the user detail list.
	us := &UserService{}
	userList := make([]models.UserDetail, 0, len(followingIDList))
	for _, followingID := range followingIDList {
		statusCode, _, user := us.GetUserDetailByUserID(currentUserID, followingID)
		if statusCode == 0 {
			userList = append(userList, *user)
		}
	}

	return 0, "get following user list successfully", userList
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

	// Try to get the following id list from cache.
	followerKey := redis.FollowerKey + strconv.FormatInt(queryUserID, 10)
	if redis.Rdb.Exists(redis.Ctx, followerKey).Val() == 1 {
		followerIDList, err := redis.Rdb.LRange(redis.Ctx, followerKey, 0, -1).Result()
		if err == nil {
			// Get the user detail list.
			us := &UserService{}
			userList := make([]models.UserDetail, 0, len(followerIDList))
			for _, followingID := range followerIDList {
				id, _ := strconv.ParseInt(followingID, 10, 64)
				statusCode, _, user := us.GetUserDetailByUserID(currentUserID, id)
				if statusCode == 0 {
					userList = append(userList, *user)
				}
			}

			redis.Rdb.Expire(redis.Ctx, followerKey, redis.RandomDay())

			return 0, "get follower user list successfully", userList
		}
	}

	// Get the follower id list.
	followerIDList, err := models.GetFollowerListByUserID(queryUserID)
	if err != nil {
		return 1, "failed to get follower list", nil
	}

	// Save the follower id list to cache.
	rabbitmq.ProduceMessage("List", "RPush", "", followerKey, "", followerIDList)

	// Get the user detail list.
	us := &UserService{}
	userList := make([]models.UserDetail, 0, len(followerIDList))
	for _, followerID := range followerIDList {
		statusCode, _, user := us.GetUserDetailByUserID(currentUserID, followerID)
		if statusCode == 0 {
			userList = append(userList, *user)
		}
	}

	return 0, "get follower user list success", userList
}

// GetFollowingCountByUserID get the number of users that a user is following.
//
//	@receiver fs *FollowService
//	@param userID int64
//	@return int64 "number of users that a user is following"
//	@return error
func (fs *FollowService) GetFollowingCountByUserID(userID int64) (int64, error) {
	followingKey := redis.FollowingKey + strconv.FormatInt(userID, 10)
	if redis.Rdb.Exists(redis.Ctx, followingKey).Val() == 1 {
		// Cache hit.
		return redis.Rdb.LLen(redis.Ctx, followingKey).Result()
	}

	return models.GetFollowingCountByUserID(userID)
}

// GetFollowerCountByUserID get the number of followers of a user.
//
//	@receiver fs *FollowService
//	@param userID int64
//	@return int64 "number of followers"
//	@return error
func (fs *FollowService) GetFollowerCountByUserID(userID int64) (int64, error) {
	followerKey := redis.FollowerKey + strconv.FormatInt(userID, 10)
	if redis.Rdb.Exists(redis.Ctx, followerKey).Val() == 1 {
		// Cache hit.
		return redis.Rdb.LLen(redis.Ctx, followerKey).Result()
	}

	return models.GetFollowerCountByUserID(userID)
}
