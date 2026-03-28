package biz

import (
	"context"
	"errors"
	"strings"

	"go-admin/internal/models"
	"gorm.io/gorm"
)

// ListLotteryInfosWithOptions 查询图纸列表并附带动物竞猜选项。
func (s *LotteryService) ListLotteryInfosWithOptions(ctx context.Context, limit int) ([]models.WLotteryInfo, map[uint][]string, error) {
	// 1) 查询图纸列表。
	items, err := s.dao.ListLotteryInfos(ctx, limit)
	// 判断条件并进入对应分支逻辑。
	if err != nil {
		// 返回当前处理结果。
		return nil, nil, err
	}

	// 2) 组装图纸ID。
	infoIDs := make([]uint, 0, len(items))
	// 循环处理当前数据集合。
	for _, item := range items {
		// 更新当前变量或字段值。
		infoIDs = append(infoIDs, item.ID)
	}

	// 3) 批量查询动物竞猜选项。
	options, err := s.dao.ListLotteryOptionsByInfoIDs(ctx, infoIDs)
	// 判断条件并进入对应分支逻辑。
	if err != nil {
		// 返回当前处理结果。
		return nil, nil, err
	}

	// 4) 按图纸ID分组。
	optionNameMap := make(map[uint][]string, len(items))
	// 循环处理当前数据集合。
	for _, opt := range options {
		// 定义并初始化当前变量。
		name := strings.TrimSpace(opt.OptionName)
		// 判断条件并进入对应分支逻辑。
		if name == "" {
			// 处理当前语句逻辑。
			continue
		}
		// 更新当前变量或字段值。
		optionNameMap[opt.LotteryInfoID] = append(optionNameMap[opt.LotteryInfoID], name)
	}

	// 返回当前处理结果。
	return items, optionNameMap, nil
}

// GetLotteryInfoByID 获取单条图纸内容。
func (s *LotteryService) GetLotteryInfoByID(ctx context.Context, id uint) (*models.WLotteryInfo, error) {
	// 返回当前处理结果。
	return s.dao.GetLotteryInfoByID(ctx, id)
}

// ResolveLotteryCategory 解析分类输入为标准分类ID+标签。
func (s *LotteryService) ResolveLotteryCategory(ctx context.Context, categoryID *uint, categoryTag *string) (uint, string, error) {
	// 优先按 category_id 解析。
	if categoryID != nil && *categoryID > 0 {
		// 定义并初始化当前变量。
		cat, err := s.dao.GetLotteryCategoryByID(ctx, *categoryID)
		// 判断条件并进入对应分支逻辑。
		if err != nil {
			// 返回当前处理结果。
			return 0, "", errors.New("category_id not found")
		}
		// 返回当前处理结果。
		return cat.ID, strings.TrimSpace(cat.CategoryKey), nil
	}

	// 兼容旧请求：仅传 category_tag 时，按 key/name 反查。
	tag := safeString(categoryTag)
	// 判断条件并进入对应分支逻辑。
	if tag != "" {
		// 定义并初始化当前变量。
		cat, err := s.dao.GetLotteryCategoryByTag(ctx, tag)
		// 判断条件并进入对应分支逻辑。
		if err == nil && cat != nil {
			// 返回当前处理结果。
			return cat.ID, strings.TrimSpace(cat.CategoryKey), nil
		}
	}

	// 返回当前处理结果。
	return 0, "", errors.New("category_id required")
}

// CreateLotteryInfo 创建图纸并写入竞猜选项。
func (s *LotteryService) CreateLotteryInfo(ctx context.Context, item *models.WLotteryInfo, optionNames []string) error {
	// 归一化选项列表，防止为空。
	normalized := normalizeOptionNames(optionNames)
	// 判断条件并进入对应分支逻辑。
	if len(normalized) == 0 {
		// 更新当前变量或字段值。
		normalized = defaultAnimalOptionNames()
	}

	// 返回当前处理结果。
	return s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// 维护“同彩种唯一当前期”。
		if item.IsCurrent == 1 && item.SpecialLotteryID > 0 {
			// 判断条件并进入对应分支逻辑。
			if err := tx.Model(&models.WLotteryInfo{}).
				// 更新当前变量或字段值。
				Where("special_lottery_id = ?", item.SpecialLotteryID).
				// 调用Update完成当前处理。
				Update("is_current", 0).Error; err != nil {
				// 返回当前处理结果。
				return err
			}
		}
		// 写入图纸。
		if err := s.dao.CreateLotteryInfoTx(tx, item); err != nil {
			// 返回当前处理结果。
			return err
		}
		// 写入动物竞猜选项。
		return s.dao.ReplaceLotteryOptionsTx(tx, item.ID, normalized)
	})
}

