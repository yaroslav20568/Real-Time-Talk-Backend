package interfaces

import "gin-real-time-talk/internal/entity"

type ChatRepository interface {
	GetByUserID(userID uint, limit int, page int, search string) ([]entity.Chat, int64, error)
	GetByID(id uint) (*entity.Chat, error)
}
