package chat_usecase

import (
	"errors"

	"gin-real-time-talk/internal/entity"
	"gin-real-time-talk/internal/entity/interfaces"
	"gin-real-time-talk/pkg/pagination"
)

type chatUsecase struct {
	chatRepo    interfaces.ChatRepository
	messageRepo interfaces.MessageRepository
}

func NewChatUsecase(chatRepo interfaces.ChatRepository, messageRepo interfaces.MessageRepository) interfaces.ChatUsecase {
	return &chatUsecase{
		chatRepo:    chatRepo,
		messageRepo: messageRepo,
	}
}

func (u *chatUsecase) GetUserChats(userID uint, limit int, page int, search string) ([]entity.Chat, int, int64, error) {
	limit = pagination.NormalizeLimit(limit)

	chats, total, err := u.chatRepo.GetByUserID(userID, limit, page, search)
	if err != nil {
		return nil, 0, 0, err
	}

	totalPages := pagination.CalculateTotalPages(total, limit)

	return chats, totalPages, total, nil
}

func (u *chatUsecase) GetChatMessages(chatID uint, userID uint, limit int, page int) ([]entity.Message, int, int64, error) {
	chat, err := u.chatRepo.GetByID(chatID)
	if err != nil {
		return nil, 0, 0, errors.New("chat not found")
	}

	if chat.UserID != userID {
		return nil, 0, 0, errors.New("access denied")
	}

	limit = pagination.NormalizeLimit(limit)

	messages, total, err := u.messageRepo.GetByChatID(chatID, limit, page)
	if err != nil {
		return nil, 0, 0, err
	}

	totalPages := pagination.CalculateTotalPages(total, limit)

	return messages, totalPages, total, nil
}
