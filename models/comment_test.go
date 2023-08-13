package models

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestCreateNewComment(t *testing.T) {
	InitDatabase(true)

	testComment := &Comment{
		UserID:     1,
		VideoID:    1,
		CreateDate: time.Now(),
	}

	returnedComment, err := CreateNewComment(testComment)
	if err != nil {
		t.Errorf("Error when creating a comment: %v", err)
	}

	assert.Equal(t, testComment.UserID, returnedComment.UserID)
	assert.Equal(t, testComment.VideoID, returnedComment.VideoID)
}

func TestGetCommentCountByVideoID(t *testing.T) {
	InitDatabase(true)

	// Create a test comment.
	CreateTestComment(1, 1)

	count, err := GetCommentCountByVideoID(1)
	if err != nil {
		t.Errorf("Error when getting comment count: %v", err)
	}

	assert.Equal(t, int64(1), count)
}

func TestDeleteCommentByCommentID(t *testing.T) {
	InitDatabase(true)

	// Create a test comment.
	testComment, _ := CreateTestComment(1, 1)

	count, err := DeleteCommentByCommentID(testComment.ID)
	if err != nil {
		t.Errorf("Error when deleting a comment: %v", err)
	}

	assert.Equal(t, int64(1), count)
}

func TestGetCommentByCommentID(t *testing.T) {
	InitDatabase(true)

	// Create a test comment.
	testComment, _ := CreateTestComment(1, 1)

	comment, err := GetCommentByCommentID(testComment.ID)
	if err != nil {
		t.Errorf("Error when getting a comment: %v", err)
	}

	assert.Equal(t, testComment.UserID, comment.UserID)
	assert.Equal(t, testComment.VideoID, comment.VideoID)
}

func TestGetCommentListByVideoID(t *testing.T) {
	InitDatabase(true)

	// Create a test comment.
	testComment, _ := CreateTestComment(1, 1)

	commentList, err := GetCommentListByVideoID(testComment.VideoID)
	if err != nil {
		t.Errorf("Error when getting comment list: %v", err)
	}

	assert.Equal(t, 1, len(commentList))
	assert.Equal(t, testComment.UserID, commentList[0].UserID)
	assert.Equal(t, testComment.VideoID, commentList[0].VideoID)
}
