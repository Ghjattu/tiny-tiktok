package models

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateNewFollowRel(t *testing.T) {
	InitDatabase(true)

	fr := &FollowRel{
		FollowerID:  1,
		FollowingID: 2,
	}

	returnedFr, err := CreateNewFollowRel(fr)
	if err != nil {
		t.Fatalf("Error creating new follow relationship: %v", err)
	}

	assert.Equal(t, fr.FollowerID, returnedFr.FollowerID)
	assert.Equal(t, fr.FollowingID, returnedFr.FollowingID)
}

func TestDeleteFollowRel(t *testing.T) {
	InitDatabase(true)

	// Create a test follow relationship.
	testFollowRel, _ := CreateTestFollowRel(1, 2)

	count, err := DeleteFollowRel(testFollowRel.FollowerID, testFollowRel.FollowingID)
	if err != nil {
		t.Fatalf("Error deleting follow relationship: %v", err)
	}

	assert.Equal(t, int64(1), count)
}

func TestCheckFollowRelExist(t *testing.T) {
	InitDatabase(true)

	// Create a test follow relationship.
	testFollowRel, _ := CreateTestFollowRel(1, 2)

	exist, err := CheckFollowRelExist(testFollowRel.FollowerID, testFollowRel.FollowingID)
	if err != nil {
		t.Fatalf("Error checking follow relationship exist: %v", err)
	}

	assert.Equal(t, true, exist)
}

func TestGetFollowingCountByUserID(t *testing.T) {
	InitDatabase(true)

	// Create a test follow relationship.
	testFollowRel, _ := CreateTestFollowRel(1, 2)

	count, err := GetFollowingCountByUserID(testFollowRel.FollowerID)
	if err != nil {
		t.Fatalf("Error getting following count: %v", err)
	}

	assert.Equal(t, int64(1), count)

	count, _ = GetFollowingCountByUserID(testFollowRel.FollowingID)

	assert.Equal(t, int64(0), count)
}

func TestGetFollowerCountByUserID(t *testing.T) {
	InitDatabase(true)

	// Create a test follow relationship.
	testFollowRel, _ := CreateTestFollowRel(1, 2)

	count, err := GetFollowerCountByUserID(testFollowRel.FollowingID)
	if err != nil {
		t.Fatalf("Error getting follower count: %v", err)
	}

	assert.Equal(t, int64(1), count)

	count, _ = GetFollowerCountByUserID(testFollowRel.FollowerID)

	assert.Equal(t, int64(0), count)
}
