package controllers

import (
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLogin(t *testing.T) {
	setup()

	t.Run("wrong password", func(t *testing.T) {
		req := httptest.NewRequest("POST",
			"http://127.0.0.1/douyin/user/login/?username=test&password=12345", nil)

		w, r := sendRequest(req)
		res := r.(*LoginResponse)

		assert.Equal(t, 200, w.Code)
		assert.Equal(t, int32(1), res.StatusCode)
	})

	t.Run("correct password", func(t *testing.T) {
		req := httptest.NewRequest("POST",
			"http://127.0.0.1/douyin/user/login/?username=test&password=123456", nil)

		w, r := sendRequest(req)
		res := r.(*LoginResponse)

		assert.Equal(t, 200, w.Code)
		assert.Equal(t, int32(0), res.StatusCode)
	})
}
