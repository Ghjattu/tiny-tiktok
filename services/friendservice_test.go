package services

import (
	"testing"

	"github.com/Ghjattu/tiny-tiktok/models"
	"github.com/stretchr/testify/assert"
)

var (
	friendService = &FriendService{}
)

func TestGetFriendListByUserID(t *testing.T) {
	setup()

	// Create a test following relationship.
	models.CreateTestFollowRel(followerUser.ID, followingUser.ID)
	models.CreateTestFollowRel(followingUser.ID, followerUser.ID)

	statusCode, _, friendList :=
		friendService.GetFriendListByUserID(0, followingUser.ID)

	assert.Equal(t, int32(0), statusCode)
	assert.Equal(t, 1, len(friendList))
	assert.Equal(t, followerUser.ID, friendList[0].ID)
}
