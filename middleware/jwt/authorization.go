package jwt

import (
	"net/http"

	"github.com/Ghjattu/tiny-tiktok/models"
	"github.com/gin-gonic/gin"
)

type Response struct {
	StatusCode int32  `json:"status_code"`
	StatusMsg  string `json:"status_msg"`
}

type UserResponse struct {
	Response
	User *models.User `json:"user"`
}

// AuthorizeGet is a middleware that checks if the token is valid
// before GET requests.
// If the token is valid, it sets the user id and name to the context.
func AuthorizeGet() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := c.Query("token")

		// Parse the token.
		userID, name, err := ValidateToken(tokenString)

		// If the token is invalid, return an error.
		if err != nil {
			c.JSON(http.StatusUnauthorized, Response{
				StatusCode: 1,
				StatusMsg:  "invalid token",
			})
			c.Abort()
			return
		}

		// If the token is valid, set the user_id and name to the context.
		c.Set("user_id", userID)
		c.Set("username", name)
		c.Next()
	}
}

// AuthorizePost is a middleware that checks if the token is valid
// before POST requests.
// If the token is valid, it sets the user id and name to the context.
func AuthorizePost() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := c.PostForm("token")

		// Parse the token.
		userID, name, err := ValidateToken(tokenString)

		// If the token is invalid, return an error.
		if err != nil {
			c.JSON(http.StatusUnauthorized, Response{
				StatusCode: 1,
				StatusMsg:  "invalid token",
			})
			c.Abort()
			return
		}

		// If the token is valid, set the user_id and name to the context.
		c.Set("user_id", userID)
		c.Set("username", name)
		c.Next()
	}
}

// AuthorizeFeed is a middleware that checks if the token is valid
// before a feed request.
// If the token is empty, it sets the user id equal to -1 and name to empty string.
// else, check if the token is valid.
func AuthorizeFeed() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := c.Query("token")

		// If the token is empty.
		if tokenString == "" {
			c.Set("user_id", int64(-1))
			c.Set("username", "")
			c.Next()
			return
		}

		// Parse the token.
		userID, name, err := ValidateToken(tokenString)

		// If the token is invalid, return an error.
		if err != nil {
			c.JSON(http.StatusUnauthorized, Response{
				StatusCode: 1,
				StatusMsg:  "invalid token",
			})
			c.Abort()
			return
		}

		// If the token is valid, set the user_id and name to the context.
		c.Set("user_id", userID)
		c.Set("username", name)
		c.Next()
	}
}
