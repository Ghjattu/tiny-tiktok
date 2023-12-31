package models

type FavoriteRel struct {
	ID      int64 `gorm:"primaryKey;autoIncrement"`
	UserID  int64 `gorm:"type:int;not null;index:user_id_idx;index:favorite_rel_idx"`
	VideoID int64 `gorm:"type:int;not null;index:video_id_idx;index:favorite_rel_idx"`
}

// CreateNewFavoriteRel	create a new favorite rel.
//
//	@param fr *FavoriteRel
//	@return *FavoriteRel
//	@return error
func CreateNewFavoriteRel(fr *FavoriteRel) (*FavoriteRel, error) {
	err := db.Model(&FavoriteRel{}).Create(fr).Error

	return fr, err
}

// DeleteFavoriteRel delete a favorite rel by user id and video id
// and return the number of rows deleted.
//
//	@param userID int64
//	@param videoID int64
//	@return int64 "number of rows deleted"
//	@return error
func DeleteFavoriteRel(userID, videoID int64) (int64, error) {
	res := db.Delete(&FavoriteRel{}, "user_id = ? AND video_id = ?", userID, videoID)

	return res.RowsAffected, res.Error
}

// CheckFavoriteRelExist check if a favorite rel exist by user id and video id.
//
//	@param userId int64
//	@param videoID int64
//	@return bool
//	@return error
func CheckFavoriteRelExist(userId, videoID int64) (bool, error) {
	var count int64
	err := db.Model(&FavoriteRel{}).
		Where("user_id = ? AND video_id = ?", userId, videoID).
		Count(&count).Error

	return count > 0, err
}

// GetFavoriteCountByVideoID get the count of favorite by video id.
//
//	@param videoID int64
//	@return int64
//	@return error
func GetFavoriteCountByVideoID(videoID int64) (int64, error) {
	var count int64 = 0
	err := db.Model(&FavoriteRel{}).
		Where("video_id = ?", videoID).
		Count(&count).Error

	return count, err
}

// GetFavoriteCountByUserID get the count of favorite by user id.
//
//	@param userID int64
//	@return int64
//	@return error
func GetFavoriteCountByUserID(userID int64) (int64, error) {
	var count int64 = 0
	err := db.Model(&FavoriteRel{}).
		Where("user_id = ?", userID).
		Count(&count).Error

	return count, err
}

// GetFavoriteVideoIDListByUserID get the user's favorite video id list by user id.
//
//	@param userID int64
//	@return []int64 "video id list"
//	@return error
func GetFavoriteVideoIDListByUserID(userID int64) ([]int64, error) {
	videoIDList := make([]int64, 0)
	err := db.Model(&FavoriteRel{}).
		Where("user_id = ?", userID).
		Pluck("video_id", &videoIDList).Error

	return videoIDList, err
}
