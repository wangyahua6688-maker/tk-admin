package biz

import (
	bizdao "go-admin/internal/dao/biz"
	"gorm.io/gorm"
)

// LotteryService 彩票业务服务。
// 该服务聚合彩种、图库、开奖记录、图库分类相关能力。
type LotteryService struct {
	dao *bizdao.LotteryDAO
	db  *gorm.DB
}

// NewLotteryService 创建彩票业务服务。
func NewLotteryService(db *gorm.DB) *LotteryService {
	return &LotteryService{
		dao: bizdao.NewLotteryDAO(db),
		db:  db,
	}
}
