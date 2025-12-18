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

	allChatsQuery := r.db.Model(&entity.Chat{}).
		Where("chats.user_id = ?", userID)

	var allChats []entity.Chat
	if err := allChatsQuery.Find(&allChats).Error; err != nil {
		return nil, 0, err
	}

	if len(allChats) == 0 {
		return []entity.Chat{}, 0, nil
	}

	allChatIDs := make([]uint, len(allChats))
	for i := range allChats {
		allChatIDs[i] = allChats[i].ID
	}

	var allOtherUsers []struct {
		ChatID uint `gorm:"column:chat_id"`
		UserID uint `gorm:"column:user_id"`
	}

	r.db.Raw(`
		SELECT DISTINCT ON (messages.chat_id) 
			messages.chat_id, 
			messages.author_id as user_id
		FROM messages
		WHERE messages.chat_id IN ? 
			AND messages.author_id != ?
		ORDER BY messages.chat_id, messages.created_at DESC
	`, allChatIDs, userID).Scan(&allOtherUsers)

	allOtherUserMap := make(map[uint]uint)
	for _, ou := range allOtherUsers {
		allOtherUserMap[ou.ChatID] = ou.UserID
	}

	var allOtherUserIDs []uint
	for _, uid := range allOtherUserMap {
		allOtherUserIDs = append(allOtherUserIDs, uid)
	}

	var allUsers []entity.User
	if len(allOtherUserIDs) > 0 {
		r.db.Where("id IN ?", allOtherUserIDs).Find(&allUsers)
	}

	allUserMap := make(map[uint]entity.User)
	for _, user := range allUsers {
		allUserMap[user.ID] = user
	}

	for i := range allChats {
		if otherUserID, exists := allOtherUserMap[allChats[i].ID]; exists {
			if otherUser, exists := allUserMap[otherUserID]; exists {
				allChats[i].User = otherUser
			}
		}
	}

	var unreadMessages []struct {
		ChatID   uint `gorm:"column:chat_id"`
		AuthorID uint `gorm:"column:author_id"`
	}

	if len(allChatIDs) > 0 {
		r.db.Model(&entity.Message{}).
			Select("chat_id, author_id").
			Where("chat_id IN ? AND is_read = ?", allChatIDs, false).
			Find(&unreadMessages)
	}

	unreadCountMap := make(map[uint]int)
	for _, msg := range unreadMessages {
		if expectedAuthorID, exists := allOtherUserMap[msg.ChatID]; exists && msg.AuthorID == expectedAuthorID {
			unreadCountMap[msg.ChatID]++
		}
	}

	for i := range allChats {
		if _, exists := allOtherUserMap[allChats[i].ID]; exists {
			allChats[i].UnreadCount = unreadCountMap[allChats[i].ID]
		} else {
			allChats[i].UnreadCount = 0
		}
	}

	var filteredChats []entity.Chat
	if search != "" {
		for i := range allChats {
			matchesSearch := false
			if allChats[i].User.ID > 0 && allChats[i].User.FullName != "" {
				matchesSearch = r.db.Where("id = ? AND full_name ILIKE ?", allChats[i].User.ID, "%"+search+"%").First(&entity.User{}).Error == nil
			}
			if !matchesSearch && allChats[i].LastMessageText != nil {
				text := *allChats[i].LastMessageText
				var matches bool
				r.db.Raw("SELECT ? ILIKE ?", text, "%"+search+"%").Scan(&matches)
				matchesSearch = matches
			}
			if matchesSearch {
				filteredChats = append(filteredChats, allChats[i])
			}
		}
	} else {
		filteredChats = allChats
	}

	for i := 0; i < len(filteredChats)-1; i++ {
		for j := i + 1; j < len(filteredChats); j++ {
			if filteredChats[i].UpdatedAt.Before(filteredChats[j].UpdatedAt) ||
				(filteredChats[i].UpdatedAt.Equal(filteredChats[j].UpdatedAt) && filteredChats[i].ID < filteredChats[j].ID) {
				filteredChats[i], filteredChats[j] = filteredChats[j], filteredChats[i]
			}
		}
	}

	total := int64(len(filteredChats))

	start := offset
	if start > len(filteredChats) {
		start = len(filteredChats)
	}
	end := start + limit
	if end > len(filteredChats) {
		end = len(filteredChats)
	}

	chats := filteredChats[start:end]

	chatIDs := make([]uint, len(chats))
	for i := range chats {
		chatIDs[i] = chats[i].ID
	}

	if len(chatIDs) > 0 {
		var chatsWithMessages []entity.Chat
		r.db.Where("id IN ?", chatIDs).
			Preload("LastMessage.Author").
			Find(&chatsWithMessages)

		chatsMap := make(map[uint]*entity.Chat)
		for i := range chats {
			chatsMap[chats[i].ID] = &chats[i]
		}

		for i := range chatsWithMessages {
			if chat, exists := chatsMap[chatsWithMessages[i].ID]; exists {
				chat.LastMessage = chatsWithMessages[i].LastMessage
				if chat.LastMessage != nil {
					chat.LastMessage.Chat = nil
				}
			}
		}
	}

	return chats, total, nil
}

func (r *chatRepository) GetByID(id uint) (*entity.Chat, error) {
	var chat entity.Chat
	err := r.db.Preload("LastMessage.Author").First(&chat, id).Error
	if err != nil {
		return nil, err
	}

	var otherUser struct {
		UserID uint `gorm:"column:user_id"`
	}

	r.db.Raw(`
		SELECT messages.author_id as user_id
		FROM messages
		WHERE messages.chat_id = ?
			AND messages.author_id != ?
		ORDER BY messages.created_at DESC
		LIMIT 1
	`, id, chat.UserID).Scan(&otherUser)

	if otherUser.UserID > 0 {
		var user entity.User
		if err := r.db.First(&user, otherUser.UserID).Error; err == nil {
			chat.User = user
		}

		var unreadCount int64
		r.db.Model(&entity.Message{}).
			Where("chat_id = ? AND author_id = ? AND is_read = ?", id, otherUser.UserID, false).
			Count(&unreadCount)
		chat.UnreadCount = int(unreadCount)
	} else {
		chat.UnreadCount = 0
	}

	if chat.LastMessage != nil {
		chat.LastMessage.Chat = nil
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
