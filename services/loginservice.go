package services

import (
	"github.com/Ghjattu/tiny-tiktok/bloomfilter"
	"github.com/Ghjattu/tiny-tiktok/middleware/jwt"
	"github.com/Ghjattu/tiny-tiktok/models"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// LOginService implements LoginInterface.
type LoginService struct{}

// Login logs in a user.
// Return user_id, status_code, status_msg, token
func (ls *LoginService) Login(username, password string) (int64, int32, string, string) {
	// Check username and password length.
	if len([]rune(username)) > 32 || len([]rune(password)) > 32 {
		return -1, 1, "username or password is too long", ""
	}

	// Check username exists.
	user, err := models.GetUserByName(username)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return -1, 1, "username not found", ""
		}
		return -1, 1, "failed to get user", ""
	}

	// Compare the password with the hashed password.
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		if err == bcrypt.ErrHashTooShort {
			return -1, 1, "invalid hashed password", ""
		}
		if err == bcrypt.ErrMismatchedHashAndPassword {
			return -1, 1, "wrong password", ""
		}
		return -1, 1, "failed to compare password", ""
	}

	// Generate a token.
	token, err := jwt.GenerateToken(user.ID, user.Name)
	if err != nil {
		return -1, 1, "failed to generate a token", ""
	}

	// Add the user id to bloom filter.
	bloomfilter.Add(bloomfilter.UserBloomFilter, user.ID)

	return user.ID, 0, "login successfully", token
}
