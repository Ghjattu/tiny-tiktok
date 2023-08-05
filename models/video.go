package models

import "time"

type Video struct {
	ID            int64     `gorm:"primary_key;auto_increment" json:"id"`
	AuthorID      int64     `json:"author_id"`
	PublishTime   time.Time `json:"publish_time"`
	PlayUrl       string    `json:"play_url"`
	CoverUrl      string    `json:"cover_url"`
	FavoriteCount int64     `json:"favorite_count"`
	CommentCount  int64     `json:"comment_count"`
	IsFavorite    bool      `json:"is_favorite"`
	Title         string    `json:"title"`
}

type VideoDetail struct {
	ID            int64  `json:"id"`
	Author        *User  `json:"author"`
	PlayUrl       string `json:"play_url"`
	CoverUrl      string `json:"cover_url"`
	FavoriteCount int64  `json:"favorite_count"`
	CommentCount  int64  `json:"comment_count"`
	IsFavorite    bool   `json:"is_favorite"`
	Title         string `json:"title"`
}

// CreateNewVideo create a new video.
//
//	@param v *Video
//	@return error
func CreateNewVideo(v *Video) (*Video, error) {
	err := db.Model(&Video{}).Create(v).Error

	return v, err
}

// GetVideoListByUserID get video list by user id.
//
//	@param userID int64
//	@return []VideoDetail
//	@return error
func GetVideoListByUserID(userID int64) ([]VideoDetail, error) {
	// Get user by user id.
	user := &User{}

	err := db.Model(&User{}).Where("id = ?", userID).First(user).Error
	if err != nil {
		return nil, err
	}

	// Hide user password.
	user.Password = ""

	// Get temporary video list by user id.
	tempVideoList := make([]Video, 0)

	err = db.Model(&Video{}).Where("author_id = ?", userID).Find(&tempVideoList).Error
	if err != nil {
		return nil, err
	}

	// Convert temporary video list to video list.
	videoList := make([]VideoDetail, 0, len(tempVideoList))

	for _, video := range tempVideoList {
		videoList = append(videoList, VideoDetail{
			ID:            video.ID,
			Author:        user,
			PlayUrl:       video.PlayUrl,
			CoverUrl:      video.CoverUrl,
			FavoriteCount: video.FavoriteCount,
			CommentCount:  video.CommentCount,
			IsFavorite:    video.IsFavorite,
			Title:         video.Title,
		})
	}

	return videoList, err
}

// GetMost30Videos get most 30 videos earlier than latest time.
//
//	@param latestTime time.Time
//	@return []VideoDetail
//	@return time.Time "the earliest publish time of the video list"
//	@return error
func GetMost30Videos(latestTime time.Time) ([]VideoDetail, time.Time, error) {
	// Get temporary video list by latest time.
	tempVideoList := make([]Video, 0)

	err := db.Model(&Video{}).
		Where("publish_time < ?", latestTime).
		Order("publish_time DESC").
		Limit(30).
		Find(&tempVideoList).Error
	if err != nil {
		return nil, time.Time{}, err
	}

	// Convert temporary video list to video list.
	videoList := make([]VideoDetail, 0, len(tempVideoList))

	for _, video := range tempVideoList {
		// Get user by user id.
		user := &User{}

		err := db.Model(&User{}).Where("id = ?", video.AuthorID).First(user).Error
		if err != nil {
			return nil, time.Time{}, err
		}

		// Hide user password.
		user.Password = ""

		videoList = append(videoList, VideoDetail{
			ID:            video.ID,
			Author:        user,
			PlayUrl:       video.PlayUrl,
			CoverUrl:      video.CoverUrl,
			FavoriteCount: video.FavoriteCount,
			CommentCount:  video.CommentCount,
			IsFavorite:    video.IsFavorite,
			Title:         video.Title,
		})
	}

	// Set the earliest publish time to current time plus one hour by default.
	earliestTime := time.Now().Add(time.Hour * 1)

	// If the video list is not empty,
	// set the earliest publish time to the publish time of the last video.
	if len(tempVideoList) > 0 {
		earliestTime = tempVideoList[len(tempVideoList)-1].PublishTime
	}

	return videoList, earliestTime, err
}
