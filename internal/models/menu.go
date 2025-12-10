package models

import (
	"time"

	"gorm.io/gorm"
)

// Menu 菜单模型
// 说明：
// 1. 菜单与权限是多对多关系（menu_permissions）。
// 2. ParentID + OrderNum 用于构建前端路由树。
type Menu struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	Title       string         `gorm:"size:200;not null" json:"title"`
	Path        string         `gorm:"size:200;not null" json:"path"`
	Icon        string         `gorm:"size:100" json:"icon"`
	ParentID    uint           `gorm:"default:0" json:"parent_id"`
	Component   string         `gorm:"size:200" json:"component"`
	OrderNum    int            `gorm:"column:order_num;default:0" json:"order_num"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
	Permissions []Permission   `gorm:"many2many:sys_menu_permissions;" json:"permissions"`
	Children    []*Menu        `gorm:"-" json:"children"` // 构建树结构时用
}

// TableName 指定 menus 表名，与 SQL 初始化脚本保持一致。
func (Menu) TableName() string {
	return "sys_menus"
}
