package app

import (
	"fmt"

	"gorm.io/gorm"
)

func Migrate(db *gorm.DB) error {
	if db == nil {
		return fmt.Errorf("database connection is not initialized")
	}

	return db.AutoMigrate()
}
