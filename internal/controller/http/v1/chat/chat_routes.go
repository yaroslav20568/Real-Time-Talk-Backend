package chat

import (
	"gin-real-time-talk/internal/entity/interfaces"
	"gin-real-time-talk/internal/usecase/chat_usecase"
	"gin-real-time-talk/internal/usecase/repository"
	"gin-real-time-talk/pkg/middleware"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func SetupChatRoutes(api *gin.RouterGroup, db *gorm.DB, authUsecase interfaces.AuthUsecase) {
	chatRepo := repository.NewChatRepository(db)
	messageRepo := repository.NewMessageRepository(db)
	chatUsecase := chat_usecase.NewChatUsecase(chatRepo, messageRepo)
	chatController := NewChatController(chatUsecase)

	chats := api.Group("/chats")
	chats.Use(middleware.AuthMiddleware(authUsecase))
	{
		chats.GET("", chatController.GetUserChats)
		chats.GET("/:id/messages", chatController.GetChatMessages)
		chats.POST("/messages", chatController.CreateMessage)
	}
}
