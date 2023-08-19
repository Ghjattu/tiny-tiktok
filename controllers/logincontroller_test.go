package controllers

import (
	"net/http/httptest"
	"testing"

	"github.com/Ghjattu/tiny-tiktok/models"
	"github.com/stretchr/testify/assert"
)

func TestLoginWithWrongPassword(t *testing.T) {
	models.Flush()

	// Register a new test user.
	registerTestUser("test", "123456")

	req := httptest.NewRequest("POST",
		"http://127.0.0.1/douyin/user/login/?username=test&password=12345", nil)

	w, r := sendRequest(req)
	res := r.(*LoginResponse)

	assert.Equal(t, 200, w.Code)
	assert.Equal(t, int32(1), res.StatusCode)
	assert.Equal(t, "wrong password", res.StatusMsg)
}

func TestLoginWithCorrectPassword(t *testing.T) {
	models.Flush()

	// Register a new test user.
	userID, _, _ := registerTestUser("test", "123456")

	req := httptest.NewRequest("POST",
		"http://127.0.0.1/douyin/user/login/?username=test&password=123456", nil)

	w, r := sendRequest(req)
	res := r.(*LoginResponse)

	assert.Equal(t, 200, w.Code)
	assert.Equal(t, int32(0), res.StatusCode)
	assert.Equal(t, "login successfully", res.StatusMsg)
	assert.Equal(t, userID, res.UserID)
}
