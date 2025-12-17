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

func (r *chatRepository) GetByUserID(userID uint, limit int, nextToken string, search string) ([]entity.Chat, string, error) {
	limit = pagination.NormalizeLimit(limit)

	query := r.db.Where(&entity.Chat{UserID: userID}).
		Preload("User").
		Preload("LastMessage")

	if search != "" {
		searchPattern := "%" + search + "%"
		query = query.Joins("User").
			Where("users.full_name ILIKE ? OR chats.last_message_text ILIKE ?", searchPattern, searchPattern)
	}

	query = query.Order("chats.updated_at DESC, chats.id DESC")

	if nextToken != "" {
		tokenData, err := pagination.ParseToken(nextToken)
		if err != nil {
			return nil, "", err
		}

		if tokenData != nil {
			if tokenData.Timestamp > 0 {
				lastUpdatedAt := time.Unix(tokenData.Timestamp, 0)
				query = query.Where("(chats.updated_at < ? OR (chats.updated_at = ? AND chats.id < ?))", lastUpdatedAt, lastUpdatedAt, tokenData.ID)
			} else {
				query = query.Where("chats.id < ?", tokenData.ID)
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
