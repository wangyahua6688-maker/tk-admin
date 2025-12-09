package models

import "time"

type RefreshTokenRecord struct {
	ID        uint `gorm:"primaryKey"`
	UserID    uint
	DeviceID  string
	Token     string
	ExpiresAt time.Time
	CreatedAt time.Time
}
