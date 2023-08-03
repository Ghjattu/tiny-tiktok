package controllers

import (
	"net/http"
	"strconv"

	"github.com/Ghjattu/tiny-tiktok/models"
	"github.com/Ghjattu/tiny-tiktok/services"
	"github.com/gin-gonic/gin"
)

type UserResponse struct {
	Response
	User *models.User `json:"user"`
}

func GetUserByUserIDAndToken(c *gin.Context) {
	userIDString := c.Query("user_id")

	// Check if the user id is a number string.
	userID, err := strconv.ParseInt(userIDString, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, UserResponse{
			Response: Response{
				StatusCode: 1,
				StatusMsg:  "incorrect user id type",
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
