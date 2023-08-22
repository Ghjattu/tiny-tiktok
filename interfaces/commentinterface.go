package interfaces

import (
	"time"

	"github.com/Ghjattu/tiny-tiktok/models"
)

type CommentInterface interface {
	// CreateNewComment creates a new comment.
	// Return status_code, status_msg, comment_detail.
	CreateNewComment(currentUserID, videoID int64, content string, timestamp time.Time) (int32, string, *models.CommentDetail)

	// DeleteCommentByCommentID deletes a comment by its id.
	// Return status_code, status_msg, comment_detail.
	DeleteCommentByCommentID(currentUserID, commentID int64) (int32, string, *models.CommentDetail)

	// GetCommentListByVideoID gets a video's comment list by its id.
	// Return status_code, status_msg, comment_detail_list.
	GetCommentListByVideoID(currentUserID, videoID int64) (int32, string, []models.CommentDetail)

	// GetCommentListByCommentIDList gets a comment list by its id list.
	// Return status_code, status_msg, comment_detail_list.
	GetCommentListByCommentIDList(currentUserID int64, commentIDList []int64) (int32, string, []models.CommentDetail)

	// GetCommentDetailByCommentID gets a comment detail by its id.
	// Return comment detail, error.
	GetCommentDetailByCommentID(currentUserID, commentID int64) (*models.CommentDetail, error)
}
