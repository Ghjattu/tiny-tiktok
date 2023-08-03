package controllers

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Ghjattu/tiny-tiktok/models"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func beforeLoginTest(req *http.Request, isInitDatabase bool) (*httptest.ResponseRecorder, *LoginResponse) {
	if isInitDatabase {
		models.InitDatabase(true)
	}

	r := gin.Default()
	r.POST("/api/douyin/user/register", Register)
	r.POST("/api/douyin/user/login", Login)

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	lr := &LoginResponse{}
	bytes, _ := io.ReadAll(w.Result().Body)
	json.Unmarshal(bytes, lr)

	return w, lr
}

func TestLoginWithNotExistName(t *testing.T) {
	req := httptest.NewRequest("POST", "http://127.0.0.1/api/douyin/user/login?username=test&password=123456", nil)

	w, lr := beforeLoginTest(req, true)

	assert.Equal(t, 200, w.Code)
	assert.Equal(t, int32(1), lr.StatusCode)
	assert.Equal(t, "username not found", lr.StatusMsg)
}

func TestLoginWithWrongPassword(t *testing.T) {
	req := httptest.NewRequest("POST", "http://127.0.0.1/api/douyin/user/register?username=test&password=123456", nil)

	w, lr := beforeLoginTest(req, true)

	assert.Equal(t, 200, w.Code)
	assert.Equal(t, int32(0), lr.StatusCode)
	assert.Equal(t, "register successfully", lr.StatusMsg)

	req = httptest.NewRequest("POST", "http://127.0.0.1/api/douyin/user/login?username=test&password=12345", nil)

	w, lr = beforeLoginTest(req, false)

	assert.Equal(t, 200, w.Code)
	assert.Equal(t, int32(1), lr.StatusCode)
	assert.Equal(t, "wrong password", lr.StatusMsg)
}

func TestLoginWithCorrectPassword(t *testing.T) {
	req := httptest.NewRequest("POST", "http://127.0.0.1/api/douyin/user/register?username=test&password=123456", nil)

	w, lr := beforeLoginTest(req, true)

	assert.Equal(t, 200, w.Code)
	assert.Equal(t, int32(0), lr.StatusCode)
	assert.Equal(t, "register successfully", lr.StatusMsg)

	req = httptest.NewRequest("POST", "http://127.0.0.1/api/douyin/user/login?username=test&password=123456", nil)

	w, lr = beforeLoginTest(req, false)

	assert.Equal(t, 200, w.Code)
	assert.Equal(t, int32(0), lr.StatusCode)
	assert.Equal(t, "login successfully", lr.StatusMsg)
}
