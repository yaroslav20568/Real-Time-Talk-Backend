package interfaces

import "gin-real-time-talk/internal/entity"

type ChatUsecase interface {
	GetUserChats(userID uint, limit int, page int, search string) ([]entity.Chat, int, int64, error)
	GetChatMessages(chatID uint, userID uint, limit int, page int) ([]entity.Message, int, int64, error)
}
