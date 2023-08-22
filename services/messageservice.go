package services

import (
	"strconv"
	"time"

	"github.com/Ghjattu/tiny-tiktok/models"
	"github.com/Ghjattu/tiny-tiktok/redis"
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
	if senderID == receiverID {
		return 1, "you can not send messages to yourself"
	}

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

	// Update the last message time between sender and receiver.
	lastMsgTimeKey := redis.LastMsgTimeKey + strconv.FormatInt(senderID, 10) + ":" +
		strconv.FormatInt(receiverID, 10)
	redis.Rdb.Set(redis.Ctx, lastMsgTimeKey, message.CreateDate.Unix(), redis.RandomDay())

	return 0, "create new message successfully"
}

// GetMessageList gets the message list between the sender and the receiver.
//
//	@receiver ms *MessageService
//	@param senderID int64
//	@param receiverID int64
//	@param preMsgTime time.Time
//	@return int32 "status code"
//	@return string "status message"
//	@return []models.MessageDetail
func (ms *MessageService) GetMessageList(senderID, receiverID int64, preMsgTime time.Time) (int32, string, []models.MessageDetail) {
	// Check if the receiver exists.
	_, err := models.GetUserByUserID(receiverID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return 1, "receiver does not exist", nil
		}
		return 1, "failed to check receiver existence", nil
	}

	// Get the message list.
	messageList, err := models.GetMessageList(senderID, receiverID, preMsgTime)
	if err != nil {
		return 1, "failed to get message list", nil
	}

	// Convert the message list to message detail list.
	messageDetailList := convertMessageToMessageDetail(messageList)

	// Update the last message time between sender and receiver.
	len := len(messageDetailList)
	if len != 0 {
		lastMsgTime := messageDetailList[len-1].CreateTime

		lastMsgTimeKey := redis.LastMsgTimeKey + strconv.FormatInt(senderID, 10) + ":" +
			strconv.FormatInt(receiverID, 10)
		redis.Rdb.Set(redis.Ctx, lastMsgTimeKey, lastMsgTime, redis.RandomDay())
	}

	return 0, "get message list successfully", messageDetailList
}

// convertMessageToMessageDetail converts the message list to message detail list.
//
//	@param messageList []models.Message
//	@return []models.MessageDetail
func convertMessageToMessageDetail(messageList []models.Message) []models.MessageDetail {
	messageDetailList := make([]models.MessageDetail, 0, len(messageList))

	for _, message := range messageList {
		messageDetail := models.MessageDetail{
			ID:         message.ID,
			SenderID:   message.SenderID,
			ReceiverID: message.ReceiverID,
			Content:    message.Content,
			CreateTime: message.CreateDate.Unix(),
		}
		messageDetailList = append(messageDetailList, messageDetail)
	}

	return messageDetailList
}
