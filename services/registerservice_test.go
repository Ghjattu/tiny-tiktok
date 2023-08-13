package services

import (
	"testing"

	"github.com/Ghjattu/tiny-tiktok/models"
	"github.com/stretchr/testify/assert"
)

var (
	registerService = &RegisterService{}
)

func TestRegisterWithEmptyUsername(t *testing.T) {
	user_id, status_code, status_msg, _ := registerService.Register("", "123456")

	assert.Equal(t, int64(-1), user_id)
	assert.Equal(t, int32(1), status_code)
	assert.Equal(t, "invalid username or password", status_msg)
}

func TestRegisterWithEmptyPassword(t *testing.T) {
	user_id, status_code, status_msg, _ := registerService.Register("test", "")

	assert.Equal(t, int64(-1), user_id)
	assert.Equal(t, int32(1), status_code)
	assert.Equal(t, "invalid username or password", status_msg)
}

func TestRegisterWithShortPassword(t *testing.T) {
	models.InitDatabase(true)

	user_id, status_code, status_msg, _ := registerService.Register("test", "123")

	assert.Equal(t, int64(-1), user_id)
	assert.Equal(t, int32(1), status_code)
	assert.Equal(t, "password is too short", status_msg)
}

func TestRegisterWithLongUsername(t *testing.T) {
	models.InitDatabase(true)

	user_id, status_code, status_msg, _ := registerService.Register(
		"1234567890123456789012345678901234567890123456789012345678901234567890", "123456")

	assert.Equal(t, int64(-1), user_id)
	assert.Equal(t, int32(1), status_code)
	assert.Equal(t, "username or password is too long", status_msg)
}

func TestRegisterWithLongPassword(t *testing.T) {
	models.InitDatabase(true)

	user_id, status_code, status_msg, _ := registerService.Register("test",
		"12345678901234567890123456789012345678901234567890123456789012345678901234567890")

	assert.Equal(t, int64(-1), user_id)
	assert.Equal(t, int32(1), status_code)
	assert.Equal(t, "username or password is too long", status_msg)
}

func TestRegisterWithExceed72BytesPassword(t *testing.T) {
	models.InitDatabase(true)

	user_id, status_code, status_msg, _ := registerService.Register("test",
		"密码密码密码密码密码密码密码密码密码密码密码密码密码密码密码")

	assert.Equal(t, int64(-1), user_id)
	assert.Equal(t, int32(1), status_code)
	assert.Equal(t, "password length exceeds 72 bytes", status_msg)
}

func TestRegisterWithRegisteredUsername(t *testing.T) {
	models.InitDatabase(true)

	// Create a new test user.
	models.CreateTestUser("test", "123456")

	user_id, status_code, status_msg, _ := registerService.Register("test", "123456")

	assert.Equal(t, int64(-1), user_id)
	assert.Equal(t, int32(1), status_code)
	assert.Equal(t, "the username has been registered", status_msg)
}

func TestRegisterWithValidUsernameAndPassword(t *testing.T) {
	models.InitDatabase(true)

	user_id, status_code, status_msg, _ := registerService.Register("test", "123456")

	assert.Equal(t, int64(1), user_id)
	assert.Equal(t, int32(0), status_code)
	assert.Equal(t, "register successfully", status_msg)
}
