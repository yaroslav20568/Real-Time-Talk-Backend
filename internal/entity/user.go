package entity

import (
	"strings"
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID                  uint       `gorm:"primaryKey" json:"id"`
	Email               string     `gorm:"uniqueIndex;not null" json:"email"`
	Password            string     `gorm:"not null" json:"-"`
	FirstName           string     `gorm:"column:first_name;not null" json:"firstName"`
	LastName            string     `gorm:"column:last_name;not null" json:"lastName"`
	FullName            string     `gorm:"column:full_name;not null" json:"fullName"`
	Photo               *string    `gorm:"type:text" json:"photo"`
	EmailVerified       bool       `gorm:"column:email_verified;default:false" json:"emailVerified"`
	TwoFactorCode       string     `gorm:"column:two_factor_code" json:"-"`
	TwoFactorExpiresAt  *time.Time `gorm:"column:two_factor_expires_at" json:"-"`
	TwoFactorVerifiedAt *time.Time `gorm:"column:two_factor_verified_at" json:"-"`
	CreatedAt           time.Time  `json:"createdAt"`
	UpdatedAt           time.Time  `json:"updatedAt"`
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

	if u.Photo != nil && strings.TrimSpace(*u.Photo) == "" {
		u.Photo = nil
	}

	return nil
}

func (User) TableName() string {
	return "users"
}
