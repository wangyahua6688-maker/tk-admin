package models

import "time"

// User 用户模型
// 说明：
// 1. 使用显式 TableName 与 DAO 中的 SQL 统一，避免单复数不一致。
// 2. Roles 为用户-角色多对多关系，关联表为 user_roles。
type User struct {
	ID           uint   `gorm:"primaryKey"`
	Username     string `gorm:"unique;not null"`
	PasswordHash string
	Email        string
	Status       int `gorm:"default:1"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
	Roles        []Role `gorm:"many2many:sys_user_roles"`
	RefreshToken string `gorm:"default:null"` // 保留字段（兼容历史实现）
}

// TableName 指定 users 表名，避免使用保留关键字 user。
func (User) TableName() string {
	return "sys_users"
}
