package controllers

import (
	"fmt"
	"net/http"
	"path/filepath"
	"strconv"

	"github.com/Ghjattu/tiny-tiktok/models"
	"github.com/Ghjattu/tiny-tiktok/services"
	"github.com/gin-gonic/gin"
)

type VideoResponse struct {
	Response
	VideoList []models.VideoDetail `json:"video_list"`
}

// Endpoint: /douyin/publish/action/
func PublishNewVideo(c *gin.Context) {
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
	username := c.GetString("username")

	// Save video to local.
	videoName := filepath.Base(data.Filename)
	finalVideoName := fmt.Sprintf("%s_%s", username, videoName)
	savedPath := filepath.Join("../public/", finalVideoName)
	if err := c.SaveUploadedFile(data, savedPath); err != nil {
		c.JSON(http.StatusOK, Response{
			StatusCode: 1,
			StatusMsg:  err.Error(),
		})
	}

	// Create new video.
	vs := &services.VideoService{}
	statusCode, statusMsg := vs.CreateNewVideo(savedPath, title, userID, username)

	c.JSON(http.StatusOK, Response{
		StatusCode: statusCode,
		StatusMsg:  statusMsg,
	})
}

// Endpoint: /douyin/publish/list/
func GetPublishListByAuthorID(c *gin.Context) {
	userIDString := c.Query("user_id")

	// Check user id is valid.
	authorID, err := strconv.ParseInt(userIDString, 10, 64)
	if err != nil {
		if numErr, ok := err.(*strconv.NumError); ok {
			if numErr.Err == strconv.ErrSyntax {
				c.JSON(http.StatusBadRequest, VideoResponse{
					Response: Response{
						StatusCode: 1,
						StatusMsg:  "invalid syntax",
					},
					VideoList: nil,
				})
			} else if numErr.Err == strconv.ErrRange {
				c.JSON(http.StatusBadRequest, VideoResponse{
					Response: Response{
						StatusCode: 1,
						StatusMsg:  "user id out of range",
					},
					VideoList: nil,
				})
			}
		} else {
			c.JSON(http.StatusOK, VideoResponse{
				Response: Response{
					StatusCode: 1,
					StatusMsg:  "unknown error",
				},
				VideoList: nil,
			})
		}
		return
	}

	// TODO:
	// Get current login user id.
	// currentUserID := c.GetInt64("user_id")

	// Get published video list by user id.
	vs := &services.VideoService{}
	statusCode, statusMsg, videoList := vs.GetPublishListByAuthorID(authorID)

	c.JSON(http.StatusOK, VideoResponse{
		Response: Response{
			StatusCode: statusCode,
			StatusMsg:  statusMsg,
		},
		VideoList: videoList,
	})
}
