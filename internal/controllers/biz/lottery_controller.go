package biz

import (
	"context"

	"go-admin/internal/models"
	bizsvc "go-admin/internal/services/biz"

	"gorm.io/gorm"
)

// specialLotteryService 定义彩种配置依赖。
type specialLotteryService interface {
	ListSpecialLotteries(ctx context.Context, limit int) ([]models.WSpecialLottery, error)
	CreateSpecialLottery(ctx context.Context, item *models.WSpecialLottery) error
	UpdateSpecialLottery(ctx context.Context, id uint, updates map[string]interface{}) error
	DeleteSpecialLottery(ctx context.Context, id uint) error
}

// lotteryInfoService 定义图纸配置依赖。
type lotteryInfoService interface {
	ListLotteryInfosWithOptions(ctx context.Context, limit int) ([]models.WLotteryInfo, map[uint][]string, error)
	GetLotteryInfoByID(ctx context.Context, id uint) (*models.WLotteryInfo, error)
	ResolveLotteryCategory(ctx context.Context, categoryID *uint, categoryTag *string) (uint, string, error)
	CreateLotteryInfo(ctx context.Context, item *models.WLotteryInfo, optionNames []string) error
	UpdateLotteryInfo(ctx context.Context, id uint, updates map[string]interface{}, updateOptions bool, optionNames []string, specialLotteryID uint, isCurrent int8) error
	DeleteLotteryInfo(ctx context.Context, id uint) error
}

// drawRecordService 定义开奖记录依赖。
type drawRecordService interface {
	ListDrawRecords(ctx context.Context, filter bizsvc.DrawRecordFilter) ([]models.WDrawRecord, error)
	GetDrawRecordByID(ctx context.Context, id uint) (*models.WDrawRecord, error)
	CreateDrawRecord(ctx context.Context, item *models.WDrawRecord) error
	UpdateDrawRecord(ctx context.Context, id uint, item *models.WDrawRecord) error
	DeleteDrawRecord(ctx context.Context, id uint) error
}

// lotteryCategoryService 定义图库分类依赖。
type lotteryCategoryService interface {
	ListLotteryCategories(ctx context.Context, keyword string) ([]models.WLotteryCategory, error)
	CreateLotteryCategory(ctx context.Context, item *models.WLotteryCategory) error
	UpdateLotteryCategory(ctx context.Context, id uint, updates map[string]interface{}) error
	DeleteLotteryCategory(ctx context.Context, id uint) error
}

// LotteryController 彩票业务控制器。
// 该控制器只负责彩票域：彩种、图库、开奖记录、图库分类。
type LotteryController struct {
	specialLotterySvc  specialLotteryService
	lotteryInfoSvc     lotteryInfoService
	drawRecordSvc      drawRecordService
	lotteryCategorySvc lotteryCategoryService
}

// LotteryControllerDeps 定义彩票控制器依赖。
type LotteryControllerDeps struct {
	SpecialLotteryService  specialLotteryService
	LotteryInfoService     lotteryInfoService
	DrawRecordService      drawRecordService
	LotteryCategoryService lotteryCategoryService
}

// NewLotteryController 创建彩票控制器。
func NewLotteryController(db *gorm.DB) *LotteryController {
	service := bizsvc.NewLotteryService(db)
	return NewLotteryControllerWithDeps(LotteryControllerDeps{
		SpecialLotteryService:  service,
		LotteryInfoService:     service,
		DrawRecordService:      service,
		LotteryCategoryService: service,
	})
}

// NewLotteryControllerWithDeps 使用显式依赖创建彩票控制器。
func NewLotteryControllerWithDeps(deps LotteryControllerDeps) *LotteryController {
	return &LotteryController{
		specialLotterySvc:  deps.SpecialLotteryService,
		lotteryInfoSvc:     deps.LotteryInfoService,
		drawRecordSvc:      deps.DrawRecordService,
		lotteryCategorySvc: deps.LotteryCategoryService,
	}
}
