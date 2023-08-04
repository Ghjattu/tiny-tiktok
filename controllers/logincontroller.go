package controllers

import (
	"github.com/Ghjattu/tiny-tiktok/services"
	"github.com/gin-gonic/gin"
)

type LoginResponse struct {
	Response
	UserID int64  `json:"user_id"`
	Token  string `json:"token"`
}

// Endpoint: /douyin/user/login/
func Login(c *gin.Context) {
	username := c.Query("username")
	password := c.Query("password")

	ls := &services.LoginService{}
	userID, statusCode, statusMsg, token := ls.Login(username, password)

	// Login successfully.
	if statusCode == 0 {
		c.JSON(200, LoginResponse{
			Response: Response{
				StatusCode: statusCode,
				StatusMsg:  statusMsg,
			},
			UserID: userID,
			Token:  token,
		})
		return
	}

	c.JSON(200, LoginResponse{
		Response: Response{
			StatusCode: statusCode,
			StatusMsg:  statusMsg,
		},
	})
}
