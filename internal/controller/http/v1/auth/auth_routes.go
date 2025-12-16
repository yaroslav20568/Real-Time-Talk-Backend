package auth

import (
	"gin-real-time-talk/internal/usecase/auth_usecase"
	"gin-real-time-talk/internal/usecase/repository"
	"gin-real-time-talk/pkg/email"
	"gin-real-time-talk/pkg/middleware"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func SetupAuthRoutes(api *gin.RouterGroup, db *gorm.DB) {
	userRepo := repository.NewUserRepository(db)
	emailService := email.NewEmailService()
	authUsecase := auth_usecase.NewAuthUsecase(userRepo, emailService)
	authController := NewAuthController(authUsecase)

	auth := api.Group("/auth")
	{
		auth.POST("/register", authController.Register)
		auth.POST("/login", authController.Login)
		auth.POST("/verify", authController.VerifyCode)
		auth.POST("/resend-code", authController.ResendCode)
		auth.POST("/refresh", authController.Refresh)
		auth.GET("/me", middleware.AuthMiddleware(authUsecase), authController.Me)
	}
}
