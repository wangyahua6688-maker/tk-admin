package rbac

import (
	"context"
	"strings"

	"go-admin-full/internal/models"
	"gorm.io/gorm"
)

// LoginLogDao 登录日志数据访问层。
type LoginLogDao struct {
	db *gorm.DB // 数据库连接
}

// NewLoginLogDao 创建登录日志 DAO。
func NewLoginLogDao(db *gorm.DB) *LoginLogDao {
	return &LoginLogDao{db: db}
}

// Create 写入一条登录日志。
func (d *LoginLogDao) Create(ctx context.Context, log *models.LoginLog) error {
	return d.db.WithContext(ctx).Create(log).Error
}

// List 分页查询登录日志，可按用户名精确筛选。
func (d *LoginLogDao) List(ctx context.Context, page, pageSize int, username string) ([]models.LoginLog, int64, error) {
	if page < 1 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 20
	}
	if pageSize > 100 {
		pageSize = 100
	}

	query := d.db.WithContext(ctx).Model(&models.LoginLog{}) // 构造查询
	if strings.TrimSpace(username) != "" {
		query = query.Where("username = ?", strings.TrimSpace(username)) // 按用户名过滤
	}

	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	var logs []models.LoginLog
	if err := query.
		// 按时间倒序展示
		Order("id DESC").
		Offset((page - 1) * pageSize).
		Limit(pageSize).
		Find(&logs).Error; err != nil {
		return nil, 0, err
	}

	return logs, total, nil
}
