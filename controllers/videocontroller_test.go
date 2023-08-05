package controllers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"strconv"
	"strings"
	"testing"

	"github.com/Ghjattu/tiny-tiktok/middleware/jwt"
	"github.com/Ghjattu/tiny-tiktok/models"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
)

var (
	serverIP   = ""
	serverPort = ""
)

func checkResponseType(req *http.Request) interface{} {
	pathSlice := strings.Split(req.URL.Path, "/")
	lastPath := pathSlice[len(pathSlice)-2]

	if lastPath == "register" {
		return &RegisterResponse{}
	} else if lastPath == "feed" {
		return &FeedResponse{}
	} else if lastPath == "list" {
		return &PublishListResponse{}
	} else if lastPath == "action" {
		return &Response{}
	} else {
		return nil
	}
}

func beforeVideoTest(req *http.Request, needInitDatabase bool, needAuth bool) (*httptest.ResponseRecorder, interface{}) {
	// Load environment variables.
	godotenv.Load("../.env")
	serverIP = os.Getenv("SERVER_IP")
	serverPort = os.Getenv("SERVER_PORT")

	if needInitDatabase {
		models.InitDatabase(true)
	}

	r := gin.Default()
	r.POST("/douyin/user/register/", Register)
	r.GET("/douyin/feed/", jwt.AuthorizationFeed(), Feed)
	if needAuth {
		r.GET("/douyin/publish/list/", jwt.AuthorizationGet(), GetPublishListByAuthorID)
		r.POST("/douyin/publish/action/", jwt.AuthorizationPost(), PublishNewVideo)
	} else {
		r.GET("/douyin/publish/list/", GetPublishListByAuthorID)
		r.POST("/douyin/publish/action/", PublishNewVideo)
	}

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	res := checkResponseType(req)
	bytes, _ := io.ReadAll(w.Result().Body)
	json.Unmarshal(bytes, res)

	return w, res
}

// testVideoAccess tests whether the video can be accessed through the URL.
func testVideoAccess(req *http.Request) *httptest.ResponseRecorder {
	r := gin.Default()
	r.Static("/static/videos", "../public/")

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	return w
}

// constructTestForm constructs a test form with a test file and form fields.
//
//	@param formFields map[string]string
//	@return *bytes.Buffer
//	@return *multipart.Writer
//	@return error
func constructTestForm(formFields map[string]string) (*bytes.Buffer, *multipart.Writer, error) {
	// Read the test video.
	file, err := os.Open("../data/bear.mp4")
	if err != nil {
		return nil, nil, err
	}
	defer file.Close()

	// Construct form data.
	form := bytes.NewBuffer([]byte(""))
	writer := multipart.NewWriter(form)
	defer writer.Close()

	// Add form fields.
	for key, value := range formFields {
		writer.WriteField(key, value)
	}

	// Add form file.
	part, err := writer.CreateFormFile("data", "bear.mp4")
	if err != nil {
		return nil, nil, err
	}

	// Copy file data to form file.
	_, err = io.Copy(part, file)
	if err != nil {
		return nil, nil, err
	}

	return form, writer, nil
}

// registerTestUser registers a new test user and returns the status_code, user_id and token.
//
//	@return int32 "status_code"
//	@return int64 "user_id"
//	@return string "token"
func registerTestUser() (int32, int64, string) {
	// Register a new user.
	req := httptest.NewRequest("POST",
		"http://127.0.0.1/douyin/user/register/?username=test&password=123456", nil)

	_, r := beforeVideoTest(req, true, false)
	res := r.(*RegisterResponse)

	return res.StatusCode, res.UserID, res.Token
}

func TestPublishNewVideoWithInvalidToken(t *testing.T) {
	// Register a new test user and get the token.
	_, _, token := registerTestUser()

	invalidToken := token + "1"

	// Construct a test form.
	formFields := map[string]string{
		"title": "Test Title",
		"token": invalidToken,
	}
	form, writer, err := constructTestForm(formFields)
	if err != nil {
		t.Fatalf("failed to construct form data: %v", err)
	}

	req := httptest.NewRequest("POST",
		"http://127.0.0.1/douyin/publish/action/", form)
	req.Header.Set("Content-Type", writer.FormDataContentType())

	w, r := beforeVideoTest(req, false, true)
	res := r.(*Response)

	assert.Equal(t, 401, w.Code)
	assert.Equal(t, int32(1), res.StatusCode)
	assert.Equal(t, "invalid token", res.StatusMsg)
}

