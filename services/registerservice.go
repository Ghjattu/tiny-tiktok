package services

import (
	"github.com/Ghjattu/tiny-tiktok/middleware"
	"github.com/Ghjattu/tiny-tiktok/models"
)

type RegisterService struct{}

// Register registers a new user.
// Return user_id, status_code, status_msg, token
func (rs *RegisterService) Register(username string, password string) (int64, int32, string, string) {
	// Check username and password is non-empty.
	if username == "" || password == "" {
		return -1, 1, "invalid username or password", ""
	}

	// Check if the username has been registered.
	user, _ := models.GetUserByName(username)
	if user.Name == username {
		return -1, 1, "the username has been registered", ""
	}

	// Create a new user.
	newUser := &models.User{
		Name:     username,
		Password: password,
	}

	// Insert the new user into the database.
	returnedUser, err := models.CreateNewUser(newUser)
	if err != nil {
		return -1, 1, "failed to create a new user", ""
	}

	// Generate a token.
	token, err := middleware.GenerateToken(returnedUser.ID, returnedUser.Name)
	if err != nil {
		return -1, 1, "failed to generate a token", ""
	}

	return int64(returnedUser.ID), 0, "register successfully", token
}
