package models

import "time"

// RefreshTokenRecord 定义RefreshTokenRecord相关结构。
type RefreshTokenRecord struct {
	// 处理当前语句逻辑。
	ID uint `gorm:"primaryKey"`
	// 处理当前语句逻辑。
	UserID uint
	// 处理当前语句逻辑。
	DeviceID string
	// 处理当前语句逻辑。
	Token string
	// 处理当前语句逻辑。
	ExpiresAt time.Time
	// 处理当前语句逻辑。
	CreatedAt time.Time
}

// TableName 指定 refresh_token_records 表名。
func (RefreshTokenRecord) TableName() string {
	return "sys_refresh_token_records"
}
