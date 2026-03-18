package models

import "gorm.io/gorm"

// Role 角色模型
type Role struct {
	// 处理当前语句逻辑。
	gorm.Model
	// 处理当前语句逻辑。
	Name string `gorm:"size:100;uniqueIndex" json:"name"`
	// 处理当前语句逻辑。
	Code string `gorm:"size:100;uniqueIndex" json:"code"`
	// 处理当前语句逻辑。
	Permissions []Permission `gorm:"many2many:sys_role_permissions;" json:"permissions"`
}

// TableName 指定 roles 表名，与 SQL 初始化脚本保持一致。
func (Role) TableName() string {
	return "sys_roles"
}
