package models

import "time"

type User struct {
	ID           uint   `gorm:"primaryKey"`
	Username     string `gorm:"unique;not null"`
	PasswordHash string
	Email        string
	Status       int `gorm:"default:1"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
	Roles        []Role `gorm:"many2many:user_roles"`
	RefreshToken string `gorm:"default:null"` // 新增刷新Token字段
}
