package services

import (
	"strconv"

	"github.com/Ghjattu/tiny-tiktok/middleware/redis"
	"github.com/Ghjattu/tiny-tiktok/models"
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
	if redis.Rdb.Exists(redis.Ctx, userKey).Val() == 1 {
		userDetail := &models.UserDetail{}
		userCache := redis.Rdb.HGetAll(redis.Ctx, userKey)
		if userCache.Err() == nil {
			if err := userCache.Scan(userDetail); err == nil {
				userDetail.IsFollow, _ = models.CheckFollowRelExist(currentUserID, userID)

				// Update expire time.
				redis.Rdb.Expire(redis.Ctx, userKey, redis.RandomDay())

				return 0, "get user successfully", userDetail
			}
		}
	}

	// Cache miss.
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

	userDetail.FollowCount, _ = models.GetFollowingCountByUserID(user.ID)
	userDetail.FollowerCount, _ = models.GetFollowerCountByUserID(user.ID)
	userDetail.IsFollow, _ = models.CheckFollowRelExist(currentUserID, user.ID)
	userDetail.WorkCount, _ = models.GetVideoCountByAuthorID(user.ID)
	userDetail.FavoriteCount, _ = models.GetFavoriteCountByUserID(user.ID)

	fs := &FavoriteService{}
	userDetail.TotalFavorited = fs.GetTotalFavoritedByUserID(user.ID)

	// Save user detail to redis.
	redis.Rdb.HSet(redis.Ctx, userKey, userDetail)
	redis.Rdb.Expire(redis.Ctx, userKey, redis.RandomDay())

	return 0, "get user successfully", userDetail
}
