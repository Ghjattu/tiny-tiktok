package controllers

import (
	"net/http"

	"github.com/Ghjattu/tiny-tiktok/models"
	"github.com/Ghjattu/tiny-tiktok/services"
	"github.com/gin-gonic/gin"
)

type FriendListResponse struct {
	Response
	UserList []models.UserDetail `json:"user_list"`
}

func FriendList(c *gin.Context) {
	queryUserID := c.GetInt64("user_id")
	currentUserID := c.GetInt64("current_user_id")

	fs := &services.FriendService{}
	statusCode, statusMsg, friendList :=
		fs.GetFriendListByUserID(currentUserID, queryUserID)

	c.JSON(http.StatusOK, FriendListResponse{
		Response: Response{
			StatusCode: statusCode,
			StatusMsg:  statusMsg,
		},
		UserList: friendList,
	})
}
