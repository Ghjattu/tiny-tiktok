package jwt

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
	_ "github.com/joho/godotenv/autoload"
)

// secret key for signing the json web token.
var key string

// tokenLifespan will determine how long each json web token will last(hour).
var tokenLifespan int

// Retrieving secret key and token's lifespan from the .env file.
// If they don't exist, set them to the default value.
func init() {
	key = os.Getenv("SECRET_KEY")
	if key == "" {
		key = "secret-key"
	}

	tokenLifespan, _ = strconv.Atoi(os.Getenv("TOKEN_LIFESPAN"))
	if tokenLifespan == 0 {
		tokenLifespan = 1
	}
}

// GenerateToken generates a json web token using HMAC-SHA256.
//
//	@param userID int64
//	@param username string
//	@return string "token"
//	@return error
func GenerateToken(userID int64, username string) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID,
		"name":    username,
		"iat":     time.Now().Unix(),
		"exp":     time.Now().Add(time.Hour * time.Duration(tokenLifespan)).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(key))
}

// ValidateToken parses and validates a token using the HMAC signing method.
// see https://pkg.go.dev/github.com/golang-jwt/jwt/v5#example-Parse-Hmac.
//
//	@param tokenString string
//	@return int64 "user_id"
//	@return string "name"
//	@return error
func ValidateToken(tokenString string) (int64, string, error) {
	token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}

		return []byte(key), nil
	})

	if clams, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		userID := int64(clams["user_id"].(float64))
		name := clams["name"].(string)

		return userID, name, nil
	}

	return -1, "", err
}
