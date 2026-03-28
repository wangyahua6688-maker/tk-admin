package biz

import (
	"context"

	"go-admin/internal/models"
)

// ListSpecialLotteries 查询彩种配置列表。
func (s *LotteryService) ListSpecialLotteries(ctx context.Context, limit int) ([]models.WSpecialLottery, error) {
	// 返回当前处理结果。
	return s.dao.ListSpecialLotteries(ctx, limit)
}

// CreateSpecialLottery 新增彩种配置。
func (s *LotteryService) CreateSpecialLottery(ctx context.Context, item *models.WSpecialLottery) error {
	// 返回当前处理结果。
	return s.dao.CreateSpecialLottery(ctx, item)
}

// UpdateSpecialLottery 更新彩种配置。
func (s *LotteryService) UpdateSpecialLottery(ctx context.Context, id uint, updates map[string]interface{}) error {
	// 返回当前处理结果。
	return s.dao.UpdateSpecialLottery(ctx, id, updates)
}

// DeleteSpecialLottery 删除彩种配置。
func (s *LotteryService) DeleteSpecialLottery(ctx context.Context, id uint) error {
	// 返回当前处理结果。
	return s.dao.DeleteSpecialLottery(ctx, id)
}
