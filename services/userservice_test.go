package services

import (
	"testing"

	"github.com/Ghjattu/tiny-tiktok/models"
	"github.com/stretchr/testify/assert"
)

func TestGetUserByUserIDWithWrongID(t *testing.T) {
	models.InitDatabase(true)

	rs := &RegisterService{}
	userID, statusCode, statusMsg, _ := rs.Register("test", "123456")

	assert.Equal(t, int32(0), statusCode)
	assert.Equal(t, "register successfully", statusMsg)

	us := &UserService{}

	statusCode, statusMsg, user := us.GetUserByUserID(userID + 1)

	assert.Equal(t, int32(1), statusCode)
	assert.Equal(t, "user not found", statusMsg)
	assert.Equal(t, (*models.User)(nil), user)
}

func TestGetUserByUserIDWithCorrectID(t *testing.T) {
	models.InitDatabase(true)

	rs := &RegisterService{}
	userID, statusCode, statusMsg, _ := rs.Register("test", "123456")

	assert.Equal(t, int32(0), statusCode)
	assert.Equal(t, "register successfully", statusMsg)

	us := &UserService{}

	statusCode, statusMsg, user := us.GetUserByUserID(userID)

	assert.Equal(t, int32(0), statusCode)
	assert.Equal(t, "get user successfully", statusMsg)
	assert.Equal(t, "test", user.Name)
	assert.Equal(t, "", user.Password)
}
