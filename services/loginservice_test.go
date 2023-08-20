package services

import (
	"testing"

	"github.com/Ghjattu/tiny-tiktok/models"
	"github.com/stretchr/testify/assert"
)

var (
	loginService = &LoginService{}
)

func TestLogin(t *testing.T) {
	models.InitDatabase(true)
	// Create a new test user.
	testUser, _ := models.CreateTestUser("test", "123456")

	t.Run("long username", func(t *testing.T) {
		userID, statusCode, statusMsg, _ := loginService.Login(
			"1234567890123456789012345678901234567890123456789012345678901234567890", "123456")

		assert.Equal(t, int64(-1), userID)
		assert.Equal(t, int32(1), statusCode)
		assert.Equal(t, "username or password is too long", statusMsg)
	})

	t.Run("long password", func(t *testing.T) {
		userID, statusCode, statusMsg, _ := loginService.Login("test",
			"1234567890123456789012345678901234567890123456789012345678901234567890")

		assert.Equal(t, int64(-1), userID)
		assert.Equal(t, int32(1), statusCode)
		assert.Equal(t, "username or password is too long", statusMsg)
	})

	t.Run("username does not exist", func(t *testing.T) {
		userID, statusCode, statusMsg, _ := loginService.Login("not_exist_name", "123456")

		assert.Equal(t, int64(-1), userID)
		assert.Equal(t, int32(1), statusCode)
		assert.Equal(t, "username not found", statusMsg)
	})

	t.Run("wrong password", func(t *testing.T) {
		userID, statusCode, statusMsg, _ := loginService.Login("test", "12345")

		assert.Equal(t, int64(-1), userID)
		assert.Equal(t, int32(1), statusCode)
		assert.Equal(t, "wrong password", statusMsg)
	})

	t.Run("correct password", func(t *testing.T) {
		userID, statusCode, statusMsg, _ := loginService.Login("test", "123456")

		assert.Equal(t, int64(1), userID)
		assert.Equal(t, int32(0), statusCode)
		assert.Equal(t, "login successfully", statusMsg)
		assert.Equal(t, testUser.ID, userID)
	})
}
