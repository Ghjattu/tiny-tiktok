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
	User *models.User `json:"user"`
}

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

	us := &services.UserService{}
	statusCode, statusMsg, user := us.GetUserByUserID(userID)

	c.JSON(http.StatusOK, UserResponse{
		Response: Response{
			StatusCode: statusCode,
			StatusMsg:  statusMsg,
		},
		User: user,
	})
}
