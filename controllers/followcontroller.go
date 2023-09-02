package controllers

import (
	"net/http"

	"github.com/Ghjattu/tiny-tiktok/models"
	"github.com/Ghjattu/tiny-tiktok/services"
	"github.com/gin-gonic/gin"
)

type UserListResponse struct {
	Response
	UserList []models.UserDetail `json:"user_list"`
}

func FollowAction(c *gin.Context) {
	followingID := c.GetInt64("to_user_id")
	actionType := c.GetInt64("action_type")
	currentUserID := c.GetInt64("current_user_id")

	statusCode := int32(1)
	statusMsg := "action type is invalid"

	fs := &services.FollowService{}
	if actionType == 1 {
		statusCode, statusMsg = fs.CreateNewFollowRel(currentUserID, followingID)
	} else if actionType == 2 {
		statusCode, statusMsg = fs.DeleteFollowRel(currentUserID, followingID)
	}

	c.JSON(http.StatusOK, Response{
		StatusCode: statusCode,
		StatusMsg:  statusMsg,
	})
}

func FollowingList(c *gin.Context) {
	queryUserID := c.GetInt64("user_id")
	currentUserID := c.GetInt64("current_user_id")

	fs := &services.FollowService{}
	statusCode, statusMsg, userList := fs.GetFollowingListByUserID(currentUserID, queryUserID)

	c.JSON(http.StatusOK, UserListResponse{
		Response: Response{
			StatusCode: statusCode,
			StatusMsg:  statusMsg,
		},
		UserList: userList,
	})
}

func FollowerList(c *gin.Context) {
	queryUserID := c.GetInt64("user_id")
	currentUserID := c.GetInt64("current_user_id")

	fs := &services.FollowService{}
	statusCode, statusMsg, userList := fs.GetFollowerListByUserID(currentUserID, queryUserID)

	c.JSON(http.StatusOK, UserListResponse{
		Response: Response{
			StatusCode: statusCode,
			StatusMsg:  statusMsg,
		},
		UserList: userList,
	})
}
