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

func (r *messageRepository) GetByChatID(chatID uint, limit int, page int) ([]entity.Message, int64, error) {
	limit = pagination.NormalizeLimit(limit)
	page = pagination.NormalizePage(page)
	offset := pagination.CalculateOffset(page, limit)

	countQuery := r.db.Model(&entity.Message{}).
		Where("chat_id = ?", chatID)

	var total int64
	if err := countQuery.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	query := r.db.Where(&entity.Message{ChatID: chatID}).
		Preload("Author").
		Preload("Chat.User").
		Order("messages.created_at DESC, messages.id DESC").
		Offset(offset).
		Limit(limit)

	var messages []entity.Message
	if err := query.Find(&messages).Error; err != nil {
		return nil, 0, err
	}

	return messages, total, nil
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
