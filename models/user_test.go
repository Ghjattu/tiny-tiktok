package models

import (
	"os"
	"testing"

	_ "github.com/joho/godotenv/autoload"
	"github.com/stretchr/testify/assert"
)

func beforeTest() {
	os.Setenv("MYSQL_USERNAME", "root")
	os.Setenv("MYSQL_PASSWORD", "")
	os.Setenv("MYSQL_IP", "127.0.0.1")
	os.Setenv("MYSQL_PORT", "3306")
	os.Setenv("MYSQL_DB_NAME", "tiktok_test")
}

func TestCreateNewUser(t *testing.T) {
	beforeTest()

	InitDatabase()

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
	beforeTest()

	InitDatabase()

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
