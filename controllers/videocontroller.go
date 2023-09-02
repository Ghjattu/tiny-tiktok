package controllers

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/Ghjattu/tiny-tiktok/ffmpeg"
	"github.com/Ghjattu/tiny-tiktok/models"
	"github.com/Ghjattu/tiny-tiktok/oss"
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
	publishTimeStr := publishTime.Format("2006-01-02-15:04:05")
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
		return
	}

	// Check title is empty or not.
	if title == "" {
		c.JSON(http.StatusOK, Response{
			StatusCode: 1,
			StatusMsg:  "video title is empty",
		})
		return
	}

	currentUserID := c.GetInt64("current_user_id")
	currentUserIDStr := fmt.Sprintf("%d", currentUserID)

	// Save video to local.
	videoName := filepath.Base(data.Filename)
	finalVideoName := fmt.Sprintf("%s_%s_%s", currentUserIDStr, publishTimeStr, videoName)
	videoSavedLocalPath := filepath.Join("../public/", finalVideoName)

	if err := c.SaveUploadedFile(data, videoSavedLocalPath); err != nil {
		c.JSON(http.StatusOK, Response{
			StatusCode: 1,
			StatusMsg:  err.Error(),
		})
		return
	}
	defer os.Remove(videoSavedLocalPath)

	// Upload video to OSS.
	if err := oss.UploadFile(finalVideoName, videoSavedLocalPath); err != nil {
		c.JSON(http.StatusOK, Response{
			StatusCode: 1,
			StatusMsg:  err.Error(),
		})
		return
	}

	// Construct play url.
	OSSEndpoint := os.Getenv("OSS_ENDPOINT")
	OSSBucketName := os.Getenv("OSS_BUCKET_NAME")
	playUrl := fmt.Sprintf("https://%s.%s/%s", OSSBucketName, OSSEndpoint, finalVideoName)

	// Get snapshot from video.
	finalCoverName := fmt.Sprintf("%s_%s_%s_cover.jpeg",
		currentUserIDStr, publishTimeStr, strings.Split(videoName, ".")[0])
	coverSavedLocalPath := filepath.Join("../public/", finalCoverName)
	coverUrl := ""

	err = ffmpeg.GetSnapshot(videoSavedLocalPath, coverSavedLocalPath, 1)
	if err == nil {
		// Upload cover to OSS.
		if err := oss.UploadFile(finalCoverName, coverSavedLocalPath); err == nil {
			coverUrl = fmt.Sprintf("https://%s.%s/%s", OSSBucketName, OSSEndpoint, finalCoverName)
		}
	}
	defer os.Remove(coverSavedLocalPath)

	// If failed to get snapshot, or failed to upload cover to OSS, use default cover.
	if coverUrl == "" {
		coverUrl = fmt.Sprintf("https://%s.%s/%s", OSSBucketName, OSSEndpoint, "default_cover.png")
	}

	// Create new video.
	vs := &services.VideoService{}
	statusCode, statusMsg := vs.CreateNewVideo(playUrl, coverUrl, title, currentUserID, publishTime)

	c.JSON(http.StatusOK, Response{
		StatusCode: statusCode,
		StatusMsg:  statusMsg,
	})
}

// Endpoint: /douyin/publish/list/
func GetPublishListByAuthorID(c *gin.Context) {
	authorID := c.GetInt64("user_id")
	// Get current login user id.
	currentUserID := c.GetInt64("current_user_id")

	// Get published video list by author id.
	vs := &services.VideoService{}
	statusCode, statusMsg, videoList := vs.GetVideoListByAuthorID(authorID, currentUserID)

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
	currentUserID := c.GetInt64("current_user_id")

	// If the latest time is empty, set it to current time.
	if latestTimeStr == "" || latestTimeStr == "0" {
		latestTimeStr = strconv.FormatInt(time.Now().Unix(), 10)
	}

	statusCode, statusMsg, latestTimeInt := utils.ParseInt64(latestTimeStr)
	// Failed to parse latest time string to int64.
	if statusCode == 1 {
		c.JSON(http.StatusOK, FeedResponse{
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
	statusCode, statusMsg, nextTime, videoList := vs.GetMost30Videos(latestTime, currentUserID)

	c.JSON(http.StatusOK, FeedResponse{
		Response: Response{
			StatusCode: statusCode,
			StatusMsg:  statusMsg,
		},
		NextTime:  nextTime,
		VideoList: videoList,
	})
}
