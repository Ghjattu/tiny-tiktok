package models

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func benchmarkFollowRelSetup() {
	InitDatabase(true)

	// Create test follow relationship.
	for i := 1; i <= 66; i++ {
		CreateTestFollowRel(int64(i%10+1), int64(i%10+2))
	}
}

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

func TestGetFollowingListByUserID(t *testing.T) {
	InitDatabase(true)

	// Create a test follow relationship.
	testFollowRel, _ := CreateTestFollowRel(1, 2)

	list, err := GetFollowingListByUserID(testFollowRel.FollowerID)
	if err != nil {
		t.Fatalf("Error getting following list: %v", err)
	}

	assert.Equal(t, 1, len(list))
	assert.Equal(t, testFollowRel.FollowingID, list[0])
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

func TestGetFollowerListByUserID(t *testing.T) {
	InitDatabase(true)

	// Create a test follow relationship.
	testFollowRel, _ := CreateTestFollowRel(1, 2)

	list, err := GetFollowerListByUserID(testFollowRel.FollowingID)
	if err != nil {
		t.Fatalf("Error getting follower list: %v", err)
	}

	assert.Equal(t, 1, len(list))
	assert.Equal(t, testFollowRel.FollowerID, list[0])
}

func BenchmarkCheckFollowRelExist(b *testing.B) {
	benchmarkFollowRelSetup()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		CheckFollowRelExist(1, 2)
	}
}

func BenchmarkGetFollowingCountByUserID(b *testing.B) {
	benchmarkFollowRelSetup()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		GetFollowingCountByUserID(1)
	}
}

func BenchmarkGetFollowingListByUserID(b *testing.B) {
	benchmarkFollowRelSetup()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		GetFollowingListByUserID(1)
	}
}

func BenchmarkGetFollowerCountByUserID(b *testing.B) {
	benchmarkFollowRelSetup()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		GetFollowerCountByUserID(2)
	}
}

func BenchmarkGetFollowerListByUserID(b *testing.B) {
	benchmarkFollowRelSetup()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		GetFollowerListByUserID(2)
	}
}
