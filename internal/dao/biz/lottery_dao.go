package biz

import "gorm.io/gorm"

// LotteryDAO 彩票业务数据访问层。
// 该 DAO 聚合彩种、图库、开奖记录、图库分类相关表访问。
type LotteryDAO struct {
	db *gorm.DB
}

// NewLotteryDAO 创建彩票业务 DAO。
func NewLotteryDAO(db *gorm.DB) *LotteryDAO {
	return &LotteryDAO{db: db}
}
