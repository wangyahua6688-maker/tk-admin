package biz

import (
	bizdao "go-admin/internal/dao/biz"
	"gorm.io/gorm"
)

// BizConfigService 聚合业务配置相关服务。
// 说明：具体方法拆分在同包不同文件中。
type BizConfigService struct {
	dao *bizdao.BizConfigDAO // 业务配置 DAO
	db  *gorm.DB             // 用于事务或复杂查询的数据库连接
}

// NewBizConfigService 创建业务配置服务。
func NewBizConfigService(db *gorm.DB) *BizConfigService {
	// 返回当前处理结果。
	return &BizConfigService{
		// 调用bizdao.NewBizConfigDAO完成当前处理。
		dao: bizdao.NewBizConfigDAO(db),
		// 处理当前语句逻辑。
		db: db,
	}
}
