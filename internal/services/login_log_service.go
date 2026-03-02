package services

import (
	"context"

	"go-admin-full/internal/dao"
	"go-admin-full/internal/models"
)

// LoginLogService 登录日志业务层。
type LoginLogService struct {
	dao *dao.LoginLogDao
}

func NewLoginLogService(d *dao.LoginLogDao) *LoginLogService {
	return &LoginLogService{dao: d}
}

// CreateLoginLog 写入登录日志。
func (s *LoginLogService) CreateLoginLog(ctx context.Context, log *models.LoginLog) error {
	return s.dao.Create(ctx, log)
}

// ListLoginLogs 分页查询登录日志。
func (s *LoginLogService) ListLoginLogs(ctx context.Context, page, pageSize int, username string) ([]models.LoginLog, int64, error) {
	return s.dao.List(ctx, page, pageSize, username)
}
