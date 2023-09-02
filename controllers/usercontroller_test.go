package controllers

import (
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetUserByUserIDAndToken(t *testing.T) {
	setup()

	url := "http://127.0.0.1/douyin/user/?user_id=" + strconv.Itoa(int(userID)) +
		"&token=" + token
	req := httptest.NewRequest("GET", url, nil)

	w, r := sendRequest(req)
	res := r.(*UserResponse)

	assert.Equal(t, 200, w.Code)
	assert.Equal(t, int32(0), res.StatusCode)
}
