package models

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func benchmarkFavoriteRelSetup() {
	InitDatabase(true)

	for i := 0; i < 66; i++ {
		CreateTestFavoriteRel(int64(i%10+1), int64(i%10+2))
	}
}

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
	testFavoriteRel, _ := CreateTestFavoriteRel(1, 1)

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
	testFavoriteRel, _ := CreateTestFavoriteRel(1, 1)

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
	testFavoriteRel, _ := CreateTestFavoriteRel(1, 1)

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
	testFavoriteRel, _ := CreateTestFavoriteRel(1, 1)

	// Get favorite count by user id.
	count, err := GetFavoriteCountByUserID(testFavoriteRel.UserID)
	if err != nil {
		t.Fatalf("Error when getting favorite count by user id: %v", err)
	}

	assert.Equal(t, int64(1), count)
}

func TestGetFavoriteVideoIDListByUserID(t *testing.T) {
	InitDatabase(true)

	// Create a new favorite rel.
	testFavoriteRel, _ := CreateTestFavoriteRel(1, 1)

	// Get favorite video id list by user id.
	videoIDList, err := GetFavoriteVideoIDListByUserID(testFavoriteRel.UserID)
	if err != nil {
		t.Fatalf("Error when getting favorite video id list by user id: %v", err)
	}

	assert.Equal(t, 1, len(videoIDList))
	assert.Equal(t, []int64{1}, videoIDList)
}

func BenchmarkCheckFavoriteRelExist(b *testing.B) {
	benchmarkFavoriteRelSetup()

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		CheckFavoriteRelExist(6, 6)
	}
}

func BenchmarkGetFavoriteCountByVideoID(b *testing.B) {
	benchmarkFavoriteRelSetup()

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		GetFavoriteCountByVideoID(6)
	}
}

func BenchmarkGetFavoriteCountByUserID(b *testing.B) {
	benchmarkFavoriteRelSetup()

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		GetFavoriteCountByUserID(6)
	}
}

func BenchmarkGetFavoriteVideoIDListByUserID(b *testing.B) {
	benchmarkFavoriteRelSetup()

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		GetFavoriteVideoIDListByUserID(6)
	}
}
