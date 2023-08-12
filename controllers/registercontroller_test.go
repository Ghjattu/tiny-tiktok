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
