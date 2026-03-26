package models

import "time"

// User 用户模型
// 说明：
// 1. 使用显式 TableName 与 DAO 中的 SQL 统一，避免单复数不一致。
// 2. Roles 为用户-角色多对多关系，关联表为 sys_user_roles。
type User struct {
	// 处理当前语句逻辑。
	ID uint `gorm:"primaryKey"`
	// 处理当前语句逻辑。
	Username string `gorm:"unique;not null"`
	// 处理当前语句逻辑。
	PasswordHash string
	// 处理当前语句逻辑。
	Email  string
	Avatar string `gorm:"size:255;default:''"` // 用户头像 URL（可为空）
	// 处理当前语句逻辑。
	Status int `gorm:"default:1"`
	// 处理当前语句逻辑。
	CreatedAt time.Time
	// 处理当前语句逻辑。
	UpdatedAt time.Time
	// 处理当前语句逻辑。
	Roles        []Role `gorm:"many2many:sys_user_roles"`
	RefreshToken string `gorm:"default:null"` // 保留字段（兼容历史实现）
}

// TableName 指定 sys_users 表名，避免使用保留关键字 user。
func (User) TableName() string {
	return "sys_users"
}
