package models

import "time"

// LoginLog 定义LoginLog相关结构。
type LoginLog struct {
	// 处理当前语句逻辑。
	ID uint `gorm:"primaryKey"`
	// 处理当前语句逻辑。
	UserID uint
	// 处理当前语句逻辑。
	Username string
	// 处理当前语句逻辑。
	IP string
	// 处理当前语句逻辑。
	Device string
	// 处理当前语句逻辑。
	Status int
	// 处理当前语句逻辑。
	CreatedAt time.Time
}

// TableName 指定 login_logs 表名。
func (LoginLog) TableName() string {
	return "sys_login_logs"
}
