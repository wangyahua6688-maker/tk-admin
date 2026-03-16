package biz

import (
	"context"
	"errors"
	"strings"

	"go-admin-full/internal/models"
	"gorm.io/gorm"
)

// ListLotteryInfosWithOptions 查询图纸列表并附带动物竞猜选项。
func (s *BizConfigService) ListLotteryInfosWithOptions(ctx context.Context, limit int) ([]models.WLotteryInfo, map[uint][]string, error) {
	// 1) 查询图纸列表。
	items, err := s.dao.ListLotteryInfos(ctx, limit)
	if err != nil {
		return nil, nil, err
	}

	// 2) 组装图纸ID。
	infoIDs := make([]uint, 0, len(items))
	for _, item := range items {
		infoIDs = append(infoIDs, item.ID)
	}

	// 3) 批量查询动物竞猜选项。
	options, err := s.dao.ListLotteryOptionsByInfoIDs(ctx, infoIDs)
	if err != nil {
		return nil, nil, err
	}

	// 4) 按图纸ID分组。
	optionNameMap := make(map[uint][]string, len(items))
	for _, opt := range options {
		name := strings.TrimSpace(opt.OptionName)
		if name == "" {
			continue
		}
		optionNameMap[opt.LotteryInfoID] = append(optionNameMap[opt.LotteryInfoID], name)
	}

	return items, optionNameMap, nil
}

// GetLotteryInfoByID 获取单条图纸内容。
func (s *BizConfigService) GetLotteryInfoByID(ctx context.Context, id uint) (*models.WLotteryInfo, error) {
	return s.dao.GetLotteryInfoByID(ctx, id)
}

// ResolveLotteryCategory 解析分类输入为标准分类ID+标签。
func (s *BizConfigService) ResolveLotteryCategory(ctx context.Context, categoryID *uint, categoryTag *string) (uint, string, error) {
	// 优先按 category_id 解析。
	if categoryID != nil && *categoryID > 0 {
		cat, err := s.dao.GetLotteryCategoryByID(ctx, *categoryID)
		if err != nil {
			return 0, "", errors.New("category_id not found")
		}
		return cat.ID, strings.TrimSpace(cat.CategoryKey), nil
	}

	// 兼容旧请求：仅传 category_tag 时，按 key/name 反查。
	tag := safeString(categoryTag)
	if tag != "" {
		cat, err := s.dao.GetLotteryCategoryByTag(ctx, tag)
		if err == nil && cat != nil {
			return cat.ID, strings.TrimSpace(cat.CategoryKey), nil
		}
	}

	return 0, "", errors.New("category_id required")
}

// CreateLotteryInfo 创建图纸并写入竞猜选项。
func (s *BizConfigService) CreateLotteryInfo(ctx context.Context, item *models.WLotteryInfo, optionNames []string) error {
	// 归一化选项列表，防止为空。
	normalized := normalizeOptionNames(optionNames)
	if len(normalized) == 0 {
		normalized = defaultAnimalOptionNames()
	}

	return s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// 维护“同彩种唯一当前期”。
		if item.IsCurrent == 1 && item.SpecialLotteryID > 0 {
			if err := tx.Model(&models.WLotteryInfo{}).
				Where("special_lottery_id = ?", item.SpecialLotteryID).
				Update("is_current", 0).Error; err != nil {
				return err
			}
		}
		// 写入图纸。
		if err := s.dao.CreateLotteryInfoTx(tx, item); err != nil {
			return err
		}
		// 写入动物竞猜选项。
		return s.dao.ReplaceLotteryOptionsTx(tx, item.ID, normalized)
	})
}

// UpdateLotteryInfo 更新图纸并在必要时替换竞猜选项。
func (s *BizConfigService) UpdateLotteryInfo(ctx context.Context, id uint, updates map[string]interface{}, updateOptions bool, optionNames []string, specialLotteryID uint, isCurrent int8) error {
	// 仅在需要时更新选项。
	normalized := normalizeOptionNames(optionNames)
	if updateOptions && len(normalized) == 0 {
		normalized = defaultAnimalOptionNames()
	}

	return s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// 维护“同彩种唯一当前期”。
		if isCurrent == 1 && specialLotteryID > 0 {
			if err := tx.Model(&models.WLotteryInfo{}).
				Where("special_lottery_id = ? AND id <> ?", specialLotteryID, id).
				Update("is_current", 0).Error; err != nil {
				return err
			}
		}
		// 更新图纸。
		if err := s.dao.UpdateLotteryInfoTx(tx, id, updates); err != nil {
			return err
		}
		// 可选：全量替换动物竞猜选项。
		if updateOptions {
			return s.dao.ReplaceLotteryOptionsTx(tx, id, normalized)
		}
		return nil
	})
}

// DeleteLotteryInfo 删除图纸并清理竞猜选项。
func (s *BizConfigService) DeleteLotteryInfo(ctx context.Context, id uint) error {
	return s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// 先删选项，避免残留。
		if err := tx.Where("lottery_info_id = ?", id).Delete(&models.WLotteryOption{}).Error; err != nil {
			return err
		}
		// 再删图纸。
		return s.dao.DeleteLotteryInfoTx(tx, id)
	})
}

// normalizeOptionNames 去重清洗动物竞猜选项，保持输入顺序。
func normalizeOptionNames(input []string) []string {
	out := make([]string, 0, len(input))
	seen := make(map[string]struct{}, len(input))
	for _, raw := range input {
		name := strings.TrimSpace(raw)
		if name == "" {
			continue
		}
		if _, ok := seen[name]; ok {
			continue
		}
		seen[name] = struct{}{}
		out = append(out, name)
	}
	return out
}

// defaultAnimalOptionNames 返回默认 12 生肖竞猜选项。
func defaultAnimalOptionNames() []string {
	return []string{"鼠", "牛", "虎", "兔", "龙", "蛇", "马", "羊", "猴", "鸡", "狗", "猪"}
}

// safeString 读取可空字符串指针并做空格裁剪。
func safeString(v *string) string {
	if v == nil {
		return ""
	}
	return strings.TrimSpace(*v)
}
