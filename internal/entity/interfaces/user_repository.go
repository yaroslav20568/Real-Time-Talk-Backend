package interfaces

import "gin-real-time-talk/internal/entity"

type UserRepository interface {
	Create(user *entity.User) error
	GetByEmail(email string) (*entity.User, error)
	GetByID(id uint) (*entity.User, error)
	Update(user *entity.User) error
}
