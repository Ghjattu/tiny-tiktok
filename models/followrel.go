package models

// Follower is someone who follows other users.
// Following is someone who is followed by other users.
// If user A follows user B, then user A becomes a follower of user B,
// and user B becomes a following of user A.
type FollowRel struct {
	ID          int64 `gorm:"primaryKey;autoIncrement;not null" json:"id"`
	FollowerID  int64 `gorm:"type:int;not null" json:"follower_id"`
	FollowingID int64 `gorm:"type:int;not null" json:"following_id"`
}

// CreateNewFollowRel creates a new follow relationship.
//
//	@param fr *FollowRel
//	@return *FollowRel
//	@return error
func CreateNewFollowRel(fr *FollowRel) (*FollowRel, error) {
	err := db.Model(&FollowRel{}).Create(fr).Error

	return fr, err
}

// DeleteFollowRel delete a follow relationship by follower id and following id
//
//	@param followerID int64
//	@param followingID int64
//	@return int64 "number of rows deleted"
//	@return error
func DeleteFollowRel(followerID, followingID int64) (int64, error) {
	res := db.Delete(&FollowRel{}, "follower_id = ? AND following_id = ?", followerID, followingID)

	return res.RowsAffected, res.Error
}

// CheckFollowRelExist check if a follow relationship exist by follower id and following id.
//
//	@param followerID int64
//	@param followingID int64
//	@return bool
//	@return error
func CheckFollowRelExist(followerID, followingID int64) (bool, error) {
	var count int64 = 0
	err := db.Model(&FollowRel{}).
		Where("follower_id = ? AND following_id = ?", followerID, followingID).
		Count(&count).Error

	return count > 0, err
}

// GetFollowingCountByUserID get the number of users that a user is following.
//
//	@param userID int64
//	@return int64 "number of users that a user is following"
//	@return error
func GetFollowingCountByUserID(userID int64) (int64, error) {
	var count int64 = 0
	err := db.Model(&FollowRel{}).Where("follower_id = ?", userID).Count(&count).Error

	return count, err
}

// GetFollowingListByUserID get the list of users that a user is following.
//
//	@param userID int64
//	@return []int64 "id list of users that a user is following"
//	@return error
func GetFollowingListByUserID(userID int64) ([]int64, error) {
	followingList := make([]int64, 0)

	err := db.Model(&FollowRel{}).
		Where("follower_id = ?", userID).
		Distinct().
		Pluck("following_id", &followingList).Error

	return followingList, err
}

// GetFollowerCountByUserID get the number of followers of a user.
//
//	@param userID int64
//	@return int64 "number of followers"
//	@return error
func GetFollowerCountByUserID(userID int64) (int64, error) {
	var count int64 = 0
	err := db.Model(&FollowRel{}).Where("following_id = ?", userID).Count(&count).Error

	return count, err
}

// GetFollowerListByUserID get the list of followers of a user.
//
//	@param userID int64
//	@return []int64 "id list of followers"
//	@return error
func GetFollowerListByUserID(userID int64) ([]int64, error) {
	followerList := make([]int64, 0)

	err := db.Model(&FollowRel{}).
		Where("following_id = ?", userID).
		Distinct().
		Pluck("follower_id", &followerList).Error

	return followerList, err
}
