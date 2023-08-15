package controllers

import (
	"net/http"

	"github.com/Ghjattu/tiny-tiktok/models"
	"github.com/Ghjattu/tiny-tiktok/services"
	"github.com/gin-gonic/gin"
)

type FavoriteListResponse struct {
	Response
	VideoList []models.VideoDetail `json:"video_list"`
}

// Endpoint: /douyin/favorite/action/
func FavoriteAction(c *gin.Context) {
	videoID := c.GetInt64("video_id")
	actionType := c.GetInt64("action_type")
	// Get user id from context.
	currentUserID := c.GetInt64("current_user_id")

	fs := &services.FavoriteService{}
	statusCode, statusMsg := fs.FavoriteAction(currentUserID, videoID, actionType)

	c.JSON(http.StatusOK, Response{
		StatusCode: statusCode,
		StatusMsg:  statusMsg,
	})
}

// Endpoint: /douyin/favorite/list/
func GetFavoriteListByUserID(c *gin.Context) {
	queryUserID := c.GetInt64("user_id")
	// Get login user id from context.
	currentUserID := c.GetInt64("current_user_id")

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
