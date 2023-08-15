package services

import (
	"time"

	"github.com/Ghjattu/tiny-tiktok/models"
	"gorm.io/gorm"
)

// MessageService implements MessageInterface.
type MessageService struct{}

// CreateNewMessage creates a new message.
//
//	@receiver ms *MessageService
//	@param senderID int64
//	@param receiverID int64
//	@param content string
//	@return int32 "status code"
//	@return string "status message"
func (ms *MessageService) CreateNewMessage(senderID, receiverID int64, content string) (int32, string) {
	// Check if the content is empty.
	if content == "" {
		return 1, "message content cannot be empty"
	}

	// Check if the receiver exists.
	_, err := models.GetUserByUserID(receiverID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return 1, "receiver does not exist"
		}
		return 1, "failed to check receiver existence"
	}

	// Create a new message.
	message := &models.Message{
		SenderID:   senderID,
		ReceiverID: receiverID,
		Content:    content,
		CreateDate: time.Now(),
	}

	_, err = models.CreateNewMessage(message)
	if err != nil {
		return 1, "failed to create new message"
	}

	return 0, "create new message successfully"
}
