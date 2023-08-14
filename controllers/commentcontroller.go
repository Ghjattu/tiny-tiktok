package controllers

import (
	"net/http"
	"time"

	"github.com/Ghjattu/tiny-tiktok/models"
	"github.com/Ghjattu/tiny-tiktok/services"
	"github.com/Ghjattu/tiny-tiktok/utils"
	"github.com/gin-gonic/gin"
)

type CommentActionResponse struct {
	Response
	Comment *models.CommentDetail `json:"comment"`
}

type CommentListResponse struct {
	Response
	CommentList []models.CommentDetail `json:"comment_list"`
}

// Endpoint: /douyin/comment/action/
func CommentAction(c *gin.Context) {
	videoIDStr := c.Query("video_id")
	actionTypeStr := c.Query("action_type")

	statusCode, statusMsg, videoID := utils.ParseInt64(videoIDStr)
	if statusCode == 1 {
		c.JSON(http.StatusBadRequest, CommentActionResponse{
			Response: Response{
				StatusCode: statusCode,
				StatusMsg:  statusMsg,
			},
			Comment: nil,
		})
		return
	}

	currentUserID := c.GetInt64("user_id")

	httpCode := http.StatusOK
	statusCode = 1
	statusMsg = "action type is invalid"
	comment := (*models.CommentDetail)(nil)

	cs := &services.CommentService{}
	if actionTypeStr == "1" {
		commentText := c.Query("comment_text")
		statusCode, statusMsg, comment =
			cs.CreateNewComment(currentUserID, videoID, commentText, time.Now())
	} else if actionTypeStr == "2" {
		commentIDStr := c.Query("comment_id")

		code, msg, commentID := utils.ParseInt64(commentIDStr)
		if code == 1 {
			httpCode = http.StatusBadRequest
			statusCode = code
			statusMsg = msg
		} else {
			statusCode, statusMsg, comment = cs.DeleteCommentByCommentID(currentUserID, commentID)
		}
	}

	c.JSON(httpCode, CommentActionResponse{
		Response: Response{
			StatusCode: statusCode,
			StatusMsg:  statusMsg,
		},
		Comment: comment,
	})
}

func CommentList(c *gin.Context) {
	videoIDStr := c.Query("video_id")

	statusCode, statusMsg, videoID := utils.ParseInt64(videoIDStr)
	if statusCode == 1 {
		c.JSON(http.StatusBadRequest, CommentListResponse{
			Response: Response{
				StatusCode: statusCode,
				StatusMsg:  statusMsg,
			},
			CommentList: nil,
		})
		return
	}

	currentUserID := c.GetInt64("user_id")

	cs := &services.CommentService{}
	statusCode, statusMsg, commentList := cs.GetCommentListByVideoID(currentUserID, videoID)

	c.JSON(http.StatusOK, CommentListResponse{
		Response: Response{
			StatusCode: statusCode,
			StatusMsg:  statusMsg,
		},
		CommentList: commentList,
	})
}
