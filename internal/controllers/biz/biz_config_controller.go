package biz

import (
	bizsvc "go-admin/internal/services/biz"

	"gorm.io/gorm"
)

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
	svc *bizsvc.BizConfigService // 业务配置服务聚合
}

// NewBizConfigController 创建BizConfigController实例。
func NewBizConfigController(db *gorm.DB) *BizConfigController {
	// 返回当前处理结果。
	return &BizConfigController{svc: bizsvc.NewBizConfigService(db)}
}
