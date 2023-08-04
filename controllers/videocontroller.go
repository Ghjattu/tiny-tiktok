package controllers

import (
	"fmt"
	"net/http"
	"path/filepath"

	"github.com/Ghjattu/tiny-tiktok/services"
	"github.com/gin-gonic/gin"
)

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
	savedPath := filepath.Join("./public/", finalVideoName)
	if err := c.SaveUploadedFile(data, savedPath); err != nil {
		c.JSON(http.StatusOK, Response{
			StatusCode: 1,
			StatusMsg:  err.Error(),
		})
	}

	vs := &services.VideoService{}
	statusCode, statusMsg := vs.CreateNewVideo(savedPath, title, userID, username)

	c.JSON(http.StatusOK, Response{
		StatusCode: statusCode,
		StatusMsg:  statusMsg,
	})
}
