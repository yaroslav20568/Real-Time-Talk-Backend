package interfaces

import "gin-real-time-talk/internal/entity"

type ChatUsecase interface {
	GetUserChats(userID uint, limit int, nextToken string, search string) ([]entity.Chat, string, error)
	GetChatMessages(chatID uint, userID uint, limit int, nextToken string) ([]entity.Message, string, error)
}
