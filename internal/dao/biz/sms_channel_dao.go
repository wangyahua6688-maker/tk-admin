package biz

import (
	"context"

	"go-admin/internal/models"
	"gorm.io/gorm"
)

// SMSChannelDAO 短信通道数据访问层。
type SMSChannelDAO struct {
	db *gorm.DB
}

// NewSMSChannelDAO 创建短信通道 DAO。
func NewSMSChannelDAO(db *gorm.DB) *SMSChannelDAO {
	return &SMSChannelDAO{db: db}
}

// ListSMSChannels 查询短信通道配置列表。
func (d *SMSChannelDAO) ListSMSChannels(ctx context.Context, status *int, limit int) ([]models.WSMSChannel, error) {
	// 定义并初始化当前变量。
	query := d.db.WithContext(ctx).Model(&models.WSMSChannel{}).Order("status DESC, id ASC")
	// 判断条件并进入对应分支逻辑。
	if limit > 0 {
		// 更新当前变量或字段值。
		query = query.Limit(limit)
	}
	// 仅当传入 0/1 时筛选状态。
	if status != nil {
		// 更新当前变量或字段值。
		query = query.Where("status = ?", *status)
	}
	// 定义并初始化当前变量。
	items := make([]models.WSMSChannel, 0)
	// 判断条件并进入对应分支逻辑。
	if err := query.Find(&items).Error; err != nil {
		// 返回当前处理结果。
		return nil, err
	}
	// 返回当前处理结果。
	return items, nil
}

// CreateSMSChannel 新增短信通道配置。
func (d *SMSChannelDAO) CreateSMSChannel(ctx context.Context, item *models.WSMSChannel) error {
	// 返回当前处理结果。
	return d.db.WithContext(ctx).Create(item).Error
}

// UpdateSMSChannel 更新短信通道配置。
func (d *SMSChannelDAO) UpdateSMSChannel(ctx context.Context, id uint, updates map[string]interface{}) error {
	// 返回当前处理结果。
	return d.db.WithContext(ctx).Model(&models.WSMSChannel{}).Where("id = ?", id).Updates(updates).Error
}

// DeleteSMSChannel 删除短信通道配置。
func (d *SMSChannelDAO) DeleteSMSChannel(ctx context.Context, id uint) error {
	// 返回当前处理结果。
	return d.db.WithContext(ctx).Delete(&models.WSMSChannel{}, id).Error
}
