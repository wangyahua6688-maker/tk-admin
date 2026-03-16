package biz

import (
	"context"

	"go-admin-full/internal/models"
)

// ListPostArticles 查询帖子列表（按官方/网友区分）。
func (d *UserOpsDAO) ListPostArticles(ctx context.Context, isOfficial int8, limit int) ([]models.WPostArticle, error) {
	query := d.db.WithContext(ctx).Model(&models.WPostArticle{}).Where("is_official = ?", isOfficial).Order("id DESC")
	if limit > 0 {
		query = query.Limit(limit)
	}
	items := make([]models.WPostArticle, 0)
	if err := query.Find(&items).Error; err != nil {
		return nil, err
	}
	return items, nil
}

// CreatePostArticle 新增帖子。
func (d *UserOpsDAO) CreatePostArticle(ctx context.Context, item *models.WPostArticle) error {
	return d.db.WithContext(ctx).Create(item).Error
}

// UpdatePostArticle 更新帖子。
func (d *UserOpsDAO) UpdatePostArticle(ctx context.Context, id uint, updates map[string]interface{}) error {
	return d.db.WithContext(ctx).Model(&models.WPostArticle{}).Where("id = ?", id).Updates(updates).Error
}

// DeletePostArticle 删除帖子。
func (d *UserOpsDAO) DeletePostArticle(ctx context.Context, id uint) error {
	return d.db.WithContext(ctx).Delete(&models.WPostArticle{}, id).Error
}

// GetPostArticleByID 查询帖子是否存在。
func (d *UserOpsDAO) GetPostArticleByID(ctx context.Context, id uint) (*models.WPostArticle, error) {
	var post models.WPostArticle
	if err := d.db.WithContext(ctx).Select("id,user_id,is_official,status").Where("id = ?", id).First(&post).Error; err != nil {
		return nil, err
	}
	return &post, nil
}
