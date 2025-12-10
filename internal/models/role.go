package models

import "gorm.io/gorm"

// Role 角色模型
type Role struct {
	gorm.Model
	Name        string       `gorm:"size:100;uniqueIndex" json:"name"`
	Code        string       `gorm:"size:100;uniqueIndex" json:"code"`
	Permissions []Permission `gorm:"many2many:sys_role_permissions;" json:"permissions"`
}

// TableName 指定 roles 表名，与 SQL 初始化脚本保持一致。
func (Role) TableName() string {
	return "sys_roles"
}
