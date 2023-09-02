package controllers

import (
	"fmt"
	"net/http/httptest"
	"testing"

	"github.com/Ghjattu/tiny-tiktok/models"
	"github.com/stretchr/testify/assert"
)

func TestFollowAction(t *testing.T) {
	setup()

	// Create a test user.
	testUser, _ := models.CreateTestUser("test2", "123456")
	testUserIDStr := fmt.Sprintf("%d", testUser.ID)

	t.Run("invalid action type", func(t *testing.T) {
		url := "http://127.0.0.1/douyin/relation/action/?to_user_id=1&action_type=3&token=" + token
		req := httptest.NewRequest("POST", url, nil)

		w, r := sendRequest(req)
		res := r.(*Response)

		assert.Equal(t, 200, w.Code)
		assert.Equal(t, int32(1), res.StatusCode)
	})

	t.Run("create new follow relationship successfully", func(t *testing.T) {
		url := "http://127.0.0.1/douyin/relation/action/?to_user_id=" + testUserIDStr +
			"&action_type=1&token=" + token
		req := httptest.NewRequest("POST", url, nil)

		w, r := sendRequest(req)
		res := r.(*Response)

		assert.Equal(t, 200, w.Code)
		assert.Equal(t, int32(0), res.StatusCode)
	})

	t.Run("delete follow relationship successfully", func(t *testing.T) {
		// Create a test follow relationship.
		models.CreateTestFollowRel(userID, testUser.ID)

		url := "http://127.0.0.1/douyin/relation/action/?to_user_id=" + testUserIDStr +
			"&action_type=2&token=" + token
		req := httptest.NewRequest("POST", url, nil)

		w, r := sendRequest(req)
		res := r.(*Response)

		assert.Equal(t, 200, w.Code)
		assert.Equal(t, int32(0), res.StatusCode)
	})
}

func TestFollowingList(t *testing.T) {
	setup()

	// Create a test user.
	testUser, _ := models.CreateTestUser("test2", "123456")
	// Create a test follow relationship.
	models.CreateTestFollowRel(userID, testUser.ID)

	url := "http://127.0.0.1/douyin/relation/follow/list/?user_id=" + userIDStr +
		"&token=" + token
	req := httptest.NewRequest("GET", url, nil)

	w, r := sendRequest(req)
	res := r.(*UserListResponse)

	assert.Equal(t, 200, w.Code)
	assert.Equal(t, int32(0), res.StatusCode)
	assert.Equal(t, 1, len(res.UserList))
}

func TestFollowerList(t *testing.T) {
	setup()

	// Create a test user.
	testUser, _ := models.CreateTestUser("test2", "123456")
	// Create a test follow relationship.
	models.CreateTestFollowRel(testUser.ID, userID)

	url := "http://127.0.0.1/douyin/relation/follower/list/?user_id=" + userIDStr +
		"&token=" + token
	req := httptest.NewRequest("GET", url, nil)

	w, r := sendRequest(req)
	res := r.(*UserListResponse)

	assert.Equal(t, 200, w.Code)
	assert.Equal(t, int32(0), res.StatusCode)
	assert.Equal(t, 1, len(res.UserList))
}
