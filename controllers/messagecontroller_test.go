package controllers

import (
	"fmt"
	"net/http/httptest"
	"testing"

	"github.com/Ghjattu/tiny-tiktok/models"
	"github.com/stretchr/testify/assert"
)

func TestMessageActionWithInvalidActionType(t *testing.T) {
	models.InitDatabase(true)

	// Register a test user.
	userID, _, token := registerTestUser("test", "123456")
	userIDStr := fmt.Sprintf("%d", userID)

	url := "http://127.0.0.1/douyin/message/action/?to_user_id=" + userIDStr +
		"&action_type=2&content=abc&token=" + token
	req := httptest.NewRequest("POST", url, nil)

	w, r := sendRequest(req)
	res := r.(*Response)

	assert.Equal(t, 200, w.Code)
	assert.Equal(t, int32(1), res.StatusCode)
	assert.Equal(t, "action type is invalid", res.StatusMsg)
}

func TestMessageAction(t *testing.T) {
	models.InitDatabase(true)

	// Register a test user.
	userID, _, token := registerTestUser("test", "123456")
	userIDStr := fmt.Sprintf("%d", userID)

	url := "http://127.0.0.1/douyin/message/action/?to_user_id=" + userIDStr +
		"&action_type=1&content=abc&token=" + token
	req := httptest.NewRequest("POST", url, nil)

	w, r := sendRequest(req)
	res := r.(*Response)

	assert.Equal(t, 200, w.Code)
	assert.Equal(t, int32(0), res.StatusCode)
	assert.Equal(t, "create new message successfully", res.StatusMsg)
}
