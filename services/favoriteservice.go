package services

import (
	"github.com/Ghjattu/tiny-tiktok/models"
	"gorm.io/gorm"
)

// FavoriteService implements FavoriteInterface.
type FavoriteService struct{}

func (fs *FavoriteService) CreateNewFavoriteRel(userID, videoID int64) (int32, string) {
	// Check if the video exist.
	_, err := models.GetVideoByID(videoID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return 1, "the video is not exist"
		}
		return 1, "get video by id failed"
	}

	// Check if the favorite relation exist.
	exist, err := models.CheckFavoriteRelExist(userID, videoID)
	if err != nil {
		return 1, "check favorite rel exist failed"
	}

	if exist {
		return 1, "you have already favorited this video"
	}

	// Create a new favorite relation.
	fr := &models.FavoriteRel{
		UserID:  userID,
		VideoID: videoID,
	}

	_, err = models.CreateNewFavoriteRel(fr)
	if err != nil {
		return 1, "favorite action failed"
	}

	return 0, "favorite action success"
}

func (fs *FavoriteService) DeleteFavoriteRel(userID, videoID int64) (int32, string) {
	_, err := models.DeleteFavoriteRel(userID, videoID)
	if err != nil {
		return 1, "unfavorite action failed"
	}

	return 0, "unfavorite action success"
}

func (fs *FavoriteService) GetFavoriteVideoListByUserID(currentUserID, queryUserID int64) (int32, string, []models.VideoDetail) {
	// Get favorite video id list by user id.
	favoriteVideoIDList, err := models.GetFavoriteVideoIDListByUserID(queryUserID)
	if err != nil {
		return 1, "failed to get favorite video id list", nil
	}

	// Get favorite video list by video id list.
	vs := &VideoService{}
	statusCode, _, favoriteVideoList := vs.GetVideoListByVideoIDList(favoriteVideoIDList, currentUserID)
	if statusCode == 1 {
		return 1, "failed to get favorite video list", nil
	}

	return 0, "get favorite video list successfully", favoriteVideoList
}

// GetTotalFavoritedByUserID returns the total number of favorited by user id.
//
//	@receiver fs FavoriteService
//	@param userID int64
//	@return int64
func (fs *FavoriteService) GetTotalFavoritedByUserID(userID int64) int64 {
	totalFavorited := int64(0)

	videoList, _ := models.GetVideoListByAuthorID(userID)
	for _, video := range videoList {
		count, _ := models.GetFavoriteCountByVideoID(video.ID)
		totalFavorited += count
	}

	return totalFavorited
}
