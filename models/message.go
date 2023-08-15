package models

import "time"

type Message struct {
	ID         int64     `gorm:"primaryKey;autoIncrement;not null"`
	SenderID   int64     `gorm:"type:int;not null"`
	ReceiverID int64     `gorm:"type:int;not null"`
	Content    string    `gorm:"type:varchar(255);not null"`
	CreateDate time.Time `gorm:"type:datetime;not null"`
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
