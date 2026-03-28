package biz

import (
	"context"

	"go-admin/internal/models"
)

// ListLotteryCategories 查询图库分类列表。
func (s *LotteryService) ListLotteryCategories(ctx context.Context, keyword string) ([]models.WLotteryCategory, error) {
	// 返回当前处理结果。
	return s.dao.ListLotteryCategories(ctx, keyword)
}

// CreateLotteryCategory 新增图库分类。
func (s *LotteryService) CreateLotteryCategory(ctx context.Context, item *models.WLotteryCategory) error {
	// 返回当前处理结果。
	return s.dao.CreateLotteryCategory(ctx, item)
}

// UpdateLotteryCategory 更新图库分类。
func (s *LotteryService) UpdateLotteryCategory(ctx context.Context, id uint, updates map[string]interface{}) error {
	// 返回当前处理结果。
	return s.dao.UpdateLotteryCategory(ctx, id, updates)
}

// DeleteLotteryCategory 删除图库分类。
func (s *LotteryService) DeleteLotteryCategory(ctx context.Context, id uint) error {
	// 返回当前处理结果。
	return s.dao.DeleteLotteryCategory(ctx, id)
}
