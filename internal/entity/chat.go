package entity

import "time"

type Chat struct {
	ID              uint      `gorm:"primaryKey" json:"id"`
	UserID          uint      `gorm:"column:user_id;not null" json:"userId"`
	User            User      `gorm:"foreignKey:UserID" json:"user"`
	LastMessageID   *uint     `gorm:"column:last_message_id" json:"lastMessageId"`
	LastMessage     *Message  `gorm:"foreignKey:LastMessageID" json:"lastMessage"`
	LastMessageText *string   `gorm:"column:last_message_text;type:text" json:"lastMessageText"`
	UnreadCount     int       `gorm:"column:unread_count;default:0" json:"unreadCount"`
	CreatedAt       time.Time `json:"createdAt"`
	UpdatedAt       time.Time `json:"updatedAt"`
}

func (Chat) TableName() string {
	return "chats"
}
