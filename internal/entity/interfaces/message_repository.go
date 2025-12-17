package interfaces

import "gin-real-time-talk/internal/entity"

type MessageRepository interface {
	GetByChatID(chatID uint, limit int, nextToken string) ([]entity.Message, string, error)
	Create(message *entity.Message) error
}
