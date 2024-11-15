package models

import (
	"time"

	"gorm.io/gorm"
)

type PasswordReset struct {
	Id        uint           `json:"id" gorm:"primaryKey"`
	Email     string         `json:"email"`
	Token     string         `json:"token" gorm:"type:varchar(100);uniqueIndex"`
	ExpiresAt time.Time      `json:"expires_at"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"deleted_at,omitempty" gorm:"index"`
}
