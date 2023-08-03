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

func beforeRegisterTest(req *http.Request, isInitDatabase bool) (*httptest.ResponseRecorder, *RegisterResponse) {
	if isInitDatabase {
		models.InitDatabase(true)
	}

	r := gin.Default()
	r.POST("/api/douyin/user/register", Register)

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	rr := &RegisterResponse{}
	bytes, _ := io.ReadAll(w.Result().Body)
	json.Unmarshal(bytes, rr)

	return w, rr
}

func TestRegisterWithEmptyUsername(t *testing.T) {
	req := httptest.NewRequest("POST", "http://127.0.0.1/api/douyin/user/register?password=123456", nil)

	w, rr := beforeRegisterTest(req, true)

	assert.Equal(t, 200, w.Code)
	assert.Equal(t, int32(1), rr.StatusCode)
	assert.Equal(t, "invalid username or password", rr.StatusMsg)
}

func TestRegisterWithEmptyPassword(t *testing.T) {
	req := httptest.NewRequest("POST", "http://127.0.0.1/api/douyin/user/register?username=test", nil)

	w, rr := beforeRegisterTest(req, true)

	assert.Equal(t, 200, w.Code)
	assert.Equal(t, int32(1), rr.StatusCode)
	assert.Equal(t, "invalid username or password", rr.StatusMsg)
}

func TestRegisterWithShortPassword(t *testing.T) {
	req := httptest.NewRequest("POST", "http://127.0.0.1/api/douyin/user/register?username=test&password=123", nil)

	w, rr := beforeRegisterTest(req, true)

	assert.Equal(t, 200, w.Code)
	assert.Equal(t, int32(1), rr.StatusCode)
	assert.Equal(t, "password is too short", rr.StatusMsg)
}

func TestRegisterWithLongPassword(t *testing.T) {
	req := httptest.NewRequest("POST",
		"http://127.0.0.1/api/douyin/user/register?username=test&password=12345678901234567890123456789012345678901234567890123456789012345678901234567890", nil)

	w, rr := beforeRegisterTest(req, true)

	assert.Equal(t, 200, w.Code)
	assert.Equal(t, int32(1), rr.StatusCode)
	assert.Equal(t, "password length exceeds 72 bytes", rr.StatusMsg)
}

func TestRegisterWithRegisteredUsername(t *testing.T) {
	req := httptest.NewRequest("POST", "http://127.0.0.1/api/douyin/user/register?username=test&password=123456", nil)

	w, rr := beforeRegisterTest(req, true)

	assert.Equal(t, 200, w.Code)
	assert.Equal(t, int32(0), rr.StatusCode)
	assert.Equal(t, "register successfully", rr.StatusMsg)

	req = httptest.NewRequest("POST", "http://127.0.0.1/api/douyin/user/register?username=test&password=123456", nil)

	w, rr = beforeRegisterTest(req, false)

	assert.Equal(t, 200, w.Code)
	assert.Equal(t, int32(1), rr.StatusCode)
	assert.Equal(t, "the username has been registered", rr.StatusMsg)
}

func TestRegisterWithValidUsernameAndPassword(t *testing.T) {
	req := httptest.NewRequest("POST", "http://127.0.0.1/api/douyin/user/register?username=test&password=123456", nil)

	w, rr := beforeRegisterTest(req, true)

	assert.Equal(t, 200, w.Code)
	assert.Equal(t, int32(0), rr.StatusCode)
	assert.Equal(t, "register successfully", rr.StatusMsg)
}
