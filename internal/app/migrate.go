package app

import (
	"fmt"

	"gin-real-time-talk/internal/entity"

	"gorm.io/gorm"
)

func Migrate(db *gorm.DB) error {
	if db == nil {
		return fmt.Errorf("database connection is not initialized")
	}

	return db.AutoMigrate(
		&entity.User{},
		&entity.Chat{},
		&entity.Message{},
	)
}
