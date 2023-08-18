package controllers

import (
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/Ghjattu/tiny-tiktok/models"
	"github.com/stretchr/testify/assert"
)

func TestGetUserByUserIDAndToken(t *testing.T) {
	models.Flush()

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
}
