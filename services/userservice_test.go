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
	models.InitDatabase(true)

	statusCode, statusMsg, _ := userService.GetUserByUserID(1)

	assert.Equal(t, int32(1), statusCode)
	assert.Equal(t, "user not found", statusMsg)
}

func TestGetUserByUserIDWithCorrectID(t *testing.T) {
	models.InitDatabase(true)

	// Create a new test user.
	testUser, _ := models.CreateTestUser("test", "123456")

	statusCode, statusMsg, user := userService.GetUserByUserID(testUser.ID)

	assert.Equal(t, int32(0), statusCode)
	assert.Equal(t, "get user successfully", statusMsg)
	assert.Equal(t, testUser.Name, user.Name)
	assert.Equal(t, "", user.Password)
}

func TestGetUserDetailByUserIDWithNonExistID(t *testing.T) {
	models.InitDatabase(true)

	statusCode, statusMsg, _ := userService.GetUserDetailByUserID(1)

	assert.Equal(t, int32(1), statusCode)
	assert.Equal(t, "user not found", statusMsg)
}

func TestGetUserDetailByUserIDWithCorrectID(t *testing.T) {
	models.InitDatabase(true)

	// Create a new test user.
	testUser, _ := models.CreateTestUser("test", "123456")

	statusCode, statusMsg, userDetail := userService.GetUserDetailByUserID(testUser.ID)

	assert.Equal(t, int32(0), statusCode)
	assert.Equal(t, "get user successfully", statusMsg)
	assert.Equal(t, testUser.Name, userDetail.Name)
}
