package models

import "time"

type LoginLog struct {
	ID        uint `gorm:"primaryKey"`
	UserID    uint
	Username  string
	IP        string
	Device    string
	Status    int
	CreatedAt time.Time
}

// TableName 指定 login_logs 表名。
func (LoginLog) TableName() string {
	return "sys_login_logs"
}
