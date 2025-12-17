package interfaces

import "gin-real-time-talk/internal/entity"

type ChatRepository interface {
	GetByUserID(userID uint, limit int, nextToken string) ([]entity.Chat, string, error)
	GetByID(id uint) (*entity.Chat, error)
}
