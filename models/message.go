package models

import "time"

type Message struct {
	ID         int64     `gorm:"primaryKey;autoIncrement;not null"`
	SenderID   int64     `gorm:"type:int;not null"`
	ReceiverID int64     `gorm:"type:int;not null"`
	Content    string    `gorm:"type:varchar(255);not null"`
	CreateDate time.Time `gorm:"type:datetime;not null"`
}

type MessageDetail struct {
	ID         int64  `json:"id"`
	SenderID   int64  `json:"from_user_id"`
	ReceiverID int64  `json:"to_user_id"`
	Content    string `json:"content"`
	CreateTime int64  `json:"create_time"`
}

// CreateNewMessage create new message.
//
//	@param m *Message
//	@return *Message
//	@return error
func CreateNewMessage(m *Message) (*Message, error) {
	err := db.Model(&Message{}).Create(m).Error

	return m, err
}

// GetMessageList get message list between sender and receiver by pre message time.
//
//	@param senderID int64
//	@param receiverID int64
//	@param preMsgTime time.Time
//	@return []Message
//	@return error
func GetMessageList(senderID, receiverID int64, preMsgTime time.Time) ([]Message, error) {
	messages := make([]Message, 0)

	err := db.Model(&Message{}).
		Where("sender_id = ? AND receiver_id = ? AND create_date > ?", senderID, receiverID, preMsgTime).
		Or("sender_id = ? AND receiver_id = ? AND create_date > ?", receiverID, senderID, preMsgTime).
		Order("create_date ASC").
		Limit(20).
		Find(&messages).Error

	return messages, err
}
