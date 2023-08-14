package controllers

import (
	"net/http"

	"github.com/Ghjattu/tiny-tiktok/models"
	"github.com/Ghjattu/tiny-tiktok/services"
	"github.com/Ghjattu/tiny-tiktok/utils"
	"github.com/gin-gonic/gin"
)

type UserResponse struct {
	Response
	User *models.UserDetail `json:"user"`
}

// Endpoint: /douyin/user/
func GetUserByUserIDAndToken(c *gin.Context) {
	userIDString := c.Query("user_id")

	// Check user id is valid.
	statusCode, statusMsg, userID := utils.ParseInt64(userIDString)
	if statusCode == 1 {
		c.JSON(http.StatusBadRequest, UserResponse{
			Response: Response{
				StatusCode: statusCode,
				StatusMsg:  statusMsg,
			},
			User: nil,
		})
		return
	}

	currentUserID := c.GetInt64("user_id")

	us := &services.UserService{}
	statusCode, statusMsg, user := us.GetUserDetailByUserID(currentUserID, userID)

	c.JSON(http.StatusOK, UserResponse{
		Response: Response{
			StatusCode: statusCode,
			StatusMsg:  statusMsg,
		},
		User: user,
	})
}
