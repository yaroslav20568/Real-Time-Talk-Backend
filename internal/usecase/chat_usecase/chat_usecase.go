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

func (u *chatUsecase) GetUserChats(userID uint, limit int, nextToken string, search string) ([]entity.Chat, string, error) {
	limit = pagination.NormalizeLimit(limit)

	chats, token, err := u.chatRepo.GetByUserID(userID, limit, nextToken, search)
	if err != nil {
		return nil, "", err
	}

	return chats, token, nil
}

func (u *chatUsecase) GetChatMessages(chatID uint, userID uint, limit int, nextToken string) ([]entity.Message, string, error) {
	chat, err := u.chatRepo.GetByID(chatID)
	if err != nil {
		return nil, "", errors.New("chat not found")
	}

	if chat.UserID != userID {
		return nil, "", errors.New("access denied")
	}

	limit = pagination.NormalizeLimit(limit)

	messages, token, err := u.messageRepo.GetByChatID(chatID, limit, nextToken)
	if err != nil {
		return nil, "", err
	}

	return messages, token, nil
}

func (u *chatUsecase) CreateMessage(senderID uint, recipientID uint, text string) (*entity.Message, error) {
	if text == "" {
		return nil, errors.New("message text cannot be empty")
	}

	chat, err := u.chatRepo.FindOrCreateChatByUsers(senderID, recipientID)
	if err != nil {
		return nil, err
	}

	message := &entity.Message{
		Text:     text,
		AuthorID: senderID,
		ChatID:   chat.ID,
		IsRead:   false,
	}

	if err := u.messageRepo.Create(message); err != nil {
		return nil, err
	}

	chat.LastMessageID = &message.ID
	chat.LastMessageText = &message.Text

	if err := u.chatRepo.Update(chat); err != nil {
		return nil, err
	}

	messageWithRelations, err := u.messageRepo.GetByID(message.ID)
	if err != nil {
		return nil, err
	}

	return messageWithRelations, nil
}
