package biz

import (
	"context"

	bizdao "go-admin/internal/dao/biz"
	"go-admin/internal/models"
	"gorm.io/gorm"
)

// BroadcastService 广播业务服务。
type BroadcastService struct {
	dao *bizdao.BroadcastDAO
}

// NewBroadcastService 创建广播服务。
func NewBroadcastService(db *gorm.DB) *BroadcastService {
	return &BroadcastService{dao: bizdao.NewBroadcastDAO(db)}
}

// ListBroadcasts 查询系统广播列表。
func (s *BroadcastService) ListBroadcasts(ctx context.Context, limit int) ([]models.WBroadcast, error) {
	// 返回当前处理结果。
	return s.dao.ListBroadcasts(ctx, limit)
}

// CreateBroadcast 新增系统广播。
func (s *BroadcastService) CreateBroadcast(ctx context.Context, item *models.WBroadcast) error {
	// 返回当前处理结果。
	return s.dao.CreateBroadcast(ctx, item)
}

// UpdateBroadcast 更新系统广播。
func (s *BroadcastService) UpdateBroadcast(ctx context.Context, id uint, updates map[string]interface{}) error {
	// 返回当前处理结果。
	return s.dao.UpdateBroadcast(ctx, id, updates)
}

// DeleteBroadcast 删除系统广播。
func (s *BroadcastService) DeleteBroadcast(ctx context.Context, id uint) error {
	// 返回当前处理结果。
	return s.dao.DeleteBroadcast(ctx, id)
}
