package models

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestCreateTestUser(t *testing.T) {
	InitDatabase(true)

	testUser := &User{
		Name:     "test",
		Password: "test",
	}

	returnedUser, err := createTestUser(testUser.Name, testUser.Password)
	if err != nil {
		t.Fatalf("Error when creating a new user: %v", err)
	}

	assert.Equal(t, testUser.Name, returnedUser.Name)
}

func TestCreateTestVideo(t *testing.T) {
	InitDatabase(true)

	testVideo := &Video{
		AuthorID:    1,
		PublishTime: time.Now(),
		PlayUrl:     "test",
		Title:       "test",
	}

	returnedVideo, err := createTestVideo(testVideo.AuthorID, testVideo.PublishTime, testVideo.Title)
	if err != nil {
		t.Fatalf("Error when creating a new video: %v", err)
	}

	assert.Equal(t, testVideo.AuthorID, returnedVideo.AuthorID)
	assert.Equal(t, testVideo.PlayUrl, returnedVideo.PlayUrl)
	assert.Equal(t, testVideo.Title, returnedVideo.Title)
}
