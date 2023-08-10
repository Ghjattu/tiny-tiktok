package controllers

import (
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/Ghjattu/tiny-tiktok/models"
	"github.com/stretchr/testify/assert"
)

func TestGetUserByUserIDAndTokenWithInvalidUserID(t *testing.T) {
	models.InitDatabase(true)

	// Register a new test user.
	_, _, token := registerTestUser("test", "123456")

	req := httptest.NewRequest("GET",
		"http://127.0.0.1/douyin/user/?user_id=abc"+"&token="+token, nil)

	w, r := sendRequest(req)
	res := r.(*UserResponse)

	assert.Equal(t, 400, w.Code)
	assert.Equal(t, int32(1), res.StatusCode)
	assert.Equal(t, "invalid syntax", res.StatusMsg)
	assert.Equal(t, (*models.User)(nil), res.User)
}

func TestGetUserByUserIDAndTokenWithCorrectToken(t *testing.T) {
	models.InitDatabase(true)

	// Register a new test user.
	userID, testUser, token := registerTestUser("test", "123456")

	url := "http://127.0.0.1/douyin/user/?user_id=" + strconv.Itoa(int(userID)) +
		"&token=" + token
	req := httptest.NewRequest("GET", url, nil)

	w, r := sendRequest(req)
	res := r.(*UserResponse)

	assert.Equal(t, 200, w.Code)
	assert.Equal(t, int32(0), res.StatusCode)
	assert.Equal(t, "get user successfully", res.StatusMsg)
	assert.Equal(t, userID, res.User.ID)
	assert.Equal(t, testUser.Name, res.User.Name)
	assert.Equal(t, "", res.User.Password)
}
