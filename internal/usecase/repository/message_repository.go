package repository

import (
	"gin-real-time-talk/internal/entity"
	"gin-real-time-talk/internal/entity/interfaces"
	"gin-real-time-talk/pkg/pagination"

	"gorm.io/gorm"
)

type messageRepository struct {
	db *gorm.DB
}

func NewMessageRepository(db *gorm.DB) interfaces.MessageRepository {
	return &messageRepository{
		db: db,
	}
}

func (r *messageRepository) GetByChatID(chatID uint, limit int, nextToken string) ([]entity.Message, string, error) {
	limit = pagination.NormalizeLimit(limit)

	query := r.db.Where(&entity.Message{ChatID: chatID}).
		Preload("Author").
		Preload("Chat.User").
		Order("messages.created_at DESC, messages.id DESC")

	if nextToken != "" {
		cursorID, err := pagination.DecodeToken(nextToken)
		if err == nil && cursorID > 0 {
			var cursorMessage entity.Message
			if err := r.db.First(&cursorMessage, cursorID).Error; err == nil {
				query = query.Where("(messages.created_at, messages.id) < (?, ?)", cursorMessage.CreatedAt, cursorMessage.ID)
			}
		}
	}

	query = query.Limit(limit + 1)

	var messages []entity.Message
	if err := query.Find(&messages).Error; err != nil {
		return nil, "", err
	}

	var hasNext bool
	if len(messages) > limit {
		hasNext = true
		messages = messages[:limit]
	}

	var token string
	if hasNext && len(messages) > 0 {
		lastMessage := messages[len(messages)-1]
		token = pagination.EncodeToken(lastMessage.ID)
	}

	return messages, token, nil
}

func (r *messageRepository) GetByID(id uint) (*entity.Message, error) {
	var message entity.Message
	err := r.db.Preload("Author").Preload("Chat.User").First(&message, id).Error
	if err != nil {
		return nil, err
	}
	return &message, nil
}

func (r *messageRepository) Create(message *entity.Message) error {
	return r.db.Create(message).Error
}
