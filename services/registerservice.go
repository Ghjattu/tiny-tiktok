package services

import (
	"github.com/Ghjattu/tiny-tiktok/bloomfilter"
	"github.com/Ghjattu/tiny-tiktok/middleware/jwt"
	"github.com/Ghjattu/tiny-tiktok/models"
	"golang.org/x/crypto/bcrypt"
)

// RegisterService implements RegisterInterface.
type RegisterService struct{}

// Register registers a new user.
// Return user_id, status_code, status_msg, token
func (rs *RegisterService) Register(username, password string) (int64, int32, string, string) {
	// Check username and password is non-empty.
	if username == "" || password == "" {
		return -1, 1, "invalid username or password", ""
	}

	// Check password length.
	if len([]rune(password)) < 6 {
		return -1, 1, "password is too short", ""
	}

	// Check username and password length.
	if len([]rune(username)) > 32 || len([]rune(password)) > 32 {
		return -1, 1, "username or password is too long", ""
	}

	// Check if the username has been registered.
	user, err := models.GetUserByName(username)
	if err == nil && user.Name == username {
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
		if err == bcrypt.ErrPasswordTooLong {
			return -1, 1, "password length exceeds 72 bytes", ""
		}
		return -1, 1, "failed to create a new user", ""
	}

	// Generate a token.
	token, err := jwt.GenerateToken(returnedUser.ID, returnedUser.Name)
	if err != nil {
		return -1, 1, "failed to generate a token", ""
	}

	// Add the user id to bloom filter.
	bloomfilter.Add(bloomfilter.UserBloomFilter, returnedUser.ID)

	return returnedUser.ID, 0, "register successfully", token
}
