package biz

import (
	"context"

	"go-admin-full/internal/models"
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

// CreateDrawRecord/UpdateDrawRecord/DeleteDrawRecord 已迁移到 draw_result_compiler.go。
// 说明：
// 1. 开奖号码保存时需要同时编译玩法结果并写入多张结果分表；
// 2. 这些方法必须和规则编译器放在一起，避免 controller/service 各自维护一套逻辑。
