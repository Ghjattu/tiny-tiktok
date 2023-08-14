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

	returnedUser, err := CreateTestUser(testUser.Name, testUser.Password)
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

	returnedVideo, err := CreateTestVideo(testVideo.AuthorID, testVideo.PublishTime, testVideo.Title)
	if err != nil {
		t.Fatalf("Error when creating a new video: %v", err)
	}

	assert.Equal(t, testVideo.AuthorID, returnedVideo.AuthorID)
	assert.Equal(t, testVideo.PlayUrl, returnedVideo.PlayUrl)
	assert.Equal(t, testVideo.Title, returnedVideo.Title)
}

func TestCreateTestFavoriteRel(t *testing.T) {
	InitDatabase(true)

	fr := &FavoriteRel{
		UserID:  1,
		VideoID: 1,
	}

	returnedFr, err := CreateTestFavoriteRel(fr.UserID, fr.VideoID)
	if err != nil {
		t.Fatalf("Error when creating a new favorite rel: %v", err)
	}

	assert.Equal(t, fr.UserID, returnedFr.UserID)
	assert.Equal(t, fr.VideoID, returnedFr.VideoID)
}

func TestCreateTestComment(t *testing.T) {
	InitDatabase(true)

	comment := &Comment{
		UserID:  1,
		VideoID: 1,
	}

	returnedComment, err := CreateTestComment(comment.UserID, comment.VideoID)
	if err != nil {
		t.Fatalf("Error when creating a new comment: %v", err)
	}

	assert.Equal(t, comment.UserID, returnedComment.UserID)
	assert.Equal(t, comment.VideoID, returnedComment.VideoID)
}

func TestCreateTestFollowRel(t *testing.T) {
	InitDatabase(true)

	fr := &FollowRel{
		FollowerID:  1,
		FollowingID: 2,
	}

	returnedFr, err := CreateTestFollowRel(fr.FollowerID, fr.FollowingID)
	if err != nil {
		t.Fatalf("Error when creating a new follow rel: %v", err)
	}

	assert.Equal(t, fr.FollowerID, returnedFr.FollowerID)
	assert.Equal(t, fr.FollowingID, returnedFr.FollowingID)
}
