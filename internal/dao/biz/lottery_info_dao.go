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
	// 判断条件并进入对应分支逻辑。
	if limit > 0 {
		// 更新当前变量或字段值。
		query = query.Limit(limit)
	}
	// 定义并初始化当前变量。
	items := make([]models.WLotteryInfo, 0)
	// 判断条件并进入对应分支逻辑。
	if err := query.Find(&items).Error; err != nil {
		// 返回当前处理结果。
		return nil, err
	}
	// 返回当前处理结果。
	return items, nil
}

// GetLotteryInfoByID 查询单条图纸。
func (d *BizConfigDAO) GetLotteryInfoByID(ctx context.Context, id uint) (*models.WLotteryInfo, error) {
	// 声明当前变量。
	var item models.WLotteryInfo
	// 判断条件并进入对应分支逻辑。
	if err := d.db.WithContext(ctx).First(&item, id).Error; err != nil {
		// 返回当前处理结果。
		return nil, err
	}
	// 返回当前处理结果。
	return &item, nil
}

// CreateLotteryInfoTx 在事务中创建图纸。
func (d *BizConfigDAO) CreateLotteryInfoTx(tx *gorm.DB, item *models.WLotteryInfo) error {
	// 返回当前处理结果。
	return tx.Create(item).Error
}

// UpdateLotteryInfoTx 在事务中更新图纸。
func (d *BizConfigDAO) UpdateLotteryInfoTx(tx *gorm.DB, id uint, updates map[string]interface{}) error {
	// 返回当前处理结果。
	return tx.Model(&models.WLotteryInfo{}).Where("id = ?", id).Updates(updates).Error
}

// DeleteLotteryInfoTx 在事务中删除图纸。
func (d *BizConfigDAO) DeleteLotteryInfoTx(tx *gorm.DB, id uint) error {
	// 返回当前处理结果。
	return tx.Delete(&models.WLotteryInfo{}, id).Error
}

// ListLotteryOptionsByInfoIDs 批量查询图纸的动物竞猜选项。
func (d *BizConfigDAO) ListLotteryOptionsByInfoIDs(ctx context.Context, infoIDs []uint) ([]models.WLotteryOption, error) {
	// 判断条件并进入对应分支逻辑。
	if len(infoIDs) == 0 {
		// 返回当前处理结果。
		return []models.WLotteryOption{}, nil
	}
	// 定义并初始化当前变量。
	options := make([]models.WLotteryOption, 0)
	// 判断条件并进入对应分支逻辑。
	if err := d.db.WithContext(ctx).
		// 调用Where完成当前处理。
		Where("lottery_info_id IN ?", infoIDs).
		// 调用Order完成当前处理。
		Order("sort ASC, id ASC").
		// 调用Find完成当前处理。
		Find(&options).Error; err != nil {
		// 返回当前处理结果。
		return nil, err
	}
	// 返回当前处理结果。
	return options, nil
}

// ReplaceLotteryOptionsTx 全量替换图纸动物竞猜选项（保留旧票数）。
func (d *BizConfigDAO) ReplaceLotteryOptionsTx(tx *gorm.DB, infoID uint, optionNames []string) error {
	// 读取旧选项票数，保证编辑时票数不丢失。
	var oldRows []models.WLotteryOption
	// 判断条件并进入对应分支逻辑。
	if err := tx.Where("lottery_info_id = ?", infoID).Find(&oldRows).Error; err != nil {
		// 返回当前处理结果。
		return err
	}
	// 定义并初始化当前变量。
	oldVotes := make(map[string]int64, len(oldRows))
	// 循环处理当前数据集合。
	for _, row := range oldRows {
		// 判断条件并进入对应分支逻辑。
		if row.OptionName == "" {
			// 处理当前语句逻辑。
			continue
		}
		// 更新当前变量或字段值。
		oldVotes[row.OptionName] = row.Votes
	}

	// 删除旧选项，避免残留。
	if err := tx.Where("lottery_info_id = ?", infoID).Delete(&models.WLotteryOption{}).Error; err != nil {
		// 返回当前处理结果。
		return err
	}

	// 构造新选项。
	rows := make([]models.WLotteryOption, 0, len(optionNames))
	// 循环处理当前数据集合。
	for idx, name := range optionNames {
		// 更新当前变量或字段值。
		rows = append(rows, models.WLotteryOption{
			// 处理当前语句逻辑。
			LotteryInfoID: infoID,
			// 处理当前语句逻辑。
			OptionName: name,
			// 处理当前语句逻辑。
			Votes: oldVotes[name],
			// 处理当前语句逻辑。
			Sort: idx + 1,
		})
	}
	// 判断条件并进入对应分支逻辑。
	if len(rows) == 0 {
		// 返回当前处理结果。
		return nil
	}
	// 返回当前处理结果。
	return tx.Create(&rows).Error
}
