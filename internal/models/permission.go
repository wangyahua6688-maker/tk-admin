package models

import "gorm.io/gorm"

// Permission 权限模型
type Permission struct {
	// 处理当前语句逻辑。
	gorm.Model
	// 处理当前语句逻辑。
	Name string `gorm:"size:100;uniqueIndex" json:"name"`
	// 处理当前语句逻辑。
	Code string `gorm:"size:100;uniqueIndex" json:"code"`
	// 处理当前语句逻辑。
	Type string `gorm:"size:50" json:"type"`
	// 处理当前语句逻辑。
	Method string `gorm:"size:50" json:"method"`
	// 处理当前语句逻辑。
	Path string `gorm:"size:200" json:"path"`
}

// TableName 指定 permissions 表名，与 SQL 初始化脚本保持一致。
func (Permission) TableName() string {
	return "sys_permissions"
}
