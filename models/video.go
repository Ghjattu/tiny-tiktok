package models

type Video struct {
	ID            int64  `gorm:"primary_key;auto_increment" json:"id"`
	AuthorID      int64  `json:"author_id"`
	AuthorName    string `json:"author_name"`
	PlayUrl       string `json:"play_url"`
	CoverUrl      string `json:"cover_url"`
	FavoriteCount int64  `json:"favorite_count"`
	CommentCount  int64  `json:"comment_count"`
	IsFavorite    bool   `json:"is_favorite"`
	Title         string `json:"title"`
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
