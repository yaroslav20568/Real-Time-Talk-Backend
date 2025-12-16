package interfaces

import "gin-real-time-talk/internal/entity"

type AuthUsecase interface {
	Register(email, password, firstName, lastName string) (*entity.User, error)
	Login(email, password string) (string, string, *entity.User, error)
	SendTwoFactorCode(email string) error
	VerifyTwoFactorCode(email, code string) (string, string, *entity.User, error)
	RefreshToken(refreshToken string) (string, string, *entity.User, error)
	ValidateAccessToken(token string) (*entity.User, error)
}
