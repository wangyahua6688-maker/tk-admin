package controllers

import "gorm.io/gorm"

// UserOpsController 用户管理域入口控制器。
// 具体业务方法按功能拆分在多个文件中：
// - client user
// - post article
// - hot comment
type UserOpsController struct {
	db *gorm.DB
}

func NewUserOpsController(db *gorm.DB) *UserOpsController {
	return &UserOpsController{db: db}
}
