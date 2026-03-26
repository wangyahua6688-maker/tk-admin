package biz

import (
	"context"

	"go-admin-full/internal/models"
)

// ListOfficialPosts 查询官方发帖列表。
func (s *BizConfigService) ListOfficialPosts(ctx context.Context, limit int) ([]models.WPostArticle, error) {
	// 直通 DAO 查询列表。
	return s.dao.ListOfficialPosts(ctx, limit)
}

// CreateOfficialPost 新增官方发帖。
func (s *BizConfigService) CreateOfficialPost(ctx context.Context, item *models.WPostArticle) error {
	// 直通 DAO 新增记录。
	return s.dao.CreateOfficialPost(ctx, item)
}

// UpdateOfficialPost 更新官方发帖。
func (s *BizConfigService) UpdateOfficialPost(ctx context.Context, id uint, updates map[string]interface{}) error {
	// 直通 DAO 更新记录。
	return s.dao.UpdateOfficialPost(ctx, id, updates)
}

// DeleteOfficialPost 删除官方发帖。
func (s *BizConfigService) DeleteOfficialPost(ctx context.Context, id uint) error {
	// 直通 DAO 删除记录。
	return s.dao.DeleteOfficialPost(ctx, id)
}
