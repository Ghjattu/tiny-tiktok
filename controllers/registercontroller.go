package controllers

import (
	"net/http"

	"github.com/Ghjattu/tiny-tiktok/services"
	"github.com/gin-gonic/gin"
)

type RegisterResponse struct {
	Response
	UserID int64  `json:"user_id"`
	Token  string `json:"token"`
}

// Endpoint: /douyin/user/register/
func Register(c *gin.Context) {
	username := c.Query("username")
	password := c.Query("password")

	rs := &services.RegisterService{}
	userID, statusCode, statusMsg, token := rs.Register(username, password)

	if statusCode == 0 {
		c.JSON(http.StatusOK, RegisterResponse{
			Response: Response{
				StatusCode: statusCode,
				StatusMsg:  statusMsg,
			},
			UserID: userID,
			Token:  token,
		})

		// If register successfully, set the user_id and name to the context.
		c.Set("user_id", userID)
		c.Set("username", username)

		return
	}

	c.JSON(http.StatusOK, RegisterResponse{
		Response: Response{
			StatusCode: statusCode,
			StatusMsg:  statusMsg,
		},
		UserID: -1,
		Token:  "",
	})
}
