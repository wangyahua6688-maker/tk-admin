package biz

import (
	bizdao "go-admin-full/internal/dao/biz"
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
	return &BizConfigService{
		dao: bizdao.NewBizConfigDAO(db),
		db:  db,
	}
}
