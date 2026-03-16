package biz

import (
	bizsvc "go-admin-full/internal/services/biz"
	"gorm.io/gorm"
)

// UserOpsController 用户管理域入口控制器。
// 具体业务方法按功能拆分在多个文件中：
// - client user
// - post article
// - hot comment
type UserOpsController struct {
	svc *bizsvc.UserOpsService // 用户运营服务聚合
}

func NewUserOpsController(db *gorm.DB) *UserOpsController {
	return &UserOpsController{svc: bizsvc.NewUserOpsService(db)}
}
