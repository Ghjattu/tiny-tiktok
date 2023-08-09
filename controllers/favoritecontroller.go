package controllers

import (
	"net/http"

	"github.com/Ghjattu/tiny-tiktok/services"
	"github.com/Ghjattu/tiny-tiktok/utils"
	"github.com/gin-gonic/gin"
)

// Endpoint: /douyin/favorite/action/
func FavoriteAction(c *gin.Context) {
	videoIDString := c.Query("video_id")
	actionTypeString := c.Query("action_type")

	// Check video id is valid.
	statusCode, statusMsg, videoID := utils.ParseInt64(videoIDString)
	if statusCode == 1 {
		c.JSON(http.StatusBadRequest, Response{
			StatusCode: statusCode,
			StatusMsg:  statusMsg,
		})
		return
	}

	// Check action type is valid.
	statusCode, statusMsg, actionType := utils.ParseInt64(actionTypeString)
	if statusCode == 1 {
		c.JSON(http.StatusBadRequest, Response{
			StatusCode: statusCode,
			StatusMsg:  statusMsg,
		})
		return
	}

	// Get user id from context.
	userID := c.GetInt64("user_id")

	fs := &services.FavoriteService{}
	statusCode, statusMsg = fs.FavoriteAction(userID, videoID, actionType)

	c.JSON(http.StatusOK, Response{
		StatusCode: statusCode,
		StatusMsg:  statusMsg,
	})
}
