package interfaces

import "github.com/Ghjattu/tiny-tiktok/models"

type MessageInterface interface {
	// CreateNewMessage creates a new message.
	// Return status code, status message.
	CreateNewMessage(senderID, receiverID int64, content string) (int32, string)

	// GetMessageList gets message list between sender and receiver.
	// Return status code, status message, message detail list.
	GetMessageList(senderID, receiverID int64) (int32, string, []models.MessageDetail)
}
