package models

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateNewVideo(t *testing.T) {
	InitDatabase(true)

	testVideo := &Video{
		AuthorID:   1,
		AuthorName: "test",
		PlayUrl:    "test",
		Title:      "test",
	}

	returnedVideo, err := CreateNewVideo(testVideo)
	if err != nil {
		t.Fatalf("Error when creating a new video: %v", err)
	}

	assert.Equal(t, testVideo.AuthorID, returnedVideo.AuthorID)
	assert.Equal(t, testVideo.AuthorName, returnedVideo.AuthorName)
	assert.Equal(t, testVideo.PlayUrl, returnedVideo.PlayUrl)
	assert.Equal(t, testVideo.Title, returnedVideo.Title)
}
