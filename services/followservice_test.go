package services

import (
	"strconv"
	"testing"

	"github.com/Ghjattu/tiny-tiktok/models"
	"github.com/Ghjattu/tiny-tiktok/redis"
	"github.com/stretchr/testify/assert"
)

var (
	followService = &FollowService{}
)

func TestCreateNewFollowRel(t *testing.T) {
	setup()

	t.Run("same follower and following", func(t *testing.T) {
		statusCode, statusMsg :=
			followService.CreateNewFollowRel(followerUser.ID, followerUser.ID)

		assert.Equal(t, int32(1), statusCode)
		assert.Equal(t, "you can not follow yourself", statusMsg)
	})

	t.Run("user does not exist", func(t *testing.T) {
		statusCode, statusMsg :=
			followService.CreateNewFollowRel(followerUser.ID, followerUser.ID+100)

		assert.Equal(t, int32(1), statusCode)
		assert.Equal(t, "the user you want to follow does not exist", statusMsg)
	})

	t.Run("already followed", func(t *testing.T) {
		// Create a test follow relationship.
		models.CreateTestFollowRel(followerUser.ID+100, followerUser.ID)

		statusCode, statusMsg :=
			followService.CreateNewFollowRel(followerUser.ID+100, followerUser.ID)

		assert.Equal(t, int32(1), statusCode)
		assert.Equal(t, "you have already followed this user", statusMsg)
	})

	t.Run("follow successfully", func(t *testing.T) {
		redis.Rdb.FlushDB(redis.Ctx)

		// Insert two test users to redis.
		followerUserKey := redis.UserKey + strconv.FormatInt(followerUser.ID, 10)
		redis.Rdb.HSet(redis.Ctx, followerUserKey, followerUserDetail)
		followingUserKey := redis.UserKey + strconv.FormatInt(followingUser.ID, 10)
		redis.Rdb.HSet(redis.Ctx, followingUserKey, followingUserDetail)

		statusCode, statusMsg :=
			followService.CreateNewFollowRel(followerUser.ID, followingUser.ID)

		// Retrieve the follow count and follower count of the user from cache.
		followCount := redis.Rdb.HGet(redis.Ctx, followerUserKey, "follow_count").Val()
		followerCount := redis.Rdb.HGet(redis.Ctx, followingUserKey, "follower_count").Val()

		// Retrieve the following id list and follower id list of the user from cache.
		followingKey := redis.FollowingKey + strconv.FormatInt(followerUser.ID, 10)
		followingIDList := redis.Rdb.LRange(redis.Ctx, followingKey, 0, -1).Val()
		followerKey := redis.FollowerKey + strconv.FormatInt(followingUser.ID, 10)
		followerIDList := redis.Rdb.LRange(redis.Ctx, followerKey, 0, -1).Val()

		assert.Equal(t, int32(0), statusCode)
		assert.Equal(t, "follow success", statusMsg)
		assert.Equal(t, "1", followCount)
		assert.Equal(t, "1", followerCount)
		assert.Equal(t, 1, len(followingIDList))
		assert.Equal(t, strconv.FormatInt(followingUser.ID, 10), followingIDList[0])
		assert.Equal(t, 1, len(followerIDList))
		assert.Equal(t, strconv.FormatInt(followerUser.ID, 10), followerIDList[0])
	})
}

