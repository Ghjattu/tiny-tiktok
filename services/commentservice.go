package services

import (
	"strconv"
	"time"

	"github.com/Ghjattu/tiny-tiktok/models"
	"github.com/Ghjattu/tiny-tiktok/redis"
	"github.com/Ghjattu/tiny-tiktok/utils"
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
func (cs *CommentService) CreateNewComment(currentUserID, videoID int64, content string, timestamp time.Time) (int32, string, *models.CommentDetail) {
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

	// Update the CommentCount of the video in cache.
	videoKey := redis.VideoKey + strconv.FormatInt(videoID, 10)
	redis.HashIncrBy(videoKey, "comment_count", 1)

	// Insert the comment id into cache.
	commentVideoKey := redis.CommentsByVideoKey + strconv.FormatInt(videoID, 10)
	redis.Rdb.RPushX(redis.Ctx, commentVideoKey, comment.ID).Val()

	// Convert the comment to a comment detail.
	_, commentDetail := convertCommentToCommentDetail(currentUserID, comment)

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
func (cs *CommentService) DeleteCommentByCommentID(currentUserID, commentID int64) (int32, string, *models.CommentDetail) {
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
	_, commentDetail := convertCommentToCommentDetail(currentUserID, comment)

	// Update the CommentCount of the video in cache.
	videoKey := redis.VideoKey + strconv.FormatInt(comment.VideoID, 10)
	redis.HashIncrBy(videoKey, "comment_count", -1)

	// Delete the comment id from cache.
	commentVideoKey := redis.CommentsByVideoKey + strconv.FormatInt(comment.VideoID, 10)
	redis.Rdb.LRem(redis.Ctx, commentVideoKey, 0, comment.ID).Val()

	// Delete the comment.
	_, err = models.DeleteCommentByCommentID(commentID)
	if err != nil {
		return 1, "failed to delete comment", nil
	}

	return 0, "delete comment successfully", commentDetail
}

// GetCommentListByVideoID gets a video's comment list by its id.
//
//	@receiver cs *CommentService
//	@param currentUserID int64
//	@param videoID int64
//	@return int32 "status code"
//	@return string "status message"
//	@return []models.CommentDetail
func (cs *CommentService) GetCommentListByVideoID(currentUserID, videoID int64) (int32, string, []models.CommentDetail) {
	// Check if the video exist.
	_, err := models.GetVideoByID(videoID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return 1, "the video is not exist", nil
		}
		return 1, "failed to check if the video exist", nil
	}

	// Try to get comment id list from redis.
	commentVideoKey := redis.CommentsByVideoKey + strconv.FormatInt(videoID, 10)
	if redis.Rdb.Exists(redis.Ctx, commentVideoKey).Val() == 1 {
		// Cache hit.
		IDStrList, err := redis.Rdb.LRange(redis.Ctx, commentVideoKey, 0, -1).Result()
		if err == nil {
			commentIDList, _ := utils.ConvertStringToInt64(IDStrList)

			// Update the expire time.
			redis.Rdb.Expire(redis.Ctx, commentVideoKey, redis.RandomDay())

			return cs.GetCommentListByCommentIDList(currentUserID, commentIDList)
		}
	}

	// Cache miss or some error occurs.
	// Get comment id list by video id.
	commentIDList, err := models.GetCommentIDListByVideoID(videoID)
	if err != nil {
		return 0, "failed to get comment list", nil
	}

	commentIDStrList, _ := utils.ConvertInt64ToString(commentIDList)

	redis.Rdb.RPush(redis.Ctx, commentVideoKey, commentIDStrList)
	redis.Rdb.Expire(redis.Ctx, commentVideoKey, redis.RandomDay())

	return cs.GetCommentListByCommentIDList(currentUserID, commentIDList)
}

// GetCommentListByCommentIDList gets a comment list by its id list.
//
//	@receiver cs *CommentService
//	@param currentUserID int64
//	@param commentIDList []int64
//	@return int32 "status code"
//	@return string "status message"
//	@return []models.CommentDetail
func (cs *CommentService) GetCommentListByCommentIDList(currentUserID int64, commentIDList []int64) (int32, string, []models.CommentDetail) {
	commentDetailList := make([]models.CommentDetail, 0, len(commentIDList))

	for _, commentID := range commentIDList {
		// Try to get comment from redis.
		commentKey := redis.CommentKey + strconv.FormatInt(commentID, 10)
		result, err := redis.HashGetAll(commentKey)
		// Cache hit.
		if err == nil {
			commentCache := &redis.CommentCache{}
			if err := result.Scan(commentCache); err == nil {
				commentDetail := &models.CommentDetail{
					ID:         commentCache.ID,
					Content:    commentCache.Content,
					CreateDate: commentCache.CreateDate,
				}

				us := &UserService{}
				_, _, commentDetail.User =
					us.GetUserDetailByUserID(currentUserID, commentCache.UserID)

				commentDetailList = append(commentDetailList, *commentDetail)

				redis.Rdb.Expire(redis.Ctx, commentKey, redis.RandomDay())

				continue
			}
		}

		// Cache miss or some error occurs.
		commentDetail, err := cs.GetCommentDetailByCommentID(currentUserID, commentID)
		if err == nil {
			commentDetailList = append(commentDetailList, *commentDetail)
		}
	}

	return 0, "get comment list successfully", commentDetailList
}

// GetCommentDetailByCommentID gets a comment detail by its id.
//
//	@receiver cs *CommentService
//	@param currentUserID int64
//	@param commentID int64
//	@return *models.CommentDetail
//	@return error
func (cs *CommentService) GetCommentDetailByCommentID(currentUserID, commentID int64) (*models.CommentDetail, error) {
	comment, err := models.GetCommentByCommentID(commentID)
	if err != nil {
		return nil, err
	}

	// Convert the comment to a comment detail.
	commentDetail := &models.CommentDetail{
		ID:         comment.ID,
		Content:    comment.Content,
		CreateDate: comment.CreateDate.Format("01-02"),
	}

	us := &UserService{}
	_, _, commentDetail.User = us.GetUserDetailByUserID(currentUserID, comment.UserID)

	// Insert the comment into redis.
	commentKey := redis.CommentKey + strconv.FormatInt(commentID, 10)
	commentCache := &redis.CommentCache{
		ID:         comment.ID,
		UserID:     comment.UserID,
		Content:    comment.Content,
		CreateDate: comment.CreateDate.Format("01-02"),
	}
	redis.Rdb.HSet(redis.Ctx, commentKey, commentCache)
	redis.Rdb.Expire(redis.Ctx, commentKey, redis.RandomDay())

	return commentDetail, nil
}

// convertCommentToCommentDetail converts a comment to a comment detail.
//
//	@param currentUserID int64
//	@param comment *models.Comment
//	@return int32 "status code"
//	@return *models.CommentDetail
func convertCommentToCommentDetail(currentUserID int64, comment *models.Comment) (int32, *models.CommentDetail) {
	us := &UserService{}
	commentDetail := &models.CommentDetail{
		ID:         comment.ID,
		Content:    comment.Content,
		CreateDate: comment.CreateDate.Format("01-02"),
	}

	// Get the user detail by user id.
	_, _, user := us.GetUserDetailByUserID(currentUserID, comment.UserID)
	commentDetail.User = user

	return 0, commentDetail
}
