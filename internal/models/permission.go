package models

import "gorm.io/gorm"

// Permission 权限模型
type Permission struct {
	gorm.Model
	Name   string `gorm:"size:100;uniqueIndex" json:"name"`
	Code   string `gorm:"size:100;uniqueIndex" json:"code"`
	Type   string `gorm:"size:50" json:"type"`
	Method string `gorm:"size:50" json:"method"`
	Path   string `gorm:"size:200" json:"path"`
}

// TableName 指定 permissions 表名，与 SQL 初始化脚本保持一致。
func (Permission) TableName() string {
	return "sys_permissions"
}
