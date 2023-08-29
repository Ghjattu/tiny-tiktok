package services

import (
	"strconv"

	"github.com/Ghjattu/tiny-tiktok/models"
	"github.com/Ghjattu/tiny-tiktok/rabbitmq"
	"github.com/Ghjattu/tiny-tiktok/redis"
	"gorm.io/gorm"
)

// UserService implements UserInterface.
type UserService struct{}

// GetUserByUserID gets a user by its user id.
//
//	@receiver us *UserService
//	@param userID int64
//	@return int32 "status_code"
//	@return string "status_msg"
//	@return *models.User
func (us *UserService) GetUserByUserID(userID int64) (int32, string, *models.User) {
	user, err := models.GetUserByUserID(userID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return 1, "user not found", nil
		}
		return 1, "failed to get user", nil
	}

	// Hide user password.
	user.Password = ""

	return 0, "get user successfully", user
}

// GetUserDetailByUserID gets a user detail by its user id.
//
//	@receiver us *UserService
//	@param currentUserID int64
//	@param userID int64
//	@return int32 "status code"
//	@return string "status message"
//	@return *models.UserDetail
func (us *UserService) GetUserDetailByUserID(currentUserID, userID int64) (int32, string, *models.UserDetail) {
	// Try to get user detail from redis.
	userKey := redis.UserKey + strconv.FormatInt(userID, 10)
	result, err := redis.HashGetAll(userKey)
	if err == nil {
		// Cache hit.
		userCache := &redis.UserCache{}
		if err := result.Scan(userCache); err == nil {
			userDetail := &models.UserDetail{
				ID:              userCache.ID,
				Name:            userCache.Name,
				Avatar:          userCache.Avatar,
				BackgroundImage: userCache.BackgroundImage,
				Signature:       userCache.Signature,
				FollowCount:     userCache.FollowCount,
				FollowerCount:   userCache.FollowerCount,
				WorkCount:       userCache.WorkCount,
				FavoriteCount:   userCache.FavoriteCount,
				TotalFavorited:  userCache.TotalFavorited,
			}
			userDetail.IsFollow, _ = models.CheckFollowRelExist(currentUserID, userID)

			// Update expire time.
			redis.Rdb.Expire(redis.Ctx, userKey, redis.RandomDay())

			return 0, "get user successfully", userDetail
		}
	}

	// Cache miss or some error occurs.
	user, err := models.GetUserByUserID(userID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return 1, "user not found", &models.UserDetail{}
		}
		return 1, "failed to get user", &models.UserDetail{}
	}

	userDetail := &models.UserDetail{
		ID:              user.ID,
		Name:            user.Name,
		Avatar:          user.Avatar,
		BackgroundImage: user.BackgroundImage,
		Signature:       user.Signature,
	}

	followService := &FollowService{}
	userDetail.FollowCount, _ = followService.GetFollowingCountByUserID(user.ID)
	userDetail.FollowerCount, _ = followService.GetFollowerCountByUserID(user.ID)
	userDetail.IsFollow, _ = models.CheckFollowRelExist(currentUserID, user.ID)

	videoService := &VideoService{}
	userDetail.WorkCount, _ = videoService.GetVideoCountByAuthorID(user.ID)

	favoriteService := &FavoriteService{}
	userDetail.FavoriteCount, _ = favoriteService.GetFavoriteCountByUserID(user.ID)
	userDetail.TotalFavorited = favoriteService.GetTotalFavoritedByUserID(user.ID)

	// Save user detail to redis.
	rabbitmq.ProduceMessage("Hash", "Set", "UserCache", userKey, "", userDetail)

	return 0, "get user successfully", userDetail
}
