package controllers

import (
	"fmt"
	"net/http/httptest"
	"testing"

	"github.com/Ghjattu/tiny-tiktok/models"
	"github.com/stretchr/testify/assert"
)

func TestFriendList(t *testing.T) {
	setup()

	// Create two test user.
	testUserOne, _ := models.CreateTestUser("testOne", "123456")
	testUserOneID := fmt.Sprintf("%d", testUserOne.ID)
	testUserTwo, _ := models.CreateTestUser("testTwo", "123456")
	// Create a test follow relationship.
	models.CreateTestFollowRel(testUserOne.ID, testUserTwo.ID)
	models.CreateTestFollowRel(testUserTwo.ID, testUserOne.ID)

	url := "http://127.0.0.1/douyin/relation/friend/list/?user_id=" + testUserOneID +
		"&token=" + token
	req := httptest.NewRequest("GET", url, nil)

	w, r := sendRequest(req)
	res := r.(*FriendListResponse)

	assert.Equal(t, 200, w.Code)
	assert.Equal(t, int32(0), res.StatusCode)
	assert.Equal(t, 1, len(res.UserList))
}
