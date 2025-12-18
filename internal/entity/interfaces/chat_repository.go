package interfaces

import "gin-real-time-talk/internal/entity"

type ChatRepository interface {
	GetByUserID(userID uint, limit int, page int, search string) ([]entity.Chat, int64, error)
	GetByID(id uint) (*entity.Chat, error)
	FindOrCreateChatByUsers(senderID uint, recipientID uint) (*entity.Chat, error)
	Create(chat *entity.Chat) error
	Update(chat *entity.Chat) error
}
