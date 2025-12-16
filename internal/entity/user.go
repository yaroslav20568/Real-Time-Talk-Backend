package entity

import (
	"strings"

	"gorm.io/gorm"
)

type User struct {
	ID        uint   `gorm:"primaryKey"`
	Email     string `gorm:"uniqueIndex;not null"`
	Password  string `gorm:"not null"`
	FirstName string `gorm:"column:first_name;not null"`
	LastName  string `gorm:"column:last_name;not null"`
	FullName  string `gorm:"column:full_name;not null"`
	Photo     string `gorm:"type:text"`
}

func (u *User) BeforeSave(tx *gorm.DB) error {
	parts := []string{}
	if strings.TrimSpace(u.FirstName) != "" {
		parts = append(parts, strings.TrimSpace(u.FirstName))
	}
	if strings.TrimSpace(u.LastName) != "" {
		parts = append(parts, strings.TrimSpace(u.LastName))
	}
	u.FullName = strings.Join(parts, " ")
	return nil
}

func (User) TableName() string {
	return "users"
}

