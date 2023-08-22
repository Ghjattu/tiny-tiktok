package services

import (
	"time"

	"github.com/Ghjattu/tiny-tiktok/models"
)

var (
	testUserOne       *models.User
	testUserOneDetail *models.UserDetail
	testUserTwo       *models.User
	// testUserTwoDetail *models.UserDetail

	testVideoOne      *models.Video
	testVideoOneCache *models.VideoCache
	testVideoTwo      *models.Video
	testVideoTwoCache *models.VideoCache

	followerUser       *models.User
	followerUserDetail *models.UserDetail

	followingUser       *models.User
	followingUserDetail *models.UserDetail

	testCommentOne      *models.Comment
	testCommentOneCache *models.CommentCache
)

func setup() {
	models.InitDatabase(true)

	// Create two test users.
	testUserOne, _ = models.CreateTestUser("testOne", "123456")
	testUserOneDetail = &models.UserDetail{ID: testUserOne.ID, Name: testUserOne.Name}
	testUserTwo, _ = models.CreateTestUser("testTwo", "123456")
	// testUserTwoDetail = &models.UserDetail{ID: testUserTwo.ID, Name: testUserTwo.Name}

	// Create two test videos.
	testVideoOne, _ = models.CreateTestVideo(testUserOne.ID, time.Now(), "testOne")
	testVideoOneCache = &models.VideoCache{ID: testVideoOne.ID, AuthorID: testVideoOne.AuthorID}
	testVideoTwo, _ = models.CreateTestVideo(testUserOne.ID, time.Now(), "testTwo")
	testVideoTwoCache = &models.VideoCache{ID: testVideoTwo.ID, AuthorID: testVideoOne.AuthorID}

	// Create a test follower user.
	followerUser, _ = models.CreateTestUser("follower", "123456")
	followerUserDetail = &models.UserDetail{ID: followerUser.ID, Name: followerUser.Name}

	// Create a test following user.
	followingUser, _ = models.CreateTestUser("following", "123456")
	followingUserDetail = &models.UserDetail{ID: followingUser.ID, Name: followingUser.Name}

	// Create a test comment.
	testCommentOne, _ = models.CreateTestComment(testUserOne.ID, testVideoOne.ID)
	testCommentOneCache = &models.CommentCache{ID: testCommentOne.ID, Content: testCommentOne.Content}
}
