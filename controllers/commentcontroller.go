package controllers

import (
	"net/http"
	"time"

	"github.com/Ghjattu/tiny-tiktok/models"
	"github.com/Ghjattu/tiny-tiktok/services"
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
	videoID := c.GetInt64("video_id")
	actionType := c.GetInt64("action_type")
	currentUserID := c.GetInt64("current_user_id")

	statusCode := int32(1)
	statusMsg := "action type is invalid"
	comment := (*models.CommentDetail)(nil)

	cs := &services.CommentService{}
	if actionType == 1 {
		commentText := c.Query("comment_text")

		statusCode, statusMsg, comment =
			cs.CreateNewComment(currentUserID, videoID, commentText, time.Now())
	} else if actionType == 2 {
		commentID := c.GetInt64("comment_id")

		statusCode, statusMsg, comment = cs.DeleteCommentByCommentID(currentUserID, commentID)
	}

	c.JSON(http.StatusOK, CommentActionResponse{
		Response: Response{
			StatusCode: statusCode,
			StatusMsg:  statusMsg,
		},
		Comment: comment,
	})
}

func CommentList(c *gin.Context) {
	videoID := c.GetInt64("video_id")
	currentUserID := c.GetInt64("current_user_id")

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
