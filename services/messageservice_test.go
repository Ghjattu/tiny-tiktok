package services

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

var (
	messageService = &MessageService{}
)

func TestCreateNewMessage(t *testing.T) {
	setup()

	t.Run("same sender and receiver", func(t *testing.T) {
		statusCode, _ := messageService.CreateNewMessage(1, 1, "Hello")

		assert.Equal(t, int32(1), statusCode)
	})

	t.Run("empty content", func(t *testing.T) {
		statusCode, _ := messageService.CreateNewMessage(1, 2, "")

		assert.Equal(t, int32(1), statusCode)
	})

	t.Run("receiver does not exist", func(t *testing.T) {
		statusCode, _ := messageService.CreateNewMessage(1, 0, "Hello")

		assert.Equal(t, int32(1), statusCode)
	})

	t.Run("create new message successfully", func(t *testing.T) {
		statusCode, _ :=
			messageService.CreateNewMessage(testUserOne.ID, testUserTwo.ID, "Hello")

		assert.Equal(t, int32(0), statusCode)
	})
}

func TestGetMessageList(t *testing.T) {
	setup()

	t.Run("receiver does not exist", func(t *testing.T) {
		statusCode, _, _ := messageService.GetMessageList(1, 0, time.Now())

		assert.Equal(t, int32(1), statusCode)
	})

	t.Run("get message list successfully", func(t *testing.T) {
		// Create a test message.
		messageService.CreateNewMessage(testUserOne.ID, testUserTwo.ID, "Hello")

		timestamp := time.Now().Add(-time.Hour)
		statusCode, _, messageList :=
			messageService.GetMessageList(testUserOne.ID, testUserTwo.ID, timestamp)

		assert.Equal(t, int32(0), statusCode)
		assert.Equal(t, 1, len(messageList))
	})
}
