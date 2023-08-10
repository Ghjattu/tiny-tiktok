package models

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateNewFavoriteRel(t *testing.T) {
	InitDatabase(true)

	// Create a new favorite rel.
	fr := &FavoriteRel{
		UserID:  1,
		VideoID: 1,
	}

	favoriteRel, err := CreateNewFavoriteRel(fr)
	if err != nil {
		t.Fatalf("Error when creating a new favorite rel: %v", err)
	}

	assert.Equal(t, fr.UserID, favoriteRel.UserID)
	assert.Equal(t, fr.VideoID, favoriteRel.VideoID)
}

func TestDeleteFavoriteRel(t *testing.T) {
	InitDatabase(true)

	// Create a new favorite rel.
	fr := &FavoriteRel{
		UserID:  1,
		VideoID: 1,
	}

	_, err := CreateNewFavoriteRel(fr)
	if err != nil {
		t.Fatalf("Error when creating a new favorite rel: %v", err)
	}

	// Delete the favorite rel.
	deletedRows, err := DeleteFavoriteRel(fr.UserID, fr.VideoID)
	if err != nil {
		t.Fatalf("Error when deleting the favorite rel: %v", err)
	}

	// Check the favorite rel is deleted.
	assert.Equal(t, int64(1), deletedRows)
}

func TestCheckFavoriteRelExist(t *testing.T) {
	InitDatabase(true)

	// Create a new favorite rel.
	fr := &FavoriteRel{
		UserID:  1,
		VideoID: 1,
	}

	_, err := CreateNewFavoriteRel(fr)
	if err != nil {
		t.Fatalf("Error when creating a new favorite rel: %v", err)
	}

	// Check the favorite rel exist.
	exist, err := CheckFavoriteRelExist(fr.UserID, fr.VideoID)
	if err != nil {
		t.Fatalf("Error when checking the favorite rel exist: %v", err)
	}

	assert.True(t, exist)
}

func TestGetFavoriteCountByVideoID(t *testing.T) {
	InitDatabase(true)

	// Create a new favorite rel.
	fr := &FavoriteRel{
		UserID:  1,
		VideoID: 1,
	}

	_, err := CreateNewFavoriteRel(fr)
	if err != nil {
		t.Fatalf("Error when creating a new favorite rel: %v", err)
	}

	// Get favorite count by video id.
	count, err := GetFavoriteCountByVideoID(fr.VideoID)
	if err != nil {
		t.Fatalf("Error when getting favorite count by video id: %v", err)
	}

	assert.Equal(t, int64(1), count)
}
