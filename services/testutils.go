package services

import (
	"time"

	"github.com/Ghjattu/tiny-tiktok/bloomfilter"
	"github.com/Ghjattu/tiny-tiktok/models"
	"github.com/Ghjattu/tiny-tiktok/redis"
)

var (
	testUserOne      *models.User
	testUserOneCache *redis.UserCache
	testUserTwo      *models.User
	// testUserTwoCache *redis.UserCache

	testVideoOne      *models.Video
	testVideoOneCache *redis.VideoCache
	testVideoTwo      *models.Video
	testVideoTwoCache *redis.VideoCache

	followerUser      *models.User
	followerUserCache *redis.UserCache

	followingUser      *models.User
	followingUserCache *redis.UserCache

	testCommentOne      *models.Comment
	testCommentOneCache *redis.CommentCache

	// testMessage *models.Message
)

func setup() {
	models.InitDatabase(true)
	bloomfilter.ClearAll()

	// Create two test users.
	testUserOne, _ = models.CreateTestUser("testOne", "123456")
	testUserOneCache = &redis.UserCache{ID: testUserOne.ID, Name: testUserOne.Name}
	testUserTwo, _ = models.CreateTestUser("testTwo", "123456")
	// testUserTwoCache = &redis.UserCache{ID: testUserTwo.ID, Name: testUserTwo.Name}

	// Create two test videos.
	testVideoOne, _ = models.CreateTestVideo(testUserOne.ID, time.Now(), "testOne")
	testVideoOneCache = &redis.VideoCache{ID: testVideoOne.ID, AuthorID: testVideoOne.AuthorID}
	testVideoTwo, _ = models.CreateTestVideo(testUserOne.ID, time.Now(), "testTwo")
	testVideoTwoCache = &redis.VideoCache{ID: testVideoTwo.ID, AuthorID: testVideoOne.AuthorID}

	// Create a test follower user.
	followerUser, _ = models.CreateTestUser("follower", "123456")
	followerUserCache = &redis.UserCache{ID: followerUser.ID, Name: followerUser.Name}

	// Create a test following user.
	followingUser, _ = models.CreateTestUser("following", "123456")
	followingUserCache = &redis.UserCache{ID: followingUser.ID, Name: followingUser.Name}

	// Create a test comment.
	testCommentOne, _ = models.CreateTestComment(testUserOne.ID, testVideoOne.ID)
	testCommentOneCache = &redis.CommentCache{ID: testCommentOne.ID, Content: testCommentOne.Content}

	// Create a test message.
	models.CreateTestMessage(testUserOne.ID, testUserTwo.ID)
}

func waitForConsumer() {
	time.Sleep(100 * time.Millisecond)
}
