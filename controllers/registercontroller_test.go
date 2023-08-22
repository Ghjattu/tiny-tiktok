package controllers

import (
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRegister(t *testing.T) {
	setup()

	t.Run("register unsuccessfully", func(t *testing.T) {
		req := httptest.NewRequest("POST",
			"http://127.0.0.1/douyin/user/register/?password=123456", nil)

		w, r := sendRequest(req)
		res := r.(*RegisterResponse)

		assert.Equal(t, 200, w.Code)
		assert.Equal(t, int32(1), res.StatusCode)
	})
	t.Run("register successfully", func(t *testing.T) {
		req := httptest.NewRequest("POST",
			"http://127.0.0.1/douyin/user/register/?username=testTwo&password=123456", nil)

		w, r := sendRequest(req)
		res := r.(*RegisterResponse)

		assert.Equal(t, 200, w.Code)
		assert.Equal(t, int32(0), res.StatusCode)
	})
}
