package models

import "gorm.io/gorm"

// Role 角色模型
type Role struct {
	gorm.Model
	Name        string       `gorm:"size:100;uniqueIndex" json:"name"`
	Code        string       `gorm:"size:100;uniqueIndex" json:"code"`
	Permissions []Permission `gorm:"many2many:role_permissions;" json:"permissions"`
}
