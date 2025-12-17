package repository

import (
	"time"

	"gin-real-time-talk/internal/entity"
	"gin-real-time-talk/internal/entity/interfaces"
	"gin-real-time-talk/pkg/pagination"

	"gorm.io/gorm"
)

type chatRepository struct {
	db *gorm.DB
}

func NewChatRepository(db *gorm.DB) interfaces.ChatRepository {
	return &chatRepository{
		db: db,
	}
}

func (r *chatRepository) GetByUserID(userID uint, limit int, nextToken string) ([]entity.Chat, string, error) {
	limit = pagination.NormalizeLimit(limit)

	query := r.db.Where("user_id = ?", userID).
		Preload("User").
		Preload("LastMessage").
		Order("updated_at DESC, id DESC")

	if nextToken != "" {
		tokenData, err := pagination.ParseToken(nextToken)
		if err != nil {
			return nil, "", err
		}

		if tokenData != nil {
			if tokenData.Timestamp > 0 {
				lastUpdatedAt := time.Unix(tokenData.Timestamp, 0)
				query = query.Where("(updated_at < ? OR (updated_at = ? AND id < ?))", lastUpdatedAt, lastUpdatedAt, tokenData.ID)
			} else {
				query = query.Where("id < ?", tokenData.ID)
			}
		}
	}

	var chats []entity.Chat
	if err := query.Limit(limit + 1).Find(&chats).Error; err != nil {
		return nil, "", err
	}

	var newNextToken string
	if len(chats) > limit {
		lastChat := chats[limit-1]
		newNextToken = pagination.GenerateToken(lastChat.ID, lastChat.UpdatedAt.Unix())
		chats = chats[:limit]
	}

	return chats, newNextToken, nil
}

func (r *chatRepository) GetByID(id uint) (*entity.Chat, error) {
	var chat entity.Chat
	err := r.db.Preload("User").Preload("LastMessage").First(&chat, id).Error
	if err != nil {
		return nil, err
	}
	return &chat, nil
}
