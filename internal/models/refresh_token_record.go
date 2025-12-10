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

// TableName 指定 refresh_token_records 表名。
func (RefreshTokenRecord) TableName() string {
	return "sys_refresh_token_records"
}
