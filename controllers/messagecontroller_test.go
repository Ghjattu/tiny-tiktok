package controllers

import (
	"fmt"
	"net/http/httptest"
	"testing"

	"github.com/Ghjattu/tiny-tiktok/models"
	"github.com/Ghjattu/tiny-tiktok/redis"
	"github.com/stretchr/testify/assert"
)

func TestMessageAction(t *testing.T) {
	setup()

	t.Run("invalid action type", func(t *testing.T) {
		url := "http://127.0.0.1/douyin/message/action/?to_user_id=" + userIDStr +
			"&action_type=2&content=abc&token=" + token
		req := httptest.NewRequest("POST", url, nil)

		w, r := sendRequest(req)
		res := r.(*Response)

		assert.Equal(t, 200, w.Code)
		assert.Equal(t, int32(1), res.StatusCode)
		assert.Equal(t, "action type is invalid", res.StatusMsg)
	})

	t.Run("create new message successfully", func(t *testing.T) {
		// Create a test user.
		testUser, _ := models.CreateTestUser("test2", "123456")
		testUserIDStr := fmt.Sprintf("%d", testUser.ID)

		url := "http://127.0.0.1/douyin/message/action/?to_user_id=" + testUserIDStr +
			"&action_type=1&content=abc&token=" + token
		req := httptest.NewRequest("POST", url, nil)

		w, r := sendRequest(req)
		res := r.(*Response)

		assert.Equal(t, 200, w.Code)
		assert.Equal(t, int32(0), res.StatusCode)
	})
}

func TestMessageChat(t *testing.T) {
	setup()

	// Create a test user.
	testUser, _ := models.CreateTestUser("test2", "123456")
	testUserIDStr := fmt.Sprintf("%d", testUser.ID)
	// Create a test message.
	models.CreateTestMessage(userID, testUser.ID)

	t.Run("very large pre message time", func(t *testing.T) {
		redis.Rdb.FlushDB(redis.Ctx)

		url := "http://127.0.0.1/douyin/message/chat/?to_user_id=" + testUserIDStr +
			"&token=" + token + "&pre_msg_time=300000000000"
		req := httptest.NewRequest("GET", url, nil)

		w, r := sendRequest(req)
		res := r.(*MessageChatResponse)

		assert.Equal(t, 200, w.Code)
		assert.Equal(t, int32(0), res.StatusCode)
		assert.Equal(t, 1, len(res.MessageList))
	})

	t.Run("normal pre message time", func(t *testing.T) {
		url := "http://127.0.0.1/douyin/message/chat/?to_user_id=" + testUserIDStr +
			"&token=" + token
		req := httptest.NewRequest("GET", url, nil)

		w, r := sendRequest(req)
		res := r.(*MessageChatResponse)

		assert.Equal(t, 200, w.Code)
		assert.Equal(t, int32(0), res.StatusCode)
		assert.Equal(t, 1, len(res.MessageList))
	})
}
