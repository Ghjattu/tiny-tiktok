package services

import (
	"testing"

	"github.com/Ghjattu/tiny-tiktok/models"
	"github.com/stretchr/testify/assert"
)

var (
	userService = &UserService{}
)

func TestGetUserByUserIDWithNonExistID(t *testing.T) {
	models.Flush()

	statusCode, statusMsg, _ := userService.GetUserByUserID(1)

	assert.Equal(t, int32(1), statusCode)
	assert.Equal(t, "user not found", statusMsg)
}

func TestGetUserByUserID(t *testing.T) {
	models.Flush()

	// Create a new test user.
	testUser, _ := models.CreateTestUser("test", "123456")

	statusCode, statusMsg, user := userService.GetUserByUserID(testUser.ID)

	assert.Equal(t, int32(0), statusCode)
	assert.Equal(t, "get user successfully", statusMsg)
	assert.Equal(t, testUser.Name, user.Name)
	assert.Equal(t, "", user.Password)
}

func TestGetUserDetailByUserIDWithNonExistID(t *testing.T) {
	models.Flush()

	statusCode, statusMsg, _ := userService.GetUserDetailByUserID(1, 2)

	assert.Equal(t, int32(1), statusCode)
	assert.Equal(t, "user not found", statusMsg)
}

func TestGetUserDetailByUserID(t *testing.T) {
	models.Flush()

	// Create a new test user.
	testUser, _ := models.CreateTestUser("test", "123456")

	statusCode, statusMsg, userDetail := userService.GetUserDetailByUserID(testUser.ID+1, testUser.ID)

	assert.Equal(t, int32(0), statusCode)
	assert.Equal(t, "get user successfully", statusMsg)
	assert.Equal(t, testUser.Name, userDetail.Name)
}

func TestGetUserDetailByUserIDFromCache(t *testing.T) {
	models.Flush()

	// Create a new test user.
	testUser, _ := models.CreateTestUser("test", "123456")

	// Get user detail from database.
	userService.GetUserDetailByUserID(testUser.ID+1, testUser.ID)

	// Get user detail from cache.
	statusCode, statusMsg, userDetail := userService.GetUserDetailByUserID(testUser.ID+1, testUser.ID)

	assert.Equal(t, int32(0), statusCode)
	assert.Equal(t, "get user successfully", statusMsg)
	assert.Equal(t, testUser.Name, userDetail.Name)
}
