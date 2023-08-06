package models

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

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
	testUser, _ := createTestUser("test", "test")

	returnedUser, err := GetUserByName(testUser.Name)
	if err != nil {
		t.Fatalf("Error when getting a user: %v", err)
	}

	assert.Equal(t, testUser.Name, returnedUser.Name)
}

func TestGetUserByUserID(t *testing.T) {
	InitDatabase(true)

	// Create a new test user.
	testUser, _ := createTestUser("test", "test")

	returnedUser, err := GetUserByUserID(testUser.ID)
	if err != nil {
		t.Fatalf("Error when getting a user: %v", err)
	}

	assert.Equal(t, testUser.Name, returnedUser.Name)
}
