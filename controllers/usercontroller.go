package controllers

import (
	"net/http"

	"github.com/Ghjattu/tiny-tiktok/models"
	"github.com/Ghjattu/tiny-tiktok/services"
	"github.com/gin-gonic/gin"
)

type UserResponse struct {
	Response
	User *models.UserDetail `json:"user"`
}

// Endpoint: /douyin/user/
func GetUserByUserIDAndToken(c *gin.Context) {
	queryUserID := c.GetInt64("user_id")
	currentUserID := c.GetInt64("current_user_id")

	us := &services.UserService{}
	statusCode, statusMsg, user := us.GetUserDetailByUserID(currentUserID, queryUserID)

	c.JSON(http.StatusOK, UserResponse{
		Response: Response{
			StatusCode: statusCode,
			StatusMsg:  statusMsg,
		},
		User: user,
	})
}
