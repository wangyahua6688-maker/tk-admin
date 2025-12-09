package models

import "gorm.io/gorm"

// Permission 权限模型
type Permission struct {
	gorm.Model
	Name string `gorm:"size:100;uniqueIndex" json:"name"`
	Code string `gorm:"size:100;uniqueIndex" json:"code"`
	Type string `gorm:"size:50" json:"type"`
}
