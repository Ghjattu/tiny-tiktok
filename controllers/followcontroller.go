package controllers

import (
	"net/http"

	"github.com/Ghjattu/tiny-tiktok/services"
	"github.com/Ghjattu/tiny-tiktok/utils"
	"github.com/gin-gonic/gin"
)

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
