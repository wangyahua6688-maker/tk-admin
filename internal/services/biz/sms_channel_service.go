package biz

import (
	"context"
	"strconv"
	"strings"

	bizdao "go-admin/internal/dao/biz"
	"go-admin/internal/models"
	"gorm.io/gorm"
)

// SMSChannelService 短信通道业务服务。
type SMSChannelService struct {
	dao *bizdao.SMSChannelDAO
}

// NewSMSChannelService 创建短信通道服务。
func NewSMSChannelService(db *gorm.DB) *SMSChannelService {
	return &SMSChannelService{dao: bizdao.NewSMSChannelDAO(db)}
}

// ListSMSChannels 查询短信通道配置列表。
func (s *SMSChannelService) ListSMSChannels(ctx context.Context, status string, limit int) ([]models.WSMSChannel, error) {
	// 将字符串状态转换为可选筛选值。
	var statusFilter *int
	// 定义并初始化当前变量。
	trimmed := strings.TrimSpace(status)
	// 判断条件并进入对应分支逻辑。
	if trimmed == "0" || trimmed == "1" {
		// 定义并初始化当前变量。
		v, _ := strconv.Atoi(trimmed)
		// 更新当前变量或字段值。
		statusFilter = &v
	}
	// 调用 DAO 执行查询。
	return s.dao.ListSMSChannels(ctx, statusFilter, limit)
}

// CreateSMSChannel 新增短信通道配置。
func (s *SMSChannelService) CreateSMSChannel(ctx context.Context, item *models.WSMSChannel) error {
	// 直通 DAO 创建记录。
	return s.dao.CreateSMSChannel(ctx, item)
}

// UpdateSMSChannel 更新短信通道配置。
func (s *SMSChannelService) UpdateSMSChannel(ctx context.Context, id uint, updates map[string]interface{}) error {
	// 直通 DAO 更新记录。
	return s.dao.UpdateSMSChannel(ctx, id, updates)
}

// DeleteSMSChannel 删除短信通道配置。
func (s *SMSChannelService) DeleteSMSChannel(ctx context.Context, id uint) error {
	// 直通 DAO 删除记录。
	return s.dao.DeleteSMSChannel(ctx, id)
}
