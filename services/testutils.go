// Description: This testutils.go file contains some functions that are used in test files
// in the services package.

package services

import (
	"time"

	"github.com/Ghjattu/tiny-tiktok/models"
)

// createTestUser creates a new test user.
//
//	@param name string
//	@param password string
//	@return *models.User
//	@return error
func createTestUser(name string, password string) (*models.User, error) {
	testUser := &models.User{
		Name:     name,
		Password: password,
	}

	_, err := models.CreateNewUser(testUser)

	return testUser, err
}

// createTestVideo create a new test video.
//
//	@param authorID int64
//	@param publishTime time.Time
//	@param title string
//	@return *Video
func createTestVideo(authorID int64, publishTime time.Time, title string) (*models.Video, error) {
	testVideo := &models.Video{
		AuthorID:    authorID,
		PublishTime: publishTime,
		PlayUrl:     "test",
		Title:       title,
	}

	_, err := models.CreateNewVideo(testVideo)

	return testVideo, err
}
