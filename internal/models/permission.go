package models

import "time"

type Permission struct {
	ID        uint `gorm:"primaryKey"`
	Name      string
	Code      string `gorm:"unique;not null"`
	Type      string
	CreatedAt time.Time
	UpdatedAt time.Time
}
