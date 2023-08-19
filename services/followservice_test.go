package services

import (
	"testing"

	"github.com/Ghjattu/tiny-tiktok/middleware/redis"
	"github.com/Ghjattu/tiny-tiktok/models"
	"github.com/stretchr/testify/assert"
)

var (
	followService = &FollowService{}
)

func TestCreateNewFollowRelWithSameFollowerAndFollowing(t *testing.T) {
	models.Flush()

	statusCode, statusMsg := followService.CreateNewFollowRel(1, 1)

	assert.Equal(t, int32(1), statusCode)
	assert.Equal(t, "you can not follow yourself", statusMsg)
}

func TestCreateNewFollowRelWithNonExistUser(t *testing.T) {
	models.Flush()

	statusCode, statusMsg := followService.CreateNewFollowRel(1, 2)

	assert.Equal(t, int32(1), statusCode)
	assert.Equal(t, "the user you want to follow does not exist", statusMsg)
}

func TestCreateNewFollowRelWithExistRel(t *testing.T) {
	models.Flush()

	// Create a test user.
	testUser, _ := models.CreateTestUser("test", "123456")
	// Create a test follow relationship.
	models.CreateTestFollowRel(testUser.ID+1, testUser.ID)

	statusCode, statusMsg := followService.CreateNewFollowRel(testUser.ID+1, testUser.ID)

	assert.Equal(t, int32(1), statusCode)
	assert.Equal(t, "you have already followed this user", statusMsg)
}

func TestCreateNewFollowRel(t *testing.T) {
	models.Flush()

	// Create a test user.
	testUser, _ := models.CreateTestUser("test", "123456")

	statusCode, statusMsg := followService.CreateNewFollowRel(testUser.ID+1, testUser.ID)

	assert.Equal(t, int32(0), statusCode)
	assert.Equal(t, "follow success", statusMsg)
}

func TestCreateNewFollowRelWithRedis(t *testing.T) {
	models.Flush()

	// Create a test user.
	models.CreateTestUser("test", "123456")
	// Insert two test users to redis.
	testUser1 := &models.UserDetail{ID: 1, Name: "test"}
	testUser2 := &models.UserDetail{ID: 2, Name: "test"}
	redis.Rdb.HSet(redis.Ctx, redis.UserKey+"1", testUser1)
	redis.Rdb.HSet(redis.Ctx, redis.UserKey+"2", testUser2)

	statusCode, statusMsg := followService.CreateNewFollowRel(2, 1)
	followCount := redis.Rdb.HGet(redis.Ctx, redis.UserKey+"2", "follow_count").Val()
	followerCount := redis.Rdb.HGet(redis.Ctx, redis.UserKey+"1", "follower_count").Val()

	assert.Equal(t, int32(0), statusCode)
	assert.Equal(t, "follow success", statusMsg)
	assert.Equal(t, "1", followCount)
	assert.Equal(t, "1", followerCount)
}

func TestDeleteFollowRel(t *testing.T) {
	models.Flush()

	// Create a test follow relationship.
	models.CreateTestFollowRel(1, 2)

	statusCode, statusMsg := followService.DeleteFollowRel(1, 2)

	assert.Equal(t, int32(0), statusCode)
	assert.Equal(t, "unfollow success", statusMsg)
}

func TestDeleteFollowRelWithRedis(t *testing.T) {
	models.Flush()

	// Create a test follow relationship.
	models.CreateTestFollowRel(1, 2)
	// Insert two test users to redis.
	testUser1 := &models.UserDetail{ID: 1, Name: "test", FollowCount: 1}
	testUser2 := &models.UserDetail{ID: 2, Name: "test", FollowerCount: 1}
	redis.Rdb.HSet(redis.Ctx, redis.UserKey+"1", testUser1)
	redis.Rdb.HSet(redis.Ctx, redis.UserKey+"2", testUser2)

	statusCode, statusMsg := followService.DeleteFollowRel(1, 2)
	followCount := redis.Rdb.HGet(redis.Ctx, redis.UserKey+"1", "follow_count").Val()
	followerCount := redis.Rdb.HGet(redis.Ctx, redis.UserKey+"2", "follower_count").Val()

	assert.Equal(t, int32(0), statusCode)
	assert.Equal(t, "unfollow success", statusMsg)
	assert.Equal(t, "0", followCount)
	assert.Equal(t, "0", followerCount)
}

func TestGetFollowingListByUserIDWithNonExistUser(t *testing.T) {
	models.Flush()

	statusCode, statusMsg, _ := followService.GetFollowingListByUserID(1, 2)

	assert.Equal(t, int32(1), statusCode)
	assert.Equal(t, "the user you want to query does not exist", statusMsg)
}

func TestGetFollowingListByUserID(t *testing.T) {
	models.Flush()

	// Create two test users.
	testUserOne, _ := models.CreateTestUser("test", "123456")
	testUserTwo, _ := models.CreateTestUser("test2", "123456")
	// Create a test follow relationship.
	models.CreateTestFollowRel(testUserOne.ID, testUserTwo.ID)

	statusCode, statusMsg, userList :=
		followService.GetFollowingListByUserID(testUserTwo.ID, testUserOne.ID)

	assert.Equal(t, int32(0), statusCode)
	assert.Equal(t, "get following list success", statusMsg)
	assert.Equal(t, 1, len(userList))
	assert.Equal(t, testUserTwo.ID, userList[0].ID)
}

func TestGetFollowerListByUserIDWithNonExistUser(t *testing.T) {
	models.Flush()

	statusCode, statusMsg, _ := followService.GetFollowerListByUserID(1, 2)

	assert.Equal(t, int32(1), statusCode)
	assert.Equal(t, "the user you want to query does not exist", statusMsg)
}

func TestGetFollowerListByUserID(t *testing.T) {
	models.Flush()

	// Create two test users.
	testUserOne, _ := models.CreateTestUser("test", "123456")
	testUserTwo, _ := models.CreateTestUser("test2", "123456")
	// Create a test follow relationship.
	models.CreateTestFollowRel(testUserOne.ID, testUserTwo.ID)

	statusCode, statusMsg, userList :=
		followService.GetFollowerListByUserID(testUserOne.ID, testUserTwo.ID)

	assert.Equal(t, int32(0), statusCode)
	assert.Equal(t, "get follower list success", statusMsg)
	assert.Equal(t, 1, len(userList))
	assert.Equal(t, testUserOne.ID, userList[0].ID)
}
