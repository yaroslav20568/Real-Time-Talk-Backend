package interfaces

import "gin-real-time-talk/internal/entity"

type MessageRepository interface {
	GetByChatID(chatID uint, limit int, page int) ([]entity.Message, int64, error)
	Create(message *entity.Message) error
}