func TestPublishNewVideoWithValidToken(t *testing.T) {
	// Register a new test user and get the id and token.
	_, userID, token := registerTestUser()

	// Construct a test form.
	formFields := map[string]string{
		"title": "Test Title",
		"token": token,
	}
	form, writer, err := constructTestForm(formFields)
	if err != nil {
		t.Fatalf("failed to construct form data: %v", err)
	}

	// Publish a new video.
	req := httptest.NewRequest("POST",
		"http://127.0.0.1/douyin/publish/action/", form)
	req.Header.Set("Content-Type", writer.FormDataContentType())

	w, r := beforeVideoTest(req, false, true)
	res := r.(*Response)

	assert.Equal(t, 200, w.Code)
	assert.Equal(t, int32(0), res.StatusCode)
	assert.Equal(t, "create new video successfully", res.StatusMsg)

	// Test the video access.
	videoURL := fmt.Sprintf("http://%s:%s/static/videos/%s_bear.mp4", serverIP, serverPort, strconv.Itoa(int(userID)))
	req = httptest.NewRequest("GET", videoURL, nil)

	w = testVideoAccess(req)

	assert.Equal(t, 200, w.Code)
	assert.Equal(t, "video/mp4", w.Header().Get("Content-Type"))
}

func TestGetPublishListByAuthorIDWithEmptyID(t *testing.T) {
	req := httptest.NewRequest("GET",
		"http://127.0.0.1/douyin/publish/list/", nil)

	w, r := beforeVideoTest(req, true, false)
	res := r.(*PublishListResponse)

	assert.Equal(t, 400, w.Code)
	assert.Equal(t, int32(1), res.StatusCode)
	assert.Equal(t, "invalid syntax", res.StatusMsg)
	assert.Equal(t, 0, len(res.VideoList))
}

func TestGetPublishListByAuthorIDWithInvalidID(t *testing.T) {
	req := httptest.NewRequest("GET",
		"http://127.0.0.1/douyin/publish/list/?user_id=abc", nil)

	w, r := beforeVideoTest(req, true, false)
	res := r.(*PublishListResponse)

	assert.Equal(t, 400, w.Code)
	assert.Equal(t, int32(1), res.StatusCode)
	assert.Equal(t, "invalid syntax", res.StatusMsg)
	assert.Equal(t, 0, len(res.VideoList))
}

func TestGetPublishListByAuthorIDWithOutOfRangeID(t *testing.T) {
	req := httptest.NewRequest("GET",
		"http://127.0.0.1/douyin/publish/list/?user_id=99999999999999999999999999", nil)

	w, r := beforeVideoTest(req, true, false)
	res := r.(*PublishListResponse)

	assert.Equal(t, 400, w.Code)
	assert.Equal(t, int32(1), res.StatusCode)
	assert.Equal(t, "user id out of range", res.StatusMsg)
	assert.Equal(t, 0, len(res.VideoList))
}

func TestGetPublishListByAuthorIDWithValidID(t *testing.T) {
	// Register a new test user and get the id and token.
	_, userID, token := registerTestUser()

	// Publish a new video.
	formFields := map[string]string{
		"title": "Test Title",
		"token": token,
	}
	form, writer, err := constructTestForm(formFields)
	if err != nil {
		t.Fatalf("failed to construct form data: %v", err)
	}

	req := httptest.NewRequest("POST",
		"http://127.0.0.1/douyin/publish/action/", form)
	req.Header.Set("Content-Type", writer.FormDataContentType())

	w, r := beforeVideoTest(req, false, true)
	res := r.(*Response)

	assert.Equal(t, 200, w.Code)
	assert.Equal(t, int32(0), res.StatusCode)
	assert.Equal(t, "create new video successfully", res.StatusMsg)

	// Get publish list by author id.
	url := "http://127.0.0.1/douyin/publish/list/?user_id=" + strconv.Itoa(int(userID)) +
		"&token=" + token
	req = httptest.NewRequest("GET", url, nil)

	w, r = beforeVideoTest(req, false, true)
	res2 := r.(*PublishListResponse)

	assert.Equal(t, 200, w.Code)
	assert.Equal(t, int32(0), res2.StatusCode)
	assert.Equal(t, "get publish list successfully", res2.StatusMsg)
	assert.Equal(t, 1, len(res2.VideoList))
}

func TestFeedWithEmptyLatestTime(t *testing.T) {
	req := httptest.NewRequest("GET", "http://127.0.0.1/douyin/feed/", nil)

	w, r := beforeVideoTest(req, true, true)
	res := r.(*FeedResponse)

	assert.Equal(t, 200, w.Code)
	assert.Equal(t, int32(0), res.StatusCode)
	assert.Equal(t, "get most 30 videos successfully", res.StatusMsg)
	assert.Equal(t, 0, len(res.VideoList))
}
