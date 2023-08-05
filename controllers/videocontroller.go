package controllers

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/Ghjattu/tiny-tiktok/models"
	"github.com/Ghjattu/tiny-tiktok/services"
	"github.com/Ghjattu/tiny-tiktok/utils"
	"github.com/gin-gonic/gin"
	_ "github.com/joho/godotenv/autoload"
)

type PublishListResponse struct {
	Response
	VideoList []models.VideoDetail `json:"video_list"`
}

type FeedResponse struct {
	Response
	NextTime  int64                `json:"next_time"`
	VideoList []models.VideoDetail `json:"video_list"`
}

// Endpoint: /douyin/publish/action/
func PublishNewVideo(c *gin.Context) {
	publishTime := time.Now()
	// publishTimeUnix := publishTime.UnixNano()
	// publishTimeStr := strconv.FormatInt(publishTimeUnix, 10)

	data, err := c.FormFile("data")
	title := c.PostForm("title")

	// Failed to get video data.
	if err != nil {
		c.JSON(http.StatusOK, Response{
			StatusCode: 1,
			StatusMsg:  err.Error(),
		})
	}

	userID := c.GetInt64("user_id")

	// Save video to local.
	videoName := filepath.Base(data.Filename)
	finalVideoName := fmt.Sprintf("%s_%s", strconv.Itoa(int(userID)), videoName)
	savedPath := filepath.Join("../public/", finalVideoName)

	if err := c.SaveUploadedFile(data, savedPath); err != nil {
		c.JSON(http.StatusOK, Response{
			StatusCode: 1,
			StatusMsg:  err.Error(),
		})
	}

	// Construct play url.
	serverIP := os.Getenv("SERVER_IP")
	serverPort := os.Getenv("SERVER_PORT")
	playUrl := fmt.Sprintf("http://%s:%s/static/videos/%s", serverIP, serverPort, finalVideoName)

	// Create new video.
	vs := &services.VideoService{}
	statusCode, statusMsg := vs.CreateNewVideo(playUrl, title, userID, publishTime)

	c.JSON(http.StatusOK, Response{
		StatusCode: statusCode,
		StatusMsg:  statusMsg,
	})
}

// Endpoint: /douyin/publish/list/
func GetPublishListByAuthorID(c *gin.Context) {
	userIDString := c.Query("user_id")

	// Check user id is valid.
	statusCode, statusMsg, authorID := utils.ParseInt64(userIDString)
	if statusCode == 1 {
		c.JSON(http.StatusBadRequest, PublishListResponse{
			Response: Response{
				StatusCode: statusCode,
				StatusMsg:  statusMsg,
			},
			VideoList: nil,
		})
		return
	}

	// TODO:
	// Get current login user id.
	// currentUserID := c.GetInt64("user_id")

	// Get published video list by user id.
	vs := &services.VideoService{}
	statusCode, statusMsg, videoList := vs.GetPublishListByAuthorID(authorID)

	c.JSON(http.StatusOK, PublishListResponse{
		Response: Response{
			StatusCode: statusCode,
			StatusMsg:  statusMsg,
		},
		VideoList: videoList,
	})
}

// Endpoint: /douyin/feed/
func Feed(c *gin.Context) {
	latestTimeStr := c.Query("latest_time")

	// If the latest time is empty, set it to current time.
	if latestTimeStr == "" {
		latestTimeStr = strconv.FormatInt(time.Now().Unix(), 10)
	}

	statusCode, statusMsg, latestTimeInt := utils.ParseInt64(latestTimeStr)
	// Failed to parse latest time string to int64.
	if statusCode == 1 {
		c.JSON(http.StatusBadRequest, FeedResponse{
			Response: Response{
				StatusCode: statusCode,
				StatusMsg:  statusMsg,
			},
			NextTime:  -1,
			VideoList: nil,
		})
		return
	}

	// Convert int64 to time.Time.
	latestTime := time.Unix(latestTimeInt, 0)

	// Get most 30 videos.
	vs := &services.VideoService{}
	statusCode, statusMsg, nextTime, videoList := vs.GetMost30Videos(latestTime)

	c.JSON(http.StatusOK, FeedResponse{
		Response: Response{
			StatusCode: statusCode,
			StatusMsg:  statusMsg,
		},
		NextTime:  nextTime,
		VideoList: videoList,
	})
}
