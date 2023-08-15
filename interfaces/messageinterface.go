package interfaces

type MessageInterface interface {
	// CreateNewMessage creates a new message.
	// Return status code, status message.
	CreateNewMessage(senderID, receiverID int64, content string) (int32, string)
}
