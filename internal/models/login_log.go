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
