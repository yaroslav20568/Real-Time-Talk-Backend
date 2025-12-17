package repository

import (
	"time"

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

	query := r.db.Where("chat_id = ?", chatID).
		Preload("Author").
		Preload("Chat").
		Order("created_at DESC, id DESC")

	if nextToken != "" {
		tokenData, err := pagination.ParseToken(nextToken)
		if err != nil {
			return nil, "", err
		}

		if tokenData != nil {
			if tokenData.Timestamp > 0 {
				lastCreatedAt := time.Unix(tokenData.Timestamp, 0)
				query = query.Where("(created_at < ? OR (created_at = ? AND id < ?))", lastCreatedAt, lastCreatedAt, tokenData.ID)
			} else {
				query = query.Where("id < ?", tokenData.ID)
			}
		}
	}

	var messages []entity.Message
	if err := query.Limit(limit + 1).Find(&messages).Error; err != nil {
		return nil, "", err
	}

	var newNextToken string
	if len(messages) > limit {
		lastMessage := messages[limit-1]
		newNextToken = pagination.GenerateToken(lastMessage.ID, lastMessage.CreatedAt.Unix())
		messages = messages[:limit]
	}

	return messages, newNextToken, nil
}

func (r *messageRepository) Create(message *entity.Message) error {
	return r.db.Create(message).Error
}
