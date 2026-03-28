package biz

import (
	"context"

	bizdao "go-admin/internal/dao/biz"
	"go-admin/internal/models"
	"gorm.io/gorm"
)

// OfficialPostService 官方发帖业务服务。
type OfficialPostService struct {
	dao *bizdao.OfficialPostDAO
}

// NewOfficialPostService 创建官方发帖服务。
func NewOfficialPostService(db *gorm.DB) *OfficialPostService {
	return &OfficialPostService{dao: bizdao.NewOfficialPostDAO(db)}
}

// ListOfficialPosts 查询官方发帖列表。
func (s *OfficialPostService) ListOfficialPosts(ctx context.Context, limit int) ([]models.WPostArticle, error) {
	// 直通 DAO 查询列表。
	return s.dao.ListOfficialPosts(ctx, limit)
}

// CreateOfficialPost 新增官方发帖。
func (s *OfficialPostService) CreateOfficialPost(ctx context.Context, item *models.WPostArticle) error {
	// 直通 DAO 新增记录。
	return s.dao.CreateOfficialPost(ctx, item)
}

// UpdateOfficialPost 更新官方发帖。
func (s *OfficialPostService) UpdateOfficialPost(ctx context.Context, id uint, updates map[string]interface{}) error {
	// 直通 DAO 更新记录。
	return s.dao.UpdateOfficialPost(ctx, id, updates)
}

// DeleteOfficialPost 删除官方发帖。
func (s *OfficialPostService) DeleteOfficialPost(ctx context.Context, id uint) error {
	// 直通 DAO 删除记录。
	return s.dao.DeleteOfficialPost(ctx, id)
}
