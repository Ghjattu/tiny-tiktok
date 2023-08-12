package models

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateNewFavoriteRel(t *testing.T) {
	InitDatabase(true)

	// Create a new favorite rel.
	testFavoriteRel := &FavoriteRel{
		UserID:  1,
		VideoID: 1,
	}

	favoriteRel, err := CreateNewFavoriteRel(testFavoriteRel)
	if err != nil {
		t.Fatalf("Error when creating a new favorite rel: %v", err)
	}

	assert.Equal(t, testFavoriteRel.UserID, favoriteRel.UserID)
	assert.Equal(t, testFavoriteRel.VideoID, favoriteRel.VideoID)
}

func TestDeleteFavoriteRel(t *testing.T) {
	InitDatabase(true)

	// Create a new favorite rel.
	testFavoriteRel, _ := createTestFavoriteRel(1, 1)

	// Delete the favorite rel.
	deletedRows, err := DeleteFavoriteRel(testFavoriteRel.UserID, testFavoriteRel.VideoID)
	if err != nil {
		t.Fatalf("Error when deleting the favorite rel: %v", err)
	}

	// Check the favorite rel is deleted.
	assert.Equal(t, int64(1), deletedRows)
}

func TestCheckFavoriteRelExist(t *testing.T) {
	InitDatabase(true)

	// Create a new favorite rel.
	testFavoriteRel, _ := createTestFavoriteRel(1, 1)

	// Check the favorite rel exist.
	exist, err := CheckFavoriteRelExist(testFavoriteRel.UserID, testFavoriteRel.VideoID)
	if err != nil {
		t.Fatalf("Error when checking the favorite rel exist: %v", err)
	}

	assert.True(t, exist)
}

func TestGetFavoriteCountByVideoID(t *testing.T) {
	InitDatabase(true)

	// Create a new favorite rel.
	testFavoriteRel, _ := createTestFavoriteRel(1, 1)

	// Get favorite count by video id.
	count, err := GetFavoriteCountByVideoID(testFavoriteRel.VideoID)
	if err != nil {
		t.Fatalf("Error when getting favorite count by video id: %v", err)
	}

	assert.Equal(t, int64(1), count)
}

func TestGetFavoriteCountByUserID(t *testing.T) {
	InitDatabase(true)

	// Create a new favorite rel.
	testFavoriteRel, _ := createTestFavoriteRel(1, 1)

	// Get favorite count by user id.
	count, err := GetFavoriteCountByUserID(testFavoriteRel.UserID)
	if err != nil {
		t.Fatalf("Error when getting favorite count by user id: %v", err)
	}

	assert.Equal(t, int64(1), count)
}
