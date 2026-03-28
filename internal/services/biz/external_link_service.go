package biz

import (
	"context"

	bizdao "go-admin/internal/dao/biz"
	"go-admin/internal/models"
	"gorm.io/gorm"
)

// ExternalLinkService 外链业务服务。
type ExternalLinkService struct {
	dao *bizdao.ExternalLinkDAO
}

// NewExternalLinkService 创建外链服务。
func NewExternalLinkService(db *gorm.DB) *ExternalLinkService {
	return &ExternalLinkService{dao: bizdao.NewExternalLinkDAO(db)}
}

// ListExternalLinks 查询外链列表。
func (s *ExternalLinkService) ListExternalLinks(ctx context.Context, limit int) ([]models.WExternalLink, error) {
	// 返回当前处理结果。
	return s.dao.ListExternalLinks(ctx, limit)
}

// CreateExternalLink 新增外链。
func (s *ExternalLinkService) CreateExternalLink(ctx context.Context, item *models.WExternalLink) error {
	// 返回当前处理结果。
	return s.dao.CreateExternalLink(ctx, item)
}

// UpdateExternalLink 更新外链。
func (s *ExternalLinkService) UpdateExternalLink(ctx context.Context, id uint, updates map[string]interface{}) error {
	// 返回当前处理结果。
	return s.dao.UpdateExternalLink(ctx, id, updates)
}

// DeleteExternalLink 删除外链。
func (s *ExternalLinkService) DeleteExternalLink(ctx context.Context, id uint) error {
	// 返回当前处理结果。
	return s.dao.DeleteExternalLink(ctx, id)
}
