package controllers

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/Ghjattu/tiny-tiktok/middleware/jwt"
	"github.com/Ghjattu/tiny-tiktok/models"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func beforeUserTest(req *http.Request, isInitDatabase bool, needAuth bool) (*httptest.ResponseRecorder, *RegisterResponse, *UserResponse) {
	if isInitDatabase {
		models.InitDatabase(true)
	}

	r := gin.Default()
	r.POST("/douyin/user/register/", Register)
	r.POST("/douyin/user/login/", Login)
	if needAuth {
		r.GET("/douyin/user/", jwt.AuthorizationGet(), GetUserByUserIDAndToken)
	} else {
		r.GET("/douyin/user/", GetUserByUserIDAndToken)
	}

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	// If the request method is POST, then the response is RegisterResponse.
	if req.Method == "POST" {
		rr := &RegisterResponse{}
		bytes, _ := io.ReadAll(w.Result().Body)
		json.Unmarshal(bytes, rr)

		return w, rr, nil
	}

	// Otherwise the request method is GET, then the response is UserResponse.
	ur := &UserResponse{}
	bytes, _ := io.ReadAll(w.Result().Body)
	json.Unmarshal(bytes, ur)

	return w, nil, ur
}

func TestGetUserByUserIDAndTokenWithInvalidUserID(t *testing.T) {
	req := httptest.NewRequest("GET", "http://127.0.0.1/douyin/user/?user_id=abc", nil)

	w, _, ur := beforeUserTest(req, false, false)

	assert.Equal(t, 400, w.Code)
	assert.Equal(t, int32(1), ur.StatusCode)
	assert.Equal(t, "invalid syntax", ur.StatusMsg)
	assert.Equal(t, (*models.User)(nil), ur.User)
}

func TestGetUserByUserIDAndTokenWithNotExistUserID(t *testing.T) {
	req := httptest.NewRequest("POST",
		"http://127.0.0.1/douyin/user/register/?username=test&password=123456", nil)

	w, rr, _ := beforeUserTest(req, true, false)

	assert.Equal(t, 200, w.Code)
	assert.Equal(t, int32(0), rr.StatusCode)
	assert.Equal(t, "register successfully", rr.StatusMsg)

	userID := rr.UserID
	token := rr.Token

	url := "http://127.0.0.1/douyin/user/?user_id=" + strconv.Itoa(int(userID)+1) +
		"&token=" + token
	req = httptest.NewRequest("GET", url, nil)

	w, _, ur := beforeUserTest(req, false, false)

	assert.Equal(t, 200, w.Code)
	assert.Equal(t, int32(1), ur.StatusCode)
	assert.Equal(t, "user not found", ur.StatusMsg)
}

func TestGetUserByUserIDAndTokenWithEmptyToken(t *testing.T) {
	req := httptest.NewRequest("POST",
		"http://127.0.0.1/douyin/user/register/?username=test&password=123456", nil)

	w, rr, _ := beforeUserTest(req, true, false)

	assert.Equal(t, 200, w.Code)
	assert.Equal(t, int32(0), rr.StatusCode)
	assert.Equal(t, "register successfully", rr.StatusMsg)

	userID := rr.UserID

	url := "http://127.0.0.1/douyin/user/?user_id=" + strconv.Itoa(int(userID)) + "&token="
	req = httptest.NewRequest("GET", url, nil)

	w, _, ur := beforeUserTest(req, false, true)

	assert.Equal(t, 401, w.Code)
	assert.Equal(t, int32(1), ur.StatusCode)
	assert.Equal(t, "invalid token", ur.StatusMsg)
}

func TestGetUserByUserIDAndTokenWithInvalidToken(t *testing.T) {
	req := httptest.NewRequest("POST",
		"http://127.0.0.1/douyin/user/register/?username=test&password=123456", nil)

	w, rr, _ := beforeUserTest(req, true, false)

	assert.Equal(t, 200, w.Code)
	assert.Equal(t, int32(0), rr.StatusCode)
	assert.Equal(t, "register successfully", rr.StatusMsg)

	userID := rr.UserID
	invalidToken := rr.Token + "1"

	url := "http://127.0.0.1/douyin/user/?user_id=" + strconv.Itoa(int(userID)) +
		"&token=" + invalidToken
	req = httptest.NewRequest("GET", url, nil)

	w, _, ur := beforeUserTest(req, false, true)

	assert.Equal(t, 401, w.Code)
	assert.Equal(t, int32(1), ur.StatusCode)
	assert.Equal(t, "invalid token", ur.StatusMsg)
	assert.Equal(t, (*models.User)(nil), ur.User)
}

func TestGetUserByUserIDAndTokenWithCorrectToken(t *testing.T) {
	req := httptest.NewRequest("POST",
		"http://127.0.0.1/douyin/user/register/?username=test&password=123456", nil)

	w, rr, _ := beforeUserTest(req, true, false)

	assert.Equal(t, 200, w.Code)
	assert.Equal(t, int32(0), rr.StatusCode)
	assert.Equal(t, "register successfully", rr.StatusMsg)

	userID := rr.UserID
	token := rr.Token

	url := "http://127.0.0.1/douyin/user/?user_id=" + strconv.Itoa(int(userID)) +
		"&token=" + token
	req = httptest.NewRequest("GET", url, nil)

	w, _, ur := beforeUserTest(req, false, true)

	assert.Equal(t, 200, w.Code)
	assert.Equal(t, int32(0), ur.StatusCode)
	assert.Equal(t, "get user successfully", ur.StatusMsg)
	assert.Equal(t, userID, ur.User.ID)
	assert.Equal(t, "test", ur.User.Name)
	assert.Equal(t, "", ur.User.Password)
}
