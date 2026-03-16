package biz

import (
	"context"

	"go-admin-full/internal/models"
)

// ListSpecialLotteries 查询彩种配置列表。
func (s *BizConfigService) ListSpecialLotteries(ctx context.Context, limit int) ([]models.WSpecialLottery, error) {
	return s.dao.ListSpecialLotteries(ctx, limit)
}

// CreateSpecialLottery 新增彩种配置。
func (s *BizConfigService) CreateSpecialLottery(ctx context.Context, item *models.WSpecialLottery) error {
	return s.dao.CreateSpecialLottery(ctx, item)
}

// UpdateSpecialLottery 更新彩种配置。
func (s *BizConfigService) UpdateSpecialLottery(ctx context.Context, id uint, updates map[string]interface{}) error {
	return s.dao.UpdateSpecialLottery(ctx, id, updates)
}

// DeleteSpecialLottery 删除彩种配置。
func (s *BizConfigService) DeleteSpecialLottery(ctx context.Context, id uint) error {
	return s.dao.DeleteSpecialLottery(ctx, id)
}
