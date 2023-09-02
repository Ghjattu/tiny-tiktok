package models

import (
	"testing"

	"github.com/Ghjattu/tiny-tiktok/utils"
	"github.com/stretchr/testify/assert"
)

func benchmarkUserSetup() {
	// Create test user.
	for i := 1; i <= 66; i++ {
		CreateTestUser(utils.GenerateRandomString(3), "test")
	}
}

func TestCreateNewUser(t *testing.T) {
	InitDatabase(true)

	testUser := &User{
		Name:     "test",
		Password: "test",
	}

	returnedUser, err := CreateNewUser(testUser)
	if err != nil {
		t.Fatalf("Error when creating a new user: %v", err)
	}

	assert.Equal(t, testUser.Name, returnedUser.Name)
}

func TestGetUserByName(t *testing.T) {
	InitDatabase(true)

	// Create a new test user.
	testUser, _ := CreateTestUser("test", "test")

	returnedUser, err := GetUserByName(testUser.Name)
	if err != nil {
		t.Fatalf("Error when getting a user: %v", err)
	}

	assert.Equal(t, testUser.Name, returnedUser.Name)
}

func TestGetUserByUserID(t *testing.T) {
	InitDatabase(true)

	// Create a new test user.
	testUser, _ := CreateTestUser("test", "test")

	returnedUser, err := GetUserByUserID(testUser.ID)
	if err != nil {
		t.Fatalf("Error when getting a user: %v", err)
	}

	assert.Equal(t, testUser.Name, returnedUser.Name)
}

func BenchmarkGetUserByName(b *testing.B) {
	benchmarkUserSetup()

	username := utils.GenerateRandomString(6)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		GetUserByName(username)
	}
}

func BenchmarkGetUserByUserID(b *testing.B) {
	benchmarkUserSetup()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		GetUserByUserID(int64(6))
	}
}
