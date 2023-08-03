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

	testUser := &User{
		Name:     "test1",
		Password: "test1",
	}

	_, err := CreateNewUser(testUser)
	if err != nil {
		t.Fatalf("Error when creating a new user: %v", err)
	}

	returnedUser, err := GetUserByName(testUser.Name)
	if err != nil {
		t.Fatalf("Error when getting a user: %v", err)
	}

	assert.Equal(t, testUser.Name, returnedUser.Name)
}

func TestGetUserByUserID(t *testing.T) {
	InitDatabase(true)

	testUser := &User{
		Name:     "test",
		Password: "test",
	}

	_, err := CreateNewUser(testUser)
	if err != nil {
		t.Fatalf("Error when creating a new user: %v", err)
	}

	returnedUser, err := GetUserByUserID(testUser.ID)
	if err != nil {
		t.Fatalf("Error when getting a user: %v", err)
	}

	assert.Equal(t, testUser.Name, returnedUser.Name)
}
