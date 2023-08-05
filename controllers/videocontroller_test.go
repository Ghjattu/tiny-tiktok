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

func beforeVideoTest(req *http.Request, needInitDatabase bool, needAuth bool) (*httptest.ResponseRecorder, *RegisterResponse, *VideoResponse) {
	// Load environment variables.
	godotenv.Load("../.env")
	serverIP = os.Getenv("SERVER_IP")
	serverPort = os.Getenv("SERVER_PORT")

	if needInitDatabase {
		models.InitDatabase(true)
	}

	r := gin.Default()
	r.POST("/douyin/user/register/", Register)
	if needAuth {
		r.GET("/douyin/publish/list/", jwt.AuthorizationGet(), GetPublishListByAuthorID)
		r.POST("/douyin/publish/action/", jwt.AuthorizationPost(), PublishNewVideo)
	} else {
		r.GET("/douyin/publish/list/", GetPublishListByAuthorID)
		r.POST("/douyin/publish/action/", PublishNewVideo)
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

	// Otherwise the request method is GET, then the response is VideoResponse.
	vr := &VideoResponse{}
	bytes, _ := io.ReadAll(w.Result().Body)
	json.Unmarshal(bytes, vr)

	return w, nil, vr
}

// testVideoAccess tests the video access.
func testVideoAccess(req *http.Request) *httptest.ResponseRecorder {
	r := gin.Default()
	r.Static("/static/videos", "../public/")

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	return w
}

// createTestFile creates a temporary testing file with the given filename and content.
//
//	@param filename string
//	@param content string
//	@return *os.File
//	@return error
// func createTestFile(filename, content string) (*os.File, error) {
// 	file, err := os.CreateTemp("", filename)
// 	if err != nil {
// 		return nil, err
// 	}

// 	if _, err := file.WriteString(content); err != nil {
// 		file.Close()
// 		return nil, err
// 	}

// 	return file, nil
// }

// constructForm constructs a form data with a test file and form fields.
//
//	@param formFields map[string]string
//	@return *bytes.Buffer
//	@return *multipart.Writer
//	@return error
func constructForm(formFields map[string]string) (*bytes.Buffer, *multipart.Writer, error) {
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

// registerNewUser registers a new user and returns the status_code, user_id and token.
//
//	@return int32 "status_code"
//	@return int64 "user_id"
//	@return string "token"
func registerNewUser() (int32, int64, string) {
	// Register a new user.
	req := httptest.NewRequest("POST",
		"http://127.0.0.1/douyin/user/register/?username=test&password=123456", nil)

	_, rr, _ := beforeVideoTest(req, true, false)

	return rr.StatusCode, rr.UserID, rr.Token
}

func TestPublishNewVideoWithInvalidToken(t *testing.T) {
	// Register a new user and get the token.
	statusCode, _, token := registerNewUser()

	assert.Equal(t, int32(0), statusCode)

	invalidToken := token + "1"

	// Construct a test form.
	formFields := map[string]string{
		"title": "Test Title",
		"token": invalidToken,
	}
	form, writer, err := constructForm(formFields)
	if err != nil {
		t.Fatalf("failed to construct form data: %v", err)
	}

	req := httptest.NewRequest("POST",
		"http://127.0.0.1/douyin/publish/action/", form)
	req.Header.Set("Content-Type", writer.FormDataContentType())

	w, response, _ := beforeVideoTest(req, false, true)

	assert.Equal(t, 401, w.Code)
	assert.Equal(t, int32(1), response.StatusCode)
	assert.Equal(t, "invalid token", response.StatusMsg)
}

func TestPublishNewVideoWithValidToken(t *testing.T) {
	// Register a new user.
	status_code, userID, token := registerNewUser()

	assert.Equal(t, int32(0), status_code)

	// Construct a test form.
	formFields := map[string]string{
		"title": "Test Title",
		"token": token,
	}
	form, writer, err := constructForm(formFields)
	if err != nil {
		t.Fatalf("failed to construct form data: %v", err)
	}

	// Publish a new video.
	req := httptest.NewRequest("POST",
		"http://127.0.0.1/douyin/publish/action/", form)
	req.Header.Set("Content-Type", writer.FormDataContentType())

	w, response, _ := beforeVideoTest(req, false, true)

	assert.Equal(t, 200, w.Code)
	assert.Equal(t, int32(0), response.StatusCode)
	assert.Equal(t, "create new video successfully", response.StatusMsg)

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

	w, _, vr := beforeVideoTest(req, true, false)

	assert.Equal(t, 400, w.Code)
	assert.Equal(t, int32(1), vr.StatusCode)
	assert.Equal(t, "invalid syntax", vr.StatusMsg)
	assert.Equal(t, 0, len(vr.VideoList))
}

func TestGetPublishListByAuthorIDWithInvalidID(t *testing.T) {
	req := httptest.NewRequest("GET",
		"http://127.0.0.1/douyin/publish/list/?user_id=abc", nil)

	w, _, vr := beforeVideoTest(req, true, false)

	assert.Equal(t, 400, w.Code)
	assert.Equal(t, int32(1), vr.StatusCode)
	assert.Equal(t, "invalid syntax", vr.StatusMsg)
	assert.Equal(t, 0, len(vr.VideoList))
}

func TestGetPublishListByAuthorIDWithOutOfRangeID(t *testing.T) {
	req := httptest.NewRequest("GET",
		"http://127.0.0.1/douyin/publish/list/?user_id=99999999999999999999999999", nil)

	w, _, vr := beforeVideoTest(req, true, false)

	assert.Equal(t, 400, w.Code)
	assert.Equal(t, int32(1), vr.StatusCode)
	assert.Equal(t, "user id out of range", vr.StatusMsg)
	assert.Equal(t, 0, len(vr.VideoList))
}

func TestGetPublishListByAuthorIDWithValidID(t *testing.T) {
	// Register a new user.
	status_code, userID, token := registerNewUser()

	assert.Equal(t, int32(0), status_code)

	// Publish a new video.
	formFields := map[string]string{
		"title": "Test Title",
		"token": token,
	}
	form, writer, err := constructForm(formFields)
	if err != nil {
		t.Fatalf("failed to construct form data: %v", err)
	}

	req := httptest.NewRequest("POST",
		"http://127.0.0.1/douyin/publish/action/", form)
	req.Header.Set("Content-Type", writer.FormDataContentType())

	w, rr, _ := beforeVideoTest(req, false, true)

	assert.Equal(t, 200, w.Code)
	assert.Equal(t, int32(0), rr.StatusCode)
	assert.Equal(t, "create new video successfully", rr.StatusMsg)

	// Get publish list by author id.
	url := "http://127.0.0.1/douyin/publish/list/?user_id=" + strconv.Itoa(int(userID)) +
		"&token=" + token
	req = httptest.NewRequest("GET", url, nil)

	w, _, vr := beforeVideoTest(req, false, true)

	assert.Equal(t, 200, w.Code)
	assert.Equal(t, int32(0), vr.StatusCode)
	assert.Equal(t, "get publish list successfully", vr.StatusMsg)
	assert.Equal(t, 1, len(vr.VideoList))
}
