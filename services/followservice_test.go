package services

import (
	"testing"

	"github.com/Ghjattu/tiny-tiktok/models"
	"github.com/stretchr/testify/assert"
)

var (
	followService = &FollowService{}
)

func TestCreateNewFollowRelWithNonExistUser(t *testing.T) {
	models.InitDatabase(true)

	statusCode, statusMsg := followService.CreateNewFollowRel(1, 2)

	assert.Equal(t, int32(1), statusCode)
	assert.Equal(t, "the user you want to follow does not exist", statusMsg)
}

func TestCreateNewFollowRelWithExistRel(t *testing.T) {
	models.InitDatabase(true)

	// Create a test user.
	testUser, _ := models.CreateTestUser("test", "123456")
	// Create a test follow relationship.
	models.CreateTestFollowRel(testUser.ID+1, testUser.ID)

	statusCode, statusMsg := followService.CreateNewFollowRel(testUser.ID+1, testUser.ID)

	assert.Equal(t, int32(1), statusCode)
	assert.Equal(t, "you have already followed this user", statusMsg)
}

func TestCreateNewFollowRel(t *testing.T) {
	models.InitDatabase(true)

	// Create a test user.
	testUser, _ := models.CreateTestUser("test", "123456")

	statusCode, statusMsg := followService.CreateNewFollowRel(testUser.ID+1, testUser.ID)

	assert.Equal(t, int32(0), statusCode)
	assert.Equal(t, "follow success", statusMsg)
}

func TestDeleteFollowRel(t *testing.T) {
	models.InitDatabase(true)

	// Create a test follow relationship.
	models.CreateTestFollowRel(1, 2)

	statusCode, statusMsg := followService.DeleteFollowRel(1, 2)

	assert.Equal(t, int32(0), statusCode)
	assert.Equal(t, "unfollow success", statusMsg)
}
