package biz

import (
	"context"

	"go-admin/internal/models"
	"gorm.io/gorm"
)

// ExternalLinkDAO 外链数据访问层。
type ExternalLinkDAO struct {
	db *gorm.DB
}

// NewExternalLinkDAO 创建外链 DAO。
func NewExternalLinkDAO(db *gorm.DB) *ExternalLinkDAO {
	return &ExternalLinkDAO{db: db}
}

// ListExternalLinks 查询外链列表。
func (d *ExternalLinkDAO) ListExternalLinks(ctx context.Context, limit int) ([]models.WExternalLink, error) {
	// 判断条件并进入对应分支逻辑。
	if limit <= 0 || limit > 1000 {
		// 更新当前变量或字段值。
		limit = 200
	}
	// 声明当前变量。
	var items []models.WExternalLink
	// 定义并初始化当前变量。
	err := d.db.WithContext(ctx).
		// 调用Order完成当前处理。
		Order("position ASC, sort ASC, id DESC").
		// 调用Limit完成当前处理。
		Limit(limit).
		// 调用Find完成当前处理。
		Find(&items).Error
	// 返回当前处理结果。
	return items, err
}

// CreateExternalLink 新增外链。
func (d *ExternalLinkDAO) CreateExternalLink(ctx context.Context, item *models.WExternalLink) error {
	// 返回当前处理结果。
	return d.db.WithContext(ctx).Create(item).Error
}

// UpdateExternalLink 更新外链。
func (d *ExternalLinkDAO) UpdateExternalLink(ctx context.Context, id uint, updates map[string]interface{}) error {
	// 返回当前处理结果。
	return d.db.WithContext(ctx).Model(&models.WExternalLink{}).Where("id = ?", id).Updates(updates).Error
}

// DeleteExternalLink 删除外链。
func (d *ExternalLinkDAO) DeleteExternalLink(ctx context.Context, id uint) error {
	// 返回当前处理结果。
	return d.db.WithContext(ctx).Delete(&models.WExternalLink{}, id).Error
}
