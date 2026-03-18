package biz

import (
	"context"

	"go-admin-full/internal/models"
	"gorm.io/gorm"
)

// DrawRecordFilter 开奖区记录查询条件。
type DrawRecordFilter struct {
	SpecialLotteryID uint   // 彩种ID（0 表示不筛选）
	Keyword          string // 期号关键词
	Limit            int    // 限制条数
}

// ListDrawRecords 查询开奖记录列表。
func (s *BizConfigService) ListDrawRecords(ctx context.Context, filter DrawRecordFilter) ([]models.WDrawRecord, error) {
	// 交由 DAO 执行查询，保持 service 层专注于编排逻辑。
	return s.dao.ListDrawRecords(ctx, filter.SpecialLotteryID, filter.Keyword, filter.Limit)
}

// GetDrawRecordByID 获取单条开奖记录。
func (s *BizConfigService) GetDrawRecordByID(ctx context.Context, id uint) (*models.WDrawRecord, error) {
	// 单表查询直接下沉 DAO。
	return s.dao.GetDrawRecordByID(ctx, id)
}

// CreateDrawRecord 创建开奖记录并维护“唯一当前期”。
func (s *BizConfigService) CreateDrawRecord(ctx context.Context, item *models.WDrawRecord) error {
	// 通过事务保证“唯一当前期”与写入原子一致。
	return s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// 同彩种只允许一条 current 记录。
		if item.IsCurrent == 1 {
			// 先清除同彩种其它记录的 current 标识。
			if err := s.dao.ResetCurrentDrawRecordTx(tx, item.SpecialLotteryID, 0); err != nil {
				// 返回当前处理结果。
				return err
			}
		}
		// 写入开奖记录。
		return s.dao.CreateDrawRecordTx(tx, item)
	})
}

// UpdateDrawRecord 更新开奖记录并维护“唯一当前期”。
func (s *BizConfigService) UpdateDrawRecord(ctx context.Context, id uint, updates map[string]interface{}, specialLotteryID uint, isCurrent int8) error {
	// 通过事务保证 current 标识与更新一致。
	return s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// current 置为 1 时，先清理同彩种其它记录。
		if isCurrent == 1 && specialLotteryID > 0 {
			// 排除当前记录，避免将自身清零。
			if err := s.dao.ResetCurrentDrawRecordTx(tx, specialLotteryID, id); err != nil {
				// 返回当前处理结果。
				return err
			}
		}
		// 执行更新。
		return s.dao.UpdateDrawRecordTx(tx, id, updates)
	})
}

// DeleteDrawRecord 删除开奖记录。
func (s *BizConfigService) DeleteDrawRecord(ctx context.Context, id uint) error {
	// 删除记录无需事务编排，直接下沉 DAO。
	return s.dao.DeleteDrawRecord(ctx, id)
}
