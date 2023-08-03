package services

import (
	"testing"

	"github.com/Ghjattu/tiny-tiktok/models"
	"github.com/stretchr/testify/assert"
)

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

	rs := &RegisterService{}
	rs.Register("test", "123456")

	ls := &LoginService{}

	user_id, status_code, status_msg, _ := ls.Login("test", "12345")

	assert.Equal(t, int64(-1), user_id)
	assert.Equal(t, int32(1), status_code)
	assert.Equal(t, "wrong password", status_msg)
}

func TestLoginWithCorrectPassword(t *testing.T) {
	models.InitDatabase(true)

	rs := &RegisterService{}
	rs.Register("test", "123456")

	ls := &LoginService{}

	user_id, status_code, status_msg, _ := ls.Login("test", "123456")

	assert.Equal(t, int64(1), user_id)
	assert.Equal(t, int32(0), status_code)
	assert.Equal(t, "login successfully", status_msg)
}
