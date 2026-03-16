package biz

import (
	"context"

	"go-admin-full/internal/models"
	"gorm.io/gorm"
)

// ListLotteryInfos 查询图纸列表。
func (d *BizConfigDAO) ListLotteryInfos(ctx context.Context, limit int) ([]models.WLotteryInfo, error) {
	// 以更新时间倒序输出，保证后台编辑后能快速找到。
	query := d.db.WithContext(ctx).Model(&models.WLotteryInfo{}).Order("updated_at DESC, id DESC")
	if limit > 0 {
		query = query.Limit(limit)
	}
	items := make([]models.WLotteryInfo, 0)
	if err := query.Find(&items).Error; err != nil {
		return nil, err
	}
	return items, nil
}

// GetLotteryInfoByID 查询单条图纸。
func (d *BizConfigDAO) GetLotteryInfoByID(ctx context.Context, id uint) (*models.WLotteryInfo, error) {
	var item models.WLotteryInfo
	if err := d.db.WithContext(ctx).First(&item, id).Error; err != nil {
		return nil, err
	}
	return &item, nil
}

// CreateLotteryInfoTx 在事务中创建图纸。
func (d *BizConfigDAO) CreateLotteryInfoTx(tx *gorm.DB, item *models.WLotteryInfo) error {
	return tx.Create(item).Error
}

// UpdateLotteryInfoTx 在事务中更新图纸。
func (d *BizConfigDAO) UpdateLotteryInfoTx(tx *gorm.DB, id uint, updates map[string]interface{}) error {
	return tx.Model(&models.WLotteryInfo{}).Where("id = ?", id).Updates(updates).Error
}

// DeleteLotteryInfoTx 在事务中删除图纸。
func (d *BizConfigDAO) DeleteLotteryInfoTx(tx *gorm.DB, id uint) error {
	return tx.Delete(&models.WLotteryInfo{}, id).Error
}

// ListLotteryOptionsByInfoIDs 批量查询图纸的动物竞猜选项。
func (d *BizConfigDAO) ListLotteryOptionsByInfoIDs(ctx context.Context, infoIDs []uint) ([]models.WLotteryOption, error) {
	if len(infoIDs) == 0 {
		return []models.WLotteryOption{}, nil
	}
	options := make([]models.WLotteryOption, 0)
	if err := d.db.WithContext(ctx).
		Where("lottery_info_id IN ?", infoIDs).
		Order("sort ASC, id ASC").
		Find(&options).Error; err != nil {
		return nil, err
	}
	return options, nil
}

// ReplaceLotteryOptionsTx 全量替换图纸动物竞猜选项（保留旧票数）。
func (d *BizConfigDAO) ReplaceLotteryOptionsTx(tx *gorm.DB, infoID uint, optionNames []string) error {
	// 读取旧选项票数，保证编辑时票数不丢失。
	var oldRows []models.WLotteryOption
	if err := tx.Where("lottery_info_id = ?", infoID).Find(&oldRows).Error; err != nil {
		return err
	}
	oldVotes := make(map[string]int64, len(oldRows))
	for _, row := range oldRows {
		if row.OptionName == "" {
			continue
		}
		oldVotes[row.OptionName] = row.Votes
	}

	// 删除旧选项，避免残留。
	if err := tx.Where("lottery_info_id = ?", infoID).Delete(&models.WLotteryOption{}).Error; err != nil {
		return err
	}

	// 构造新选项。
	rows := make([]models.WLotteryOption, 0, len(optionNames))
	for idx, name := range optionNames {
		rows = append(rows, models.WLotteryOption{
			LotteryInfoID: infoID,
			OptionName:    name,
			Votes:         oldVotes[name],
			Sort:          idx + 1,
		})
	}
	if len(rows) == 0 {
		return nil
	}
	return tx.Create(&rows).Error
}
