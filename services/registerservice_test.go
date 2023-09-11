package services

import (
	"testing"

	"github.com/Ghjattu/tiny-tiktok/bloomfilter"
	"github.com/Ghjattu/tiny-tiktok/models"
	"github.com/stretchr/testify/assert"
)

var (
	registerService = &RegisterService{}
)

func TestRegister(t *testing.T) {
	models.InitDatabase(true)
	bloomfilter.ClearAll()
	// Create a test user.
	testUser, _ := models.CreateTestUser("test", "123456")

	t.Run("empty username", func(t *testing.T) {
		UserID, statusCode, statusMsg, _ := registerService.Register("", "123456")

		assert.Equal(t, int64(-1), UserID)
		assert.Equal(t, int32(1), statusCode)
		assert.Equal(t, "invalid username or password", statusMsg)
	})

	t.Run("empty password", func(t *testing.T) {
		UserID, statusCode, statusMsg, _ := registerService.Register("test", "")

		assert.Equal(t, int64(-1), UserID)
		assert.Equal(t, int32(1), statusCode)
		assert.Equal(t, "invalid username or password", statusMsg)
	})

	t.Run("too short password", func(t *testing.T) {
		UserID, statusCode, statusMsg, _ := registerService.Register("test", "123")

		assert.Equal(t, int64(-1), UserID)
		assert.Equal(t, int32(1), statusCode)
		assert.Equal(t, "password is too short", statusMsg)
	})

	t.Run("too long username", func(t *testing.T) {
		UserID, statusCode, statusMsg, _ := registerService.Register(
			"1234567890123456789012345678901234567890123456789012345678901234567890", "123456")

		assert.Equal(t, int64(-1), UserID)
		assert.Equal(t, int32(1), statusCode)
		assert.Equal(t, "username or password is too long", statusMsg)
	})

	t.Run("too long password", func(t *testing.T) {
		UserID, statusCode, statusMsg, _ := registerService.Register("test",
			"12345678901234567890123456789012345678901234567890123456789012345678901234567890")

		assert.Equal(t, int64(-1), UserID)
		assert.Equal(t, int32(1), statusCode)
		assert.Equal(t, "username or password is too long", statusMsg)
	})

	t.Run("exceed 72 bytes password", func(t *testing.T) {
		UserID, statusCode, statusMsg, _ := registerService.Register(testUser.Name+"1",
			"密码密码密码密码密码密码密码密码密码密码密码密码密码密码密码")

		assert.Equal(t, int64(-1), UserID)
		assert.Equal(t, int32(1), statusCode)
		assert.Equal(t, "password length exceeds 72 bytes", statusMsg)
	})

	t.Run("registered username", func(t *testing.T) {
		UserID, statusCode, statusMsg, _ := registerService.Register(testUser.Name, "123456")

		assert.Equal(t, int64(-1), UserID)
		assert.Equal(t, int32(1), statusCode)
		assert.Equal(t, "the username has been registered", statusMsg)
	})

	t.Run("valid username and password", func(t *testing.T) {
		id, statusCode, statusMsg, _ := registerService.Register(testUser.Name+"1", "123456")
		exist := bloomfilter.CheckInt64Exist(bloomfilter.UserBloomFilter, id)

		assert.Equal(t, int32(0), statusCode)
		assert.Equal(t, "register successfully", statusMsg)
		assert.True(t, exist)
	})
}
