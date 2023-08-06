package services

import (
	"testing"

	"github.com/Ghjattu/tiny-tiktok/models"
	"github.com/stretchr/testify/assert"
)

func TestLoginWithLongUsername(t *testing.T) {
	models.InitDatabase(true)

	ls := &LoginService{}

	user_id, status_code, status_msg, _ := ls.Login(
		"1234567890123456789012345678901234567890123456789012345678901234567890", "123456")

	assert.Equal(t, int64(-1), user_id)
	assert.Equal(t, int32(1), status_code)
	assert.Equal(t, "username or password is too long", status_msg)
}

func TestLoginWithLongPassword(t *testing.T) {
	models.InitDatabase(true)

	ls := &LoginService{}

	user_id, status_code, status_msg, _ := ls.Login("test",
		"1234567890123456789012345678901234567890123456789012345678901234567890")

	assert.Equal(t, int64(-1), user_id)
	assert.Equal(t, int32(1), status_code)
	assert.Equal(t, "username or password is too long", status_msg)
}

func TestLoginWithNotExistName(t *testing.T) {
	models.InitDatabase(true)

	ls := &LoginService{}

	user_id, status_code, status_msg, _ := ls.Login("test", "123456")

	assert.Equal(t, int64(-1), user_id)
	assert.Equal(t, int32(1), status_code)
	assert.Equal(t, "username not found", status_msg)
}

func TestLoginWithWrongPassword(t *testing.T) {
	models.InitDatabase(true)

	// Create a new test user.
	createTestUser("test", "123456")

	ls := &LoginService{}

	user_id, status_code, status_msg, _ := ls.Login("test", "12345")

	assert.Equal(t, int64(-1), user_id)
	assert.Equal(t, int32(1), status_code)
	assert.Equal(t, "wrong password", status_msg)
}

func TestLoginWithCorrectPassword(t *testing.T) {
	models.InitDatabase(true)

	// Create a new test user.
	testUser, _ := createTestUser("test", "123456")

	ls := &LoginService{}

	user_id, status_code, status_msg, _ := ls.Login("test", "123456")

	assert.Equal(t, int64(1), user_id)
	assert.Equal(t, int32(0), status_code)
	assert.Equal(t, "login successfully", status_msg)
	assert.Equal(t, testUser.ID, user_id)
}
