package models

import "time"

type Video struct {
	ID          int64     `gorm:"primaryKey;autoIncrement"`
	AuthorID    int64     `gorm:"type:int;not null"`
	PublishTime time.Time `gorm:"not null"`
	PlayUrl     string    `gorm:"type:varchar(255);not null"`
	CoverUrl    string    `gorm:"type:varchar(255);not null"`
	Title       string    `gorm:"type:varchar(255);not null"`
}

type VideoDetail struct {
	ID            int64       `json:"id" redis:"id"`
	Author        *UserDetail `json:"author"`
	PlayUrl       string      `json:"play_url" redis:"play_url"`
	CoverUrl      string      `json:"cover_url" redis:"cover_url"`
	FavoriteCount int64       `json:"favorite_count" redis:"favorite_count"`
	CommentCount  int64       `json:"comment_count" redis:"comment_count"`
	IsFavorite    bool        `json:"is_favorite"`
	Title         string      `json:"title" redis:"title"`
}

// CreateNewVideo create a new video.
//
//	@param v *Video
//	@return error
func CreateNewVideo(v *Video) (*Video, error) {
	err := db.Model(&Video{}).Create(v).Error

	return v, err
}

// TODO: should be deleted.
// GetVideoListByAuthorID get video list by user id.
//
//	@param authorID int64
//	@return []Video
//	@return error
func GetVideoListByAuthorID(authorID int64) ([]Video, error) {
	videoList := make([]Video, 0)

	err := db.Model(&Video{}).Where("author_id = ?", authorID).Find(&videoList).Error

	return videoList, err
}

// GetVideoIDListByAuthorID get video id list by author id.
//
//	@param authorID int64
//	@return []int64 "video id list"
//	@return error
func GetVideoIDListByAuthorID(authorID int64) ([]int64, error) {
	videoIDList := make([]int64, 0)

	err := db.Model(&Video{}).
		Where("author_id = ?", authorID).
		Pluck("id", &videoIDList).Error

	return videoIDList, err
}

// GetAuthorIDByVideoID get author id by video id.
//
//	@param videoID int64
//	@return int64 "author id"
//	@return error
func GetAuthorIDByVideoID(videoID int64) (int64, error) {
	var authorID int64 = 0

	err := db.Model(&Video{}).Where("id = ?", videoID).Pluck("author_id", &authorID).Error

	return authorID, err
}

// GetMost30Videos get most 30 videos earlier than latest time.
//
//	@param latestTime time.Time
//	@return []int64 "video id list"
//	@return time.Time "the earliest publish time of the video list"
//	@return error
func GetMost30Videos(latestTime time.Time) ([]int64, time.Time, error) {
	videoList := make([]Video, 0, 30)

	err := db.Model(&Video{}).
		Where("publish_time < ?", latestTime).
		Order("publish_time DESC").
		Limit(30).
		Find(&videoList).Error
	if err != nil {
		return nil, time.Time{}, err
	}

	// Set the earliest publish time to current time plus one hour by default.
	earliestTime := time.Now().Add(time.Hour * 1)

	// If the video list is not empty,
	// set the earliest publish time to the publish time of the last video.
	if len(videoList) > 0 {
		earliestTime = videoList[len(videoList)-1].PublishTime
	}

	// Get the video id list.
	videoIDList := make([]int64, 0, len(videoList))
	for _, video := range videoList {
		videoIDList = append(videoIDList, video.ID)
	}

	return videoIDList, earliestTime, err
}

// GetVideoByID get video by video id.
//
//	@param videoID int64
//	@return *Video
//	@return error
func GetVideoByID(videoID int64) (*Video, error) {
	video := &Video{}

	err := db.Model(&Video{}).Where("id = ?", videoID).First(video).Error

	return video, err
}

// GetVideoCountByAuthorID get the number of videos published by author id.
//
//	@param authorID int64
//	@return int64
//	@return error
func GetVideoCountByAuthorID(authorID int64) (int64, error) {
	var count int64 = 0

	err := db.Model(&Video{}).Where("author_id = ?", authorID).Count(&count).Error

	return count, err
}
