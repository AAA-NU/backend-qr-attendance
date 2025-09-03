package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type QRCode struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	UUID      string    `gorm:"uniqueIndex;not null" json:"uuid"`
	CreatedAt time.Time `json:"created_at"`
	ExpiresAt time.Time `json:"expires_at"`
	IsActive  bool      `gorm:"default:true" json:"is_active"`
}

func (q *QRCode) BeforeCreate(tx *gorm.DB) error {
	q.UUID = uuid.New().String()
	return nil
}

func (q *QRCode) IsExpired() bool {
	return time.Now().After(q.ExpiresAt)
}
