package controllers

import (
	"fmt"
	"net/http/httptest"
	"testing"

	"github.com/Ghjattu/tiny-tiktok/models"
	"github.com/stretchr/testify/assert"
)

func TestFollowActionWithInvalidUserID(t *testing.T) {
	models.InitDatabase(true)

	// Register a test user.
	_, _, token := registerTestUser("test", "123456")

	url := "http://127.0.0.1/douyin/relation/action/?to_user_id=abd&action_type=1&token=" + token
	req := httptest.NewRequest("POST", url, nil)

	w, r := sendRequest(req)
	res := r.(*Response)

	assert.Equal(t, 400, w.Code)
	assert.Equal(t, int32(1), res.StatusCode)
}

func TestFollowActionWithInvalidActionType(t *testing.T) {
	models.InitDatabase(true)

	// Register a test user.
	_, _, token := registerTestUser("test", "123456")

	url := "http://127.0.0.1/douyin/relation/action/?to_user_id=1&action_type=3&token=" + token
	req := httptest.NewRequest("POST", url, nil)

	w, r := sendRequest(req)
	res := r.(*Response)

	assert.Equal(t, 200, w.Code)
	assert.Equal(t, int32(1), res.StatusCode)
}

func TestFollowActionWithActionTypeOne(t *testing.T) {
	models.InitDatabase(true)

	// Register a test user.
	_, _, token := registerTestUser("test", "123456")
	// Create a test user.
	testUser, _ := models.CreateTestUser("test2", "123456")
	testUserIDStr := fmt.Sprintf("%d", testUser.ID)

	url := "http://127.0.0.1/douyin/relation/action/?to_user_id=" + testUserIDStr +
		"&action_type=1&token=" + token
	req := httptest.NewRequest("POST", url, nil)

	w, r := sendRequest(req)
	res := r.(*Response)

	assert.Equal(t, 200, w.Code)
	assert.Equal(t, int32(0), res.StatusCode)
}
