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
	testUser := &User{
		Name:     "test",
		Password: "test",
	}

	returnedTestUser, err := CreateNewUser(testUser)
	if err != nil {
		t.Fatalf("Error when creating a new user: %v", err)
	}

	assert.Equal(t, testUser.Name, returnedTestUser.Name)

	// Create a new test video.
	testVideo := &Video{
		AuthorID:    returnedTestUser.ID,
		PublishTime: time.Now(),
		PlayUrl:     "test",
		Title:       "test",
	}

	returnedTestVideo, err := CreateNewVideo(testVideo)
	if err != nil {
		t.Fatalf("Error when creating a new video: %v", err)
	}

	assert.Equal(t, testVideo.AuthorID, returnedTestVideo.AuthorID)

	// Get video list by test user id.
	videoList, err := GetVideoListByUserID(returnedTestUser.ID)
	if err != nil {
		t.Fatalf("Error when getting video list: %v", err)
	}

	assert.Equal(t, 1, len(videoList))
	assert.Equal(t, testVideo.AuthorID, videoList[0].Author.ID)
	assert.Equal(t, testVideo.PlayUrl, videoList[0].PlayUrl)
	assert.Equal(t, testVideo.Title, videoList[0].Title)
}
