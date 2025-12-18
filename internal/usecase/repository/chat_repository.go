package repository

import (
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

func (r *chatRepository) GetByUserID(userID uint, limit int, page int, search string) ([]entity.Chat, int64, error) {
	limit = pagination.NormalizeLimit(limit)
	page = pagination.NormalizePage(page)
	offset := pagination.CalculateOffset(page, limit)

	query := r.db.Model(&entity.Chat{}).
		Where("chats.user_id = ?", userID)

	countQuery := r.db.Model(&entity.Chat{}).
		Where("chats.user_id = ?", userID)

	if search != "" {
		searchPattern := "%" + search + "%"
		query = query.Joins("User").
			Where("users.full_name ILIKE ? OR chats.last_message_text ILIKE ?", searchPattern, searchPattern)
		countQuery = countQuery.Joins("User").
			Where("users.full_name ILIKE ? OR chats.last_message_text ILIKE ?", searchPattern, searchPattern)
	}

	var total int64
	if err := countQuery.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	query = query.Order("chats.updated_at DESC, chats.id DESC").
		Preload("User").
		Preload("LastMessage").
		Offset(offset).
		Limit(limit)

	var chats []entity.Chat
	if err := query.Find(&chats).Error; err != nil {
		return nil, 0, err
	}

	return chats, total, nil
}

func (r *chatRepository) GetByID(id uint) (*entity.Chat, error) {
	var chat entity.Chat
	err := r.db.Preload("User").Preload("LastMessage").First(&chat, id).Error
	if err != nil {
		return nil, err
	}
	return &chat, nil
}

func (r *chatRepository) FindOrCreateChatByUsers(senderID uint, recipientID uint) (*entity.Chat, error) {
	var chat entity.Chat
	
	err := r.db.Where("chats.user_id = ?", senderID).
		Joins("JOIN messages ON messages.chat_id = chats.id").
		Where("messages.author_id = ? OR messages.author_id = ?", senderID, recipientID).
		Group("chats.id").
		First(&chat).Error

	if err == nil {
		return &chat, nil
	}

	if err != gorm.ErrRecordNotFound {
		return nil, err
	}

	newChat := entity.Chat{
		UserID:      senderID,
		UnreadCount: 0,
	}

	if err := r.db.Create(&newChat).Error; err != nil {
		return nil, err
	}

	return &newChat, nil
}

func (r *chatRepository) Create(chat *entity.Chat) error {
	return r.db.Create(chat).Error
}

func (r *chatRepository) Update(chat *entity.Chat) error {
	return r.db.Save(chat).Error
}
