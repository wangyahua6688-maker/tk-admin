package rbac

import (
	"context"
	"strings"

	"go-admin/internal/models"
	"gorm.io/gorm"
)

// LoginLogDao 登录日志数据访问层。
type LoginLogDao struct {
	db *gorm.DB // 数据库连接
}

// NewLoginLogDao 创建登录日志 DAO。
func NewLoginLogDao(db *gorm.DB) *LoginLogDao {
	// 返回当前处理结果。
	return &LoginLogDao{db: db}
}

// Create 写入一条登录日志。
func (d *LoginLogDao) Create(ctx context.Context, log *models.LoginLog) error {
	// 返回当前处理结果。
	return d.db.WithContext(ctx).Create(log).Error
}

// List 分页查询登录日志，可按用户名精确筛选。
func (d *LoginLogDao) List(ctx context.Context, page, pageSize int, username string) ([]models.LoginLog, int64, error) {
	// 判断条件并进入对应分支逻辑。
	if page < 1 {
		// 更新当前变量或字段值。
		page = 1
	}
	// 判断条件并进入对应分支逻辑。
	if pageSize <= 0 {
		// 更新当前变量或字段值。
		pageSize = 20
	}
	// 判断条件并进入对应分支逻辑。
	if pageSize > 100 {
		// 更新当前变量或字段值。
		pageSize = 100
	}

	query := d.db.WithContext(ctx).Model(&models.LoginLog{}) // 构造查询
	// 判断条件并进入对应分支逻辑。
	if strings.TrimSpace(username) != "" {
		query = query.Where("username = ?", strings.TrimSpace(username)) // 按用户名过滤
	}

	// 声明当前变量。
	var total int64
	// 判断条件并进入对应分支逻辑。
	if err := query.Count(&total).Error; err != nil {
		// 返回当前处理结果。
		return nil, 0, err
	}

	// 声明当前变量。
	var logs []models.LoginLog
	// 判断条件并进入对应分支逻辑。
	if err := query.
		// 按时间倒序展示
		Order("id DESC").
		// 调用Offset完成当前处理。
		Offset((page - 1) * pageSize).
		// 调用Limit完成当前处理。
		Limit(pageSize).
		// 调用Find完成当前处理。
		Find(&logs).Error; err != nil {
		// 返回当前处理结果。
		return nil, 0, err
	}

	// 返回当前处理结果。
	return logs, total, nil
}
