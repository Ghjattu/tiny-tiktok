package services

import (
	"github.com/Ghjattu/tiny-tiktok/models"
	"gorm.io/gorm"
)

// FavoriteService implements FavoriteInterface.
type FavoriteService struct{}

func (fs *FavoriteService) FavoriteAction(userID, videoID, actionType int64) (int32, string) {
	// Check if the action type is valid.
	if actionType != 1 && actionType != 2 {
		return 1, "action type is invalid"
	}

	// Check if the video exist.
	_, err := models.GetVideoByID(videoID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return 1, "the video is not exist"
		}
		return 1, "get video by id failed"
	}

	// If action type is 1, create a new favorite relation.
	if actionType == 1 {
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

	// Otherwise action type is 2, delete the favorite relation.
	_, err = models.DeleteFavoriteRel(userID, videoID)
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
