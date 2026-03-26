package biz

import (
	"context"
	"strings"

	"go-admin/internal/models"
	"gorm.io/gorm"
)

// ListDrawRecords 查询开奖记录列表（支持彩种与关键字筛选）。
func (d *BizConfigDAO) ListDrawRecords(ctx context.Context, specialLotteryID uint, keyword string, limit int) ([]models.WDrawRecord, error) {
	// 构建基础查询并设置排序。
	query := d.db.WithContext(ctx).Model(&models.WDrawRecord{}).Order("draw_at DESC, id DESC")
	// 限制数量，防止一次性拉取过多数据。
	if limit > 0 {
		// 更新当前变量或字段值。
		query = query.Limit(limit)
	}
	// 彩种筛选。
	if specialLotteryID > 0 {
		// 更新当前变量或字段值。
		query = query.Where("special_lottery_id = ?", specialLotteryID)
	}
	// 关键字筛选（期号模糊匹配）。
	keyword = strings.TrimSpace(keyword)
	// 判断条件并进入对应分支逻辑。
	if keyword != "" {
		// 更新当前变量或字段值。
		query = query.Where("issue LIKE ?", "%"+keyword+"%")
	}

	// 执行查询。
	items := make([]models.WDrawRecord, 0)
	// 判断条件并进入对应分支逻辑。
	if err := query.Find(&items).Error; err != nil {
		// 返回当前处理结果。
		return nil, err
	}
	// 返回当前处理结果。
	return items, nil
}

// GetDrawRecordByID 查询单条开奖记录。
func (d *BizConfigDAO) GetDrawRecordByID(ctx context.Context, id uint) (*models.WDrawRecord, error) {
	// 查询单条记录。
	var item models.WDrawRecord
	// 判断条件并进入对应分支逻辑。
	if err := d.db.WithContext(ctx).First(&item, id).Error; err != nil {
		// 返回当前处理结果。
		return nil, err
	}
	// 返回当前处理结果。
	return &item, nil
}

// CreateDrawRecordTx 在事务中创建开奖记录。
func (d *BizConfigDAO) CreateDrawRecordTx(tx *gorm.DB, item *models.WDrawRecord) error {
	// 直接写入开奖记录。
	return tx.Create(item).Error
}

// UpdateDrawRecordTx 在事务中更新开奖记录。
func (d *BizConfigDAO) UpdateDrawRecordTx(tx *gorm.DB, id uint, updates map[string]interface{}) error {
	// 按主键更新记录。
	return tx.Model(&models.WDrawRecord{}).Where("id = ?", id).Updates(updates).Error
}

// ResetCurrentDrawRecordTx 将指定彩种的 current 标记清零。
func (d *BizConfigDAO) ResetCurrentDrawRecordTx(tx *gorm.DB, specialLotteryID uint, excludeID uint) error {
	// 清理同彩种的 current 标识（排除当前记录）。
	query := tx.Model(&models.WDrawRecord{}).Where("special_lottery_id = ?", specialLotteryID)
	// 判断条件并进入对应分支逻辑。
	if excludeID > 0 {
		// 更新当前变量或字段值。
		query = query.Where("id <> ?", excludeID)
	}
	// 返回当前处理结果。
	return query.Update("is_current", 0).Error
}

// DeleteDrawRecord 删除开奖记录。
func (d *BizConfigDAO) DeleteDrawRecord(ctx context.Context, id uint) error {
	// 直接按主键删除。
	return d.db.WithContext(ctx).Delete(&models.WDrawRecord{}, id).Error
}
