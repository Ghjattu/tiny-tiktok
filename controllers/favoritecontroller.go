package controllers

import (
	"net/http"

	"github.com/Ghjattu/tiny-tiktok/models"
	"github.com/Ghjattu/tiny-tiktok/services"
	"github.com/Ghjattu/tiny-tiktok/utils"
	"github.com/gin-gonic/gin"
)

type FavoriteListResponse struct {
	Response
	VideoList []models.VideoDetail `json:"video_list"`
}

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

// Endpoint: /douyin/favorite/list/
func GetFavoriteListByUserID(c *gin.Context) {
	queryUserIDStr := c.Query("user_id")

	// Check user id is valid.
	statusCode, statusMsg, queryUserID := utils.ParseInt64(queryUserIDStr)
	if statusCode == 1 {
		c.JSON(http.StatusBadRequest, Response{
			StatusCode: statusCode,
			StatusMsg:  statusMsg,
		})
		return
	}

	// Get login user id from context.
	currentUserID := c.GetInt64("user_id")

	// Get user's favorite video list by user id.
	fs := &services.FavoriteService{}
	statusCode, statusMsg, videoList := fs.GetFavoriteVideoListByUserID(currentUserID, queryUserID)

	c.JSON(http.StatusOK, FavoriteListResponse{
		Response: Response{
			StatusCode: statusCode,
			StatusMsg:  statusMsg,
		},
		VideoList: videoList,
	})
}
