package services

import (
	"testing"

	"github.com/Ghjattu/tiny-tiktok/models"
	"github.com/stretchr/testify/assert"
)

func TestRegisterWithEmptyUsername(t *testing.T) {
	rs := &RegisterService{}

	user_id, status_code, status_msg, _ := rs.Register("", "123456")

	assert.Equal(t, int64(-1), user_id)
	assert.Equal(t, int32(1), status_code)
	assert.Equal(t, "invalid username or password", status_msg)
}

func TestRegisterWithEmptyPassword(t *testing.T) {
	rs := &RegisterService{}

	user_id, status_code, status_msg, _ := rs.Register("test", "")

	assert.Equal(t, int64(-1), user_id)
	assert.Equal(t, int32(1), status_code)
	assert.Equal(t, "invalid username or password", status_msg)
}

func TestRegisterWithRegisteredUsername(t *testing.T) {
	models.InitDatabase(true)

	rs := &RegisterService{}

	rs.Register("test", "123456")
	user_id, status_code, status_msg, _ := rs.Register("test", "123456")

	assert.Equal(t, int64(-1), user_id)
	assert.Equal(t, int32(1), status_code)
	assert.Equal(t, "the username has been registered", status_msg)
}

func TestRegisterWithValidUsernameAndPassword(t *testing.T) {
	models.InitDatabase(true)

	rs := &RegisterService{}

	user_id, status_code, status_msg, _ := rs.Register("test", "123456")

	assert.Equal(t, int64(1), user_id)
	assert.Equal(t, int32(0), status_code)
	assert.Equal(t, "register successfully", status_msg)
}
