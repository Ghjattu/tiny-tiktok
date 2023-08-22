package models

import "time"

type Comment struct {
	ID         int64     `gorm:"primaryKey;autoIncrement"`
	UserID     int64     `gorm:"type:int;not null"`
	VideoID    int64     `gorm:"type:int;not null"`
	Content    string    `gorm:"type:varchar(255);not null"`
	CreateDate time.Time `gorm:"not null"`
}

type CommentDetail struct {
	ID         int64       `json:"id"`
	User       *UserDetail `json:"user"`
	Content    string      `json:"content"`
	CreateDate string      `json:"create_date"`
}

type CommentCache struct {
	ID         int64  `redis:"id"`
	UserID     int64  `redis:"user_id"`
	Content    string `redis:"content"`
	CreateDate string `redis:"create_date"`
}

// CreateNewComment creates a new comment.
//
//	@param c *Comment
//	@return *Comment
//	@return error
func CreateNewComment(c *Comment) (*Comment, error) {
	err := db.Model(&Comment{}).Create(c).Error

	return c, err
}

// GetCommentCountByVideoID returns the number of comments for a video.
//
//	@param videoID int64
//	@return int64 "the number of comments"
//	@return error
func GetCommentCountByVideoID(videoID int64) (int64, error) {
	var count int64 = 0

	err := db.Model(&Comment{}).Where("video_id = ?", videoID).Count(&count).Error

	return count, err
}

// DeleteCommentByCommentID deletes a comment by its id.
//
//	@param commentID int64
//	@return int64 "the number of rows deleted"
//	@return error
func DeleteCommentByCommentID(commentID int64) (int64, error) {
	res := db.Delete(&Comment{}, "id = ?", commentID)

	return res.RowsAffected, res.Error
}

// GetCommentByCommentID gets a comment by its id.
//
//	@param commentID int64
//	@return *Comment
//	@return error
func GetCommentByCommentID(commentID int64) (*Comment, error) {
	comment := &Comment{}

	err := db.Model(&Comment{}).Where("id = ?", commentID).First(comment).Error

	return comment, err
}

// GetCommentIDListByVideoID returns a video's comment id list by its id.
//
//	@param videoID int64
//	@return []int64 "comment id list"
//	@return error
func GetCommentIDListByVideoID(videoID int64) ([]int64, error) {
	commentIDList := make([]int64, 0)

	err := db.Model(&Comment{}).
		Where("video_id = ?", videoID).
		Pluck("id", &commentIDList).Error

	return commentIDList, err
}
