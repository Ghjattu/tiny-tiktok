package controllers

import (
	"net/http"

	"github.com/Ghjattu/tiny-tiktok/models"
	"github.com/Ghjattu/tiny-tiktok/services"
	"github.com/Ghjattu/tiny-tiktok/utils"
	"github.com/gin-gonic/gin"
)

type UserListResponse struct {
	Response
	UserList []models.UserDetail `json:"user_list"`
}

func FollowAction(c *gin.Context) {
	followingIDStr := c.Query("to_user_id")
	actionTypeStr := c.Query("action_type")

	// Check if the user id is valid.
	statusCode, statusMsg, followingID := utils.ParseInt64(followingIDStr)
	if statusCode == 1 {
		c.JSON(http.StatusBadRequest, Response{
			StatusCode: statusCode,
			StatusMsg:  statusMsg,
		})
		return
	}

	currentUserID := c.GetInt64("user_id")

	statusCode = 1
	statusMsg = "action type is invalid"

	fs := &services.FollowService{}
	if actionTypeStr == "1" {
		statusCode, statusMsg = fs.CreateNewFollowRel(currentUserID, followingID)
	} else if actionTypeStr == "2" {
		statusCode, statusMsg = fs.DeleteFollowRel(currentUserID, followingID)
	}

	c.JSON(http.StatusOK, Response{
		StatusCode: statusCode,
		StatusMsg:  statusMsg,
	})
}

func FollowingList(c *gin.Context) {
	userIDStr := c.Query("user_id")

	// Check if the user id is valid.
	statusCode, statusMsg, userID := utils.ParseInt64(userIDStr)
	if statusCode == 1 {
		c.JSON(http.StatusBadRequest, UserListResponse{
			Response: Response{
				StatusCode: statusCode,
				StatusMsg:  statusMsg,
			},
			UserList: nil,
		})
		return
	}

	currentUserID := c.GetInt64("user_id")

	fs := &services.FollowService{}
	statusCode, statusMsg, userList := fs.GetFollowingListByUserID(currentUserID, userID)

	c.JSON(http.StatusOK, UserListResponse{
		Response: Response{
			StatusCode: statusCode,
			StatusMsg:  statusMsg,
		},
		UserList: userList,
	})
}

func FollowerList(c *gin.Context) {
	userIDStr := c.Query("user_id")

	// Check if the user id is valid.
	statusCode, statusMsg, userID := utils.ParseInt64(userIDStr)
	if statusCode == 1 {
		c.JSON(http.StatusBadRequest, UserListResponse{
			Response: Response{
				StatusCode: statusCode,
				StatusMsg:  statusMsg,
			},
			UserList: nil,
		})
		return
	}

	currentUserID := c.GetInt64("user_id")

	fs := &services.FollowService{}
	statusCode, statusMsg, userList := fs.GetFollowerListByUserID(currentUserID, userID)

	c.JSON(http.StatusOK, UserListResponse{
		Response: Response{
			StatusCode: statusCode,
			StatusMsg:  statusMsg,
		},
		UserList: userList,
	})
}