// UpdateLotteryInfo 更新图纸并在必要时替换竞猜选项。
func (s *LotteryService) UpdateLotteryInfo(ctx context.Context, id uint, updates map[string]interface{}, updateOptions bool, optionNames []string, specialLotteryID uint, isCurrent int8) error {
	// 仅在需要时更新选项。
	normalized := normalizeOptionNames(optionNames)
	// 判断条件并进入对应分支逻辑。
	if updateOptions && len(normalized) == 0 {
		// 更新当前变量或字段值。
		normalized = defaultAnimalOptionNames()
	}

	// 返回当前处理结果。
	return s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// 维护“同彩种唯一当前期”。
		if isCurrent == 1 && specialLotteryID > 0 {
			// 判断条件并进入对应分支逻辑。
			if err := tx.Model(&models.WLotteryInfo{}).
				// 更新当前变量或字段值。
				Where("special_lottery_id = ? AND id <> ?", specialLotteryID, id).
				// 调用Update完成当前处理。
				Update("is_current", 0).Error; err != nil {
				// 返回当前处理结果。
				return err
			}
		}
		// 更新图纸。
		if err := s.dao.UpdateLotteryInfoTx(tx, id, updates); err != nil {
			// 返回当前处理结果。
			return err
		}
		// 可选：全量替换动物竞猜选项。
		if updateOptions {
			// 返回当前处理结果。
			return s.dao.ReplaceLotteryOptionsTx(tx, id, normalized)
		}
		// 返回当前处理结果。
		return nil
	})
}

// DeleteLotteryInfo 删除图纸并清理竞猜选项。
func (s *LotteryService) DeleteLotteryInfo(ctx context.Context, id uint) error {
	// 返回当前处理结果。
	return s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// 先删选项，避免残留。
		if err := tx.Where("lottery_info_id = ?", id).Delete(&models.WLotteryOption{}).Error; err != nil {
			// 返回当前处理结果。
			return err
		}
		// 再删图纸。
		return s.dao.DeleteLotteryInfoTx(tx, id)
	})
}

// normalizeOptionNames 去重清洗动物竞猜选项，保持输入顺序。
func normalizeOptionNames(input []string) []string {
	// 定义并初始化当前变量。
	out := make([]string, 0, len(input))
	// 定义并初始化当前变量。
	seen := make(map[string]struct{}, len(input))
	// 循环处理当前数据集合。
	for _, raw := range input {
		// 定义并初始化当前变量。
		name := strings.TrimSpace(raw)
		// 判断条件并进入对应分支逻辑。
		if name == "" {
			// 处理当前语句逻辑。
			continue
		}
		// 判断条件并进入对应分支逻辑。
		if _, ok := seen[name]; ok {
			// 处理当前语句逻辑。
			continue
		}
		// 更新当前变量或字段值。
		seen[name] = struct{}{}
		// 更新当前变量或字段值。
		out = append(out, name)
	}
	// 返回当前处理结果。
	return out
}

// defaultAnimalOptionNames 返回默认 12 生肖竞猜选项。
func defaultAnimalOptionNames() []string {
	// 返回当前处理结果。
	return []string{"鼠", "牛", "虎", "兔", "龙", "蛇", "马", "羊", "猴", "鸡", "狗", "猪"}
}

// safeString 读取可空字符串指针并做空格裁剪。
func safeString(v *string) string {
	// 判断条件并进入对应分支逻辑。
	if v == nil {
		// 返回当前处理结果。
		return ""
	}
	// 返回当前处理结果。
	return strings.TrimSpace(*v)
}
