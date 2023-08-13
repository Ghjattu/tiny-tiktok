package interfaces

import (
	"time"

	"github.com/Ghjattu/tiny-tiktok/models"
)

type CommentInterface interface {
	// CreateNewComment creates a new comment.
	// Return status_code, status_msg, comment_detail.
	CreateNewComment(currentUserID int64, videoID int64, content string, timestamp time.Time) (int32, string, *models.CommentDetail)

	// DeleteCommentByCommentID deletes a comment by its id.
	// Return status_code, status_msg, comment_detail.
	DeleteCommentByCommentID(currentUserID int64, commentID int64) (int32, string, *models.CommentDetail)
}