func TestDeleteFollowRel(t *testing.T) {
	setup()

	t.Run("follow relationship does not exist", func(t *testing.T) {
		statusCode, statusMsg :=
			followService.DeleteFollowRel(followerUser.ID, followingUser.ID)

		assert.Equal(t, int32(1), statusCode)
		assert.Equal(t, "you have not followed this user", statusMsg)
	})

	t.Run("unfollow successfully", func(t *testing.T) {
		redis.Rdb.FlushDB(redis.Ctx)

		// Insert two test users to redis.
		followerUserKey := redis.UserKey + strconv.FormatInt(followerUser.ID, 10)
		redis.Rdb.HSet(redis.Ctx, followerUserKey, followerUserDetail)
		followingUserKey := redis.UserKey + strconv.FormatInt(followingUser.ID, 10)
		redis.Rdb.HSet(redis.Ctx, followingUserKey, followingUserDetail)

		// Create a test follow relationship.
		followService.CreateNewFollowRel(followerUser.ID, followingUser.ID)

		statusCode, statusMsg := followService.DeleteFollowRel(followerUser.ID, followingUser.ID)

		// Retrieve the follow count and follower count of the user from cache.
		followCount := redis.Rdb.HGet(redis.Ctx, followerUserKey, "follow_count").Val()
		followerCount := redis.Rdb.HGet(redis.Ctx, followingUserKey, "follower_count").Val()

		// Retrieve the following id list and follower id list of the user from cache.
		followingKey := redis.FollowingKey + strconv.FormatInt(followerUser.ID, 10)
		followingIDList := redis.Rdb.LRange(redis.Ctx, followingKey, 0, -1).Val()
		followerKey := redis.FollowerKey + strconv.FormatInt(followingUser.ID, 10)
		followerIDList := redis.Rdb.LRange(redis.Ctx, followerKey, 0, -1).Val()

		assert.Equal(t, int32(0), statusCode)
		assert.Equal(t, "unfollow success", statusMsg)
		assert.Equal(t, "0", followCount)
		assert.Equal(t, "0", followerCount)
		assert.Equal(t, 0, len(followingIDList))
		assert.Equal(t, 0, len(followerIDList))
	})
}

func TestGetFollowingListByUserID(t *testing.T) {
	setup()

	// Create a test follow relationship.
	models.CreateTestFollowRel(followerUser.ID, followingUser.ID)

	t.Run("user does not exist", func(t *testing.T) {
		statusCode, statusMsg, _ :=
			followService.GetFollowingListByUserID(0, followerUser.ID+100)

		assert.Equal(t, int32(1), statusCode)
		assert.Equal(t, "the user you want to query does not exist", statusMsg)
	})

	t.Run("get following list successfully with cache miss", func(t *testing.T) {
		redis.Rdb.FlushDB(redis.Ctx)

		statusCode, _, userList :=
			followService.GetFollowingListByUserID(0, followerUser.ID)

		assert.Equal(t, int32(0), statusCode)
		assert.Equal(t, 1, len(userList))
		assert.Equal(t, followingUser.Name, userList[0].Name)
	})

	t.Run("get following list successfully with cache hit", func(t *testing.T) {
		redis.Rdb.FlushDB(redis.Ctx)

		// Insert following id list to redis.
		followingKey := redis.FollowingKey + strconv.FormatInt(followerUser.ID, 10)
		redis.Rdb.RPush(redis.Ctx, followingKey, followingUser.ID)

		statusCode, _, userList :=
			followService.GetFollowingListByUserID(0, followerUser.ID)

		assert.Equal(t, int32(0), statusCode)
		assert.Equal(t, 1, len(userList))
		assert.Equal(t, followingUser.Name, userList[0].Name)
	})
}

func TestGetFollowerListByUserID(t *testing.T) {
	setup()

	// Create a test follow relationship.
	models.CreateTestFollowRel(followerUser.ID, followingUser.ID)

	t.Run("user does not exist", func(t *testing.T) {
		statusCode, statusMsg, _ :=
			followService.GetFollowerListByUserID(0, followingUser.ID+100)

		assert.Equal(t, int32(1), statusCode)
		assert.Equal(t, "the user you want to query does not exist", statusMsg)
	})

	t.Run("get follower list successfully with cache miss", func(t *testing.T) {
		redis.Rdb.FlushDB(redis.Ctx)

		statusCode, _, userList :=
			followService.GetFollowerListByUserID(0, followingUser.ID)

		assert.Equal(t, int32(0), statusCode)
		assert.Equal(t, 1, len(userList))
		assert.Equal(t, followerUser.Name, userList[0].Name)
	})

	t.Run("get follower list successfully with cache hit", func(t *testing.T) {
		redis.Rdb.FlushDB(redis.Ctx)

		// Insert follower id list to redis.
		followerKey := redis.FollowerKey + strconv.FormatInt(followingUser.ID, 10)
		redis.Rdb.RPush(redis.Ctx, followerKey, followerUser.ID)

		statusCode, _, userList :=
			followService.GetFollowerListByUserID(0, followingUser.ID)

		assert.Equal(t, int32(0), statusCode)
		assert.Equal(t, 1, len(userList))
		assert.Equal(t, followerUser.Name, userList[0].Name)
	})
}
