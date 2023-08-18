package services

import (
	"testing"

	"github.com/Ghjattu/tiny-tiktok/models"
	"github.com/stretchr/testify/assert"
)

var (
	messageService = &MessageService{}
)

func TestCreateNewMessageWithSameSenderAndReceiver(t *testing.T) {
	models.InitDatabase(true)

	statusCode, statusMsg := messageService.CreateNewMessage(1, 1, "Hello")

	assert.Equal(t, int32(1), statusCode)
	assert.Equal(t, "you can not send messages to yourself", statusMsg)
}

func TestCreateNewMessageWithEmptyContent(t *testing.T) {
	models.InitDatabase(true)

	// Create a test user.
	testUser, _ := models.CreateTestUser("test", "123456")

	statusCode, statusMsg := messageService.CreateNewMessage(testUser.ID+1, testUser.ID, "")

	assert.Equal(t, int32(1), statusCode)
	assert.Equal(t, "message content cannot be empty", statusMsg)
}

func TestCreateNewMessageWithNonExistUser(t *testing.T) {
	models.InitDatabase(true)

	statusCode, statusMsg := messageService.CreateNewMessage(1, 2, "Hello")

	assert.Equal(t, int32(1), statusCode)
	assert.Equal(t, "receiver does not exist", statusMsg)
}

func TestCreateNewMessage(t *testing.T) {
	models.InitDatabase(true)

	// Create a test user.
	testUser, _ := models.CreateTestUser("test", "123456")

	statusCode, statusMsg := messageService.CreateNewMessage(testUser.ID+1, testUser.ID, "Hello")

	assert.Equal(t, int32(0), statusCode)
	assert.Equal(t, "create new message successfully", statusMsg)
}

func TestGetMessageListWithNonExistUser(t *testing.T) {
	models.InitDatabase(true)

	statusCode, statusMsg, _ := messageService.GetMessageList(1, 2)

	assert.Equal(t, int32(1), statusCode)
	assert.Equal(t, "receiver does not exist", statusMsg)
}

func TestGetMessageList(t *testing.T) {
	models.InitDatabase(true)

	// Create two test users.
	testUserOne, _ := models.CreateTestUser("testOne", "123456")
	testUserTwo, _ := models.CreateTestUser("testTwo", "123456")
	// Create a test message.
	models.CreateTestMessage(testUserOne.ID, testUserTwo.ID)

	statusCode, statusMsg, messageList :=
		messageService.GetMessageList(testUserOne.ID, testUserTwo.ID)

	assert.Equal(t, int32(0), statusCode)
	assert.Equal(t, "get message list successfully", statusMsg)
	assert.Equal(t, 1, len(messageList))

}
