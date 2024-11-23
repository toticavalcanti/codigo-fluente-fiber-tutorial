package models

import (
	"time"

	"gorm.io/gorm"
)

type PasswordReset struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	Email     string         `json:"email" gorm:"not null;type:varchar(255)"`
	Token     string         `json:"token" gorm:"not null;type:varchar(100);uniqueIndex"`
	ExpiresAt time.Time      `json:"expires_at" gorm:"type:timestamp;not null"`
	CreatedAt time.Time      `json:"created_at" gorm:"autoCreateTime;type:timestamp"`
	UpdatedAt time.Time      `json:"updated_at" gorm:"autoUpdateTime;type:timestamp"`
	DeletedAt gorm.DeletedAt `json:"deleted_at,omitempty" gorm:"index"`
}
