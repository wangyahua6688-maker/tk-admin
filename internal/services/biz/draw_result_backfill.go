package biz

import (
	"context"
	"fmt"

	"go-admin/internal/models"

	"gorm.io/gorm"
)

// BackfillDrawResultData 为历史开奖记录补齐主表派生字段和玩法结果分表。
func (s *LotteryService) BackfillDrawResultData(ctx context.Context, specialLotteryID uint, force bool) (int, error) {
	rows := make([]models.WDrawRecord, 0)
	query := s.db.WithContext(ctx).Model(&models.WDrawRecord{})
	if specialLotteryID > 0 {
		query = query.Where("special_lottery_id = ?", specialLotteryID)
	}
	if !force {
		query = query.Where(
			`draw_labels = '' OR color_labels = '' OR zodiac_labels = '' OR wuxing_labels = '' OR
			 special_single_double = '' OR special_big_small = '' OR sum_single_double = '' OR sum_big_small = '' OR
			 special_code = '' OR normal_code = '' OR zheng1 = '' OR zheng2 = '' OR zheng3 = '' OR
			 zheng4 = '' OR zheng5 = '' OR zheng6 = ''`,
		)
	}
	if err := query.Order("id ASC").Find(&rows).Error; err != nil {
		return 0, err
	}

	count := 0
	for idx := range rows {
		current := rows[idx]
		next := current
		bundle, err := hydrateDrawRecordDerivedFields(&next)
		if err != nil {
			return count, fmt.Errorf("compile draw record %d (%s) failed: %w", current.ID, current.Issue, err)
		}
		if err := s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
			if err := s.dao.UpdateDrawRecordTx(tx, current.ID, buildDrawRecordUpdateMap(&next)); err != nil {
				return err
			}
			return upsertDrawResultTablesTx(tx, current.ID, bundle)
		}); err != nil {
			return count, fmt.Errorf("backfill draw record %d (%s) failed: %w", current.ID, current.Issue, err)
		}
		count++
	}

	return count, nil
}
