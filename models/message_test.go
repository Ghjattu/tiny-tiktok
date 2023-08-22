package models

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestCreateNewMessage(t *testing.T) {
	InitDatabase(true)

	message := &Message{
		SenderID:   1,
		ReceiverID: 2,
		Content:    "Hello World",
		CreateDate: time.Now(),
	}

	returnedMessage, err := CreateNewMessage(message)
	if err != nil {
		t.Fatalf("Error when creating new message: %v", err)
	}

	assert.Equal(t, message.SenderID, returnedMessage.SenderID)
	assert.Equal(t, message.ReceiverID, returnedMessage.ReceiverID)
	assert.Equal(t, message.Content, returnedMessage.Content)
}

func TestGetMessageList(t *testing.T) {
	InitDatabase(true)

	// Create a test message.
	CreateTestMessage(1, 2)
	time.Sleep(time.Second * 2)
	messageTwo, _ := CreateTestMessage(1, 2)

	messageList, err := GetMessageList(1, 2, messageTwo.CreateDate.Add(-time.Second))
	if err != nil {
		t.Fatalf("Error when getting message list: %v", err)
	}

	assert.Equal(t, 1, len(messageList))
}
