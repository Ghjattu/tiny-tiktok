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

// AuthorizationGet is a middleware that checks if the token is valid
// before GET requests.
// If the token is valid, it sets the user id and name to the context.
func AuthorizationGet() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := c.Query("token")

		// Parse the token.
		userID, name, err := ValidateToken(tokenString)

		// If the token is invalid, return an error.
		if err != nil {
			c.JSON(http.StatusUnauthorized, UserResponse{
				Response: Response{
					StatusCode: 1,
					StatusMsg:  "invalid token",
				},
				User: nil,
			})
			c.Abort()
			return
		}

		// If the token is valid, set the user_id and name to the context.
		c.Set("user_id", userID)
		c.Set("name", name)
		c.Next()
	}
}
