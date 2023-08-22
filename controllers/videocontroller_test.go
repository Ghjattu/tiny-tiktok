package controllers

import (
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPublishNewVideo(t *testing.T) {
	setup()

	t.Run("empty title", func(t *testing.T) {
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
	})

	t.Run("create new video successfully", func(t *testing.T) {
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
	})
}

func TestGetPublishListByAuthorID(t *testing.T) {
	setup()

	// Get publish list by author id.
	url := "http://127.0.0.1/douyin/publish/list/?user_id=" + userIDStr +
		"&token=" + token
	req := httptest.NewRequest("GET", url, nil)

	w, r := sendRequest(req)
	res := r.(*PublishListResponse)

	assert.Equal(t, 200, w.Code)
	assert.Equal(t, int32(0), res.StatusCode)
	assert.Equal(t, 0, len(res.VideoList))
}

func TestFeed(t *testing.T) {
	setup()

	t.Run("invalid latest time", func(t *testing.T) {
		req := httptest.NewRequest("GET", "http://127.0.0.1/douyin/feed/?latest_time=abc", nil)

		w, r := sendRequest(req)
		res := r.(*FeedResponse)

		assert.Equal(t, 200, w.Code)
		assert.Equal(t, int32(1), res.StatusCode)
		assert.Equal(t, 0, len(res.VideoList))
	})

	t.Run("get video list successfully", func(t *testing.T) {
		req := httptest.NewRequest("GET", "http://127.0.0.1/douyin/feed/", nil)

		w, r := sendRequest(req)
		res := r.(*FeedResponse)

		assert.Equal(t, 200, w.Code)
		assert.Equal(t, int32(0), res.StatusCode)
		assert.Equal(t, 0, len(res.VideoList))
	})
}
