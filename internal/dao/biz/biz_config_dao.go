package biz

import "gorm.io/gorm"

// BizConfigDAO 聚合业务配置相关的 DAO。
// 说明：具体方法拆分在同包不同文件中，避免单文件过大。
type BizConfigDAO struct {
	db *gorm.DB // 数据库连接
}

// NewBizConfigDAO 创建业务配置 DAO。
func NewBizConfigDAO(db *gorm.DB) *BizConfigDAO {
	// 返回当前处理结果。
	return &BizConfigDAO{db: db}
}
