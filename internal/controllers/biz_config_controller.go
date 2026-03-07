package controllers

import "gorm.io/gorm"

// BizConfigController 业务配置域入口控制器。
// 具体能力按功能拆分在多个文件中：
// - banner
// - broadcast
// - special_lottery
// - draw_record
// - lottery_info(图库内容)
// - official post
// - external link
type BizConfigController struct {
	db *gorm.DB
}

func NewBizConfigController(db *gorm.DB) *BizConfigController {
	return &BizConfigController{db: db}
}
