package services

import (
	"time"

	"github.com/Ghjattu/tiny-tiktok/models"
	"gorm.io/gorm"
)

// CommentService implements CommentInterface.
type CommentService struct{}

// CreateNewComment creates a new comment.
//
//	@receiver cs *CommentService
//	@param currentUserID int64
//	@param videoID int64
//	@param content string
//	@param timestamp time.Time
//	@return int32 "status code"
//	@return string "status message"
//	@return *models.CommentDetail
func (cs *CommentService) CreateNewComment(currentUserID int64, videoID int64, content string, timestamp time.Time) (int32, string, *models.CommentDetail) {
	// Check if the content is empty.
	if content == "" {
		return 1, "comment text cannot be empty", nil
	}

	// Check if the video exist.
	_, err := models.GetVideoByID(videoID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return 1, "the video is not exist", nil
		}
		return 1, "failed to check if the video exist", nil
	}

	// Create a new comment.
	comment := &models.Comment{
		UserID:     currentUserID,
		VideoID:    videoID,
		Content:    content,
		CreateDate: timestamp,
	}

	_, err = models.CreateNewComment(comment)
	if err != nil {
		return 1, "failed to create new comment", nil
	}

	// Convert the comment to a comment detail.
	_, commentDetail := convertCommentToCommentDetail(comment)

	return 0, "create new comment successfully", commentDetail
}

// DeleteCommentByCommentID deletes a comment by its id.
//
//	@receiver cs *CommentService
//	@param currentUserID int64
//	@param commentID int64
//	@return int32 "status code"
//	@return string "status message"
//	@return *models.CommentDetail
func (cs *CommentService) DeleteCommentByCommentID(currentUserID int64, commentID int64) (int32, string, *models.CommentDetail) {
	// Check if the comment exist.
	comment, err := models.GetCommentByCommentID(commentID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return 1, "the comment is not exist", nil
		}
		return 1, "failed to check if the comment exist", nil
	}

	// Check if the current user have permission to delete this comment.
	video, _ := models.GetVideoByID(comment.VideoID)
	if currentUserID != comment.UserID && currentUserID != video.AuthorID {
		return 1, "you do not have permission to delete this comment", nil
	}

	// Convert the comment to a comment detail.
	_, commentDetail := convertCommentToCommentDetail(comment)

	// Delete the comment.
	_, err = models.DeleteCommentByCommentID(commentID)
	if err != nil {
		return 1, "failed to delete comment", nil
	}

	return 0, "delete comment successfully", commentDetail
}

// convertCommentToCommentDetail converts a comment to a comment detail.
//
//	@param comment *models.Comment
//	@return int32 "status code"
//	@return *models.CommentDetail
func convertCommentToCommentDetail(comment *models.Comment) (int32, *models.CommentDetail) {
	us := &UserService{}
	commentDetail := &models.CommentDetail{
		ID:         comment.ID,
		Content:    comment.Content,
		CreateDate: comment.CreateDate.Format("01-02"),
	}

	// Get the user detail by user id.
	_, _, user := us.GetUserDetailByUserID(comment.UserID)
	commentDetail.User = user

	return 0, commentDetail
}
