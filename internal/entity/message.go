package entity

import "time"

type Message struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Text      string    `gorm:"type:text;not null" json:"text"`
	AuthorID  uint      `gorm:"column:author_id;not null" json:"authorId"`
	Author    User      `gorm:"foreignKey:AuthorID" json:"author"`
	ChatID    uint      `gorm:"column:chat_id;not null" json:"chatId"`
	Chat      *Chat     `gorm:"foreignKey:ChatID" json:"chat,omitempty"`
	IsRead    bool      `gorm:"column:is_read;default:false" json:"isRead"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

func (Message) TableName() string {
	return "messages"
}
