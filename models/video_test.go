package models

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestCreateNewVideo(t *testing.T) {
	InitDatabase(true)

	testVideo := &Video{
		AuthorID:    1,
		PublishTime: time.Now(),
		PlayUrl:     "test",
		Title:       "test",
	}

	returnedVideo, err := CreateNewVideo(testVideo)
	if err != nil {
		t.Fatalf("Error when creating a new video: %v", err)
	}

	assert.Equal(t, testVideo.AuthorID, returnedVideo.AuthorID)
	assert.Equal(t, testVideo.PlayUrl, returnedVideo.PlayUrl)
	assert.Equal(t, testVideo.Title, returnedVideo.Title)
}

func TestGetVideoListByUserID(t *testing.T) {
	InitDatabase(true)

	// Create a new test user.
	testUser, _ := createTestUser("test", "test")

	// Create a new test video.
	testVideo, _ := createTestVideo(testUser.ID, time.Now(), "test")

	// Get video list by test user id.
	videoList, err := GetVideoListByUserID(testUser.ID)
	if err != nil {
		t.Fatalf("Error when getting video list: %v", err)
	}

	assert.Equal(t, 1, len(videoList))
	assert.Equal(t, testVideo.AuthorID, videoList[0].Author.ID)
	assert.Equal(t, testVideo.PlayUrl, videoList[0].PlayUrl)
	assert.Equal(t, testVideo.Title, videoList[0].Title)
}

func TestGetMost30Videos(t *testing.T) {
	InitDatabase(true)

	// Create a new test user.
	testUser, _ := createTestUser("test", "test")

	// Construct three timestamps.
	videoOneTimestamp := time.Now()
	middleTimestamp := time.Now().Add(time.Second * 5)
	videoTwoTimestamp := time.Now().Add(time.Second * 10)

	// Create two new test videos.
	testVideoOne, _ := createTestVideo(testUser.ID, videoOneTimestamp, "testOne")
	createTestVideo(testUser.ID, videoTwoTimestamp, "testTwo")

	// Check the results.
	videoList, earliestTime, err := GetMost30Videos(middleTimestamp)
	if err != nil {
		t.Fatalf("Error when getting most 30 videos: %v", err)
	}

	assert.Equal(t, 1, len(videoList))
	assert.Equal(t, testVideoOne.Title, videoList[0].Title)
	assert.Equal(t, videoOneTimestamp.Unix(), earliestTime.Unix())
}

func TestGetMost30VideosWithEmptyVideoList(t *testing.T) {
	InitDatabase(true)

	videoList, _, err := GetMost30Videos(time.Now())
	if err != nil {
		t.Fatalf("Error when getting most 30 videos: %v", err)
	}

	assert.Equal(t, 0, len(videoList))
}
