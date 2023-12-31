package models

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func benchmarkVideoSetup() {
	InitDatabase(true)

	for i := 0; i < 66; i++ {
		CreateTestVideo(int64(i%10+1), time.Now(), "test")
	}
}

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

func TestGetVideoIDListByAuthorID(t *testing.T) {
	InitDatabase(true)

	// Create a test video.
	testVideo, _ := CreateTestVideo(1, time.Now(), "test")

	videoIDList, err := GetVideoIDListByAuthorID(testVideo.AuthorID)
	if err != nil {
		t.Fatalf("Error when getting video id list by author id: %v", err)
	}

	assert.Equal(t, 1, len(videoIDList))
	assert.Equal(t, testVideo.ID, videoIDList[0])
}

func TestGetAuthorIDByVideoID(t *testing.T) {
	InitDatabase(true)

	// Create a test video.
	testVideo, _ := CreateTestVideo(1, time.Now(), "test")

	authorID, err := GetAuthorIDByVideoID(testVideo.ID)
	if err != nil {
		t.Fatalf("Error when getting author id by video id: %v", err)
	}

	assert.Equal(t, testVideo.AuthorID, authorID)
}

func TestGetMost30Videos(t *testing.T) {
	InitDatabase(true)

	// Create a new test user.
	testUser, _ := CreateTestUser("test", "test")

	// Construct three timestamps.
	videoOneTimestamp := time.Now()
	middleTimestamp := time.Now().Add(time.Second * 5)
	videoTwoTimestamp := time.Now().Add(time.Second * 10)

	// Create two new test videos.
	CreateTestVideo(testUser.ID, videoOneTimestamp, "testOne")
	CreateTestVideo(testUser.ID, videoTwoTimestamp, "testTwo")

	// Check the results.
	videoIDList, earliestTime, err := GetMost30Videos(middleTimestamp)
	if err != nil {
		t.Fatalf("Error when getting most 30 videos: %v", err)
	}

	assert.Equal(t, 1, len(videoIDList))
	assert.Equal(t, videoOneTimestamp.Unix(), earliestTime.Unix())
}

func TestGetMost30VideosWithEmptyVideoList(t *testing.T) {
	InitDatabase(true)

	videoList, _, err := GetMost30Videos(time.Now())
	if err != nil {
		t.Fatalf("Error when getting most 30 videos: %v", err)
	}

	assert.Equal(t, 0, len(videoList))
}

func TestGetVideoByID(t *testing.T) {
	InitDatabase(true)

	// Create a new test video.
	testVideo, _ := CreateTestVideo(1, time.Now(), "test")

	returnedVideo, err := GetVideoByID(testVideo.ID)
	if err != nil {
		t.Fatalf("Error when getting video by id: %v", err)
	}

	assert.Equal(t, testVideo.AuthorID, returnedVideo.AuthorID)
	assert.Equal(t, testVideo.Title, returnedVideo.Title)
}

func TestGetVideoCountByAuthorID(t *testing.T) {
	InitDatabase(true)

	// Create a new test video.
	CreateTestVideo(1, time.Now(), "test")

	count, err := GetVideoCountByAuthorID(1)
	if err != nil {
		t.Fatalf("Error when getting video count by author id: %v", err)
	}

	assert.Equal(t, int64(1), count)
}

func BenchmarkGetVideoIDListByAuthorID(b *testing.B) {
	benchmarkVideoSetup()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		GetVideoIDListByAuthorID(6)
	}
}

func BenchmarkGetAuthorIDByVideoID(b *testing.B) {
	benchmarkVideoSetup()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		GetAuthorIDByVideoID(6)
	}
}

func BenchmarkGetVideoByID(b *testing.B) {
	benchmarkVideoSetup()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		GetVideoByID(int64(6))
	}
}

func BenchmarkGetVideoCountByAuthorID(b *testing.B) {
	benchmarkVideoSetup()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		GetVideoCountByAuthorID(int64(6))
	}
}
