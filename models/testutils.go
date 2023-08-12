// Description: This testutils.go file contains some functions
// that are used in *_test.go files.

package models

import "time"

// createTestUser create a new test user.
//
//	@param name string
//	@param password string
//	@return *User
//	@return error
func CreateTestUser(name string, password string) (*User, error) {
	testUser := &User{
		Name:     name,
		Password: password,
	}

	_, err := CreateNewUser(testUser)

	return testUser, err
}

// createTestVideo create a new test video.
//
//	@param authorID int64
//	@param publishTime time.Time
//	@param title string
//	@return *Video
func CreateTestVideo(authorID int64, publishTime time.Time, title string) (*Video, error) {
	testVideo := &Video{
		AuthorID:    authorID,
		PublishTime: publishTime,
		PlayUrl:     "test",
		Title:       title,
	}

	returnedTestVideo, err := CreateNewVideo(testVideo)

	return returnedTestVideo, err
}

// createTestFavoriteRel create a new test favorite rel.
//
//	@param userID int64
//	@param videoID int64
//	@return *FavoriteRel
//	@return error
func CreateTestFavoriteRel(userID int64, videoID int64) (*FavoriteRel, error) {
	testFavoriteRel := &FavoriteRel{
		UserID:  userID,
		VideoID: videoID,
	}

	_, err := CreateNewFavoriteRel(testFavoriteRel)

	return testFavoriteRel, err
}
