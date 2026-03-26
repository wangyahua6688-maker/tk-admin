package models

import (
	"time"

	"gorm.io/gorm"
)

// Menu 菜单模型
// 说明：
// 1. 菜单与权限是多对多关系（sys_menu_permissions）。
// 2. ParentID + OrderNum 用于构建前端路由树。
type Menu struct {
	// 处理当前语句逻辑。
	ID uint `gorm:"primaryKey" json:"id"`
	// 处理当前语句逻辑。
	Title string `gorm:"size:200;not null" json:"title"`
	// 处理当前语句逻辑。
	Path string `gorm:"size:200;not null" json:"path"`
	// 处理当前语句逻辑。
	Icon string `gorm:"size:100" json:"icon"`
	// 处理当前语句逻辑。
	ParentID uint `gorm:"default:0" json:"parent_id"`
	// 处理当前语句逻辑。
	Component string `gorm:"size:200" json:"component"`
	// 处理当前语句逻辑。
	OrderNum int `gorm:"column:order_num;default:0" json:"order_num"`
	// 处理当前语句逻辑。
	CreatedAt time.Time `json:"created_at"`
	// 处理当前语句逻辑。
	UpdatedAt time.Time `json:"updated_at"`
	// 处理当前语句逻辑。
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
	// 处理当前语句逻辑。
	Permissions []Permission `gorm:"many2many:sys_menu_permissions;" json:"permissions"`
	Children    []*Menu      `gorm:"-" json:"children"` // 构建树结构时用
}

// TableName 指定 menus 表名，与 SQL 初始化脚本保持一致。
func (Menu) TableName() string {
	return "sys_menus"
}
