package models

import "gorm.io/gorm"

type Menu struct {
	gorm.Model
	Title    string `gorm:"size:200" json:"title"`
	Path     string `gorm:"size:200" json:"path"`
	Icon     string `gorm:"size:100" json:"icon"`
	ParentID *uint  `json:"parent_id"`
	Order    int    `json:"order"`

	Permissions []Permission `gorm:"many2many:menu_permissions;" json:"permissions"`
	Children    []*Menu      `gorm:"-" json:"children"` // 构建树结构时用
}
