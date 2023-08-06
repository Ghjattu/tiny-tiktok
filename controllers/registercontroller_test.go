package controllers

import (
	"net/http/httptest"
	"testing"

	"github.com/Ghjattu/tiny-tiktok/models"
	"github.com/stretchr/testify/assert"
)

func TestRegisterWithEmptyUsername(t *testing.T) {
	models.InitDatabase(true)

	req := httptest.NewRequest("POST",
		"http://127.0.0.1/douyin/user/register/?password=123456", nil)

	w, r := sendRequest(req)
	res := r.(*RegisterResponse)

	assert.Equal(t, 200, w.Code)
	assert.Equal(t, int32(1), res.StatusCode)
	assert.Equal(t, "invalid username or password", res.StatusMsg)
}

func TestRegisterWithEmptyPassword(t *testing.T) {
	models.InitDatabase(true)

	req := httptest.NewRequest("POST",
		"http://127.0.0.1/douyin/user/register/?username=test", nil)

	w, r := sendRequest(req)
	res := r.(*RegisterResponse)

	assert.Equal(t, 200, w.Code)
	assert.Equal(t, int32(1), res.StatusCode)
	assert.Equal(t, "invalid username or password", res.StatusMsg)
}

func TestRegisterWithShortPassword(t *testing.T) {
	models.InitDatabase(true)

	req := httptest.NewRequest("POST",
		"http://127.0.0.1/douyin/user/register/?username=test&password=123", nil)

	w, r := sendRequest(req)
	res := r.(*RegisterResponse)

	assert.Equal(t, 200, w.Code)
	assert.Equal(t, int32(1), res.StatusCode)
	assert.Equal(t, "password is too short", res.StatusMsg)
}

func TestRegisterWithLongPassword(t *testing.T) {
	models.InitDatabase(true)

	req := httptest.NewRequest("POST",
		"http://127.0.0.1/douyin/user/register/?username=test&password=12345678901234567890123456789012345678901234567890123456789012345678901234567890", nil)

	w, r := sendRequest(req)
	res := r.(*RegisterResponse)

	assert.Equal(t, 200, w.Code)
	assert.Equal(t, int32(1), res.StatusCode)
	assert.Equal(t, "username or password is too long", res.StatusMsg)
}

func TestRegisterWithRegisteredUsername(t *testing.T) {
	models.InitDatabase(true)

	// Register a new test user.
	registerTestUser("test", "123456")

	req := httptest.NewRequest("POST",
		"http://127.0.0.1/douyin/user/register/?username=test&password=123456", nil)

	w, r := sendRequest(req)
	res := r.(*RegisterResponse)

	assert.Equal(t, 200, w.Code)
	assert.Equal(t, int32(1), res.StatusCode)
	assert.Equal(t, "the username has been registered", res.StatusMsg)
}

func TestRegisterWithValidUsernameAndPassword(t *testing.T) {
	models.InitDatabase(true)

	req := httptest.NewRequest("POST",
		"http://127.0.0.1/douyin/user/register/?username=test&password=123456", nil)

	w, r := sendRequest(req)
	res := r.(*RegisterResponse)

	assert.Equal(t, 200, w.Code)
	assert.Equal(t, int32(0), res.StatusCode)
	assert.Equal(t, "register successfully", res.StatusMsg)
}
