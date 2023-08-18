package controllers

import (
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/Ghjattu/tiny-tiktok/models"
	"github.com/stretchr/testify/assert"
)

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

func TestPublishNewVideoWithEmptyTitle(t *testing.T) {
	models.Flush()

	// Register a new test user.
	_, _, token := registerTestUser("test", "123456")

	// Construct a test form.
	formFields := map[string]string{
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

	w, r := sendRequest(req)
	res := r.(*Response)

	assert.Equal(t, 200, w.Code)
	assert.Equal(t, int32(1), res.StatusCode)
	assert.Equal(t, "video title is empty", res.StatusMsg)
}

func TestPublishNewVideo(t *testing.T) {
	models.Flush()

	// Register a new test user.
	_, _, token := registerTestUser("test", "123456")

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

	w, r := sendRequest(req)
	res := r.(*Response)

	assert.Equal(t, 200, w.Code)
	assert.Equal(t, int32(0), res.StatusCode)
	assert.Equal(t, "create new video successfully", res.StatusMsg)
}

func TestGetPublishListByAuthorID(t *testing.T) {
	models.Flush()

	// Register a new test user.
	userID, _, token := registerTestUser("test", "123456")
	userIDStr := fmt.Sprintf("%d", userID)

	// Get publish list by author id.
	url := "http://127.0.0.1/douyin/publish/list/?user_id=" + userIDStr +
		"&token=" + token
	req := httptest.NewRequest("GET", url, nil)

	w, r := sendRequest(req)
	res := r.(*PublishListResponse)

	assert.Equal(t, 200, w.Code)
	assert.Equal(t, int32(0), res.StatusCode)
	assert.Equal(t, "get publish list successfully", res.StatusMsg)
	assert.Equal(t, 0, len(res.VideoList))
}

func TestFeedWithInvalidLatestTime(t *testing.T) {
	req := httptest.NewRequest("GET", "http://127.0.0.1/douyin/feed/?latest_time=abc", nil)

	w, r := sendRequest(req)
	res := r.(*FeedResponse)

	assert.Equal(t, 200, w.Code)
	assert.Equal(t, int32(1), res.StatusCode)
	assert.Equal(t, "invalid syntax", res.StatusMsg)
	assert.Equal(t, 0, len(res.VideoList))
}

func TestFeedWithEmptyLatestTime(t *testing.T) {
	models.Flush()

	req := httptest.NewRequest("GET", "http://127.0.0.1/douyin/feed/", nil)

	w, r := sendRequest(req)
	res := r.(*FeedResponse)

	assert.Equal(t, 200, w.Code)
	assert.Equal(t, int32(0), res.StatusCode)
	assert.Equal(t, "get most 30 videos successfully", res.StatusMsg)
	assert.Equal(t, 0, len(res.VideoList))
}
