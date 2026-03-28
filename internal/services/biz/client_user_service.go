package biz

import (
	"context"

	bizdao "go-admin/internal/dao/biz"
	"go-admin/internal/models"
	"gorm.io/gorm"
)

// ClientUserService 客户端用户业务服务。
type ClientUserService struct {
	dao *bizdao.ClientUserDAO
}

// NewClientUserService 创建客户端用户服务。
func NewClientUserService(db *gorm.DB) *ClientUserService {
	return &ClientUserService{dao: bizdao.NewClientUserDAO(db)}
}

// ListClientUsers 查询客户端用户列表。
func (s *ClientUserService) ListClientUsers(ctx context.Context, userType string, limit int) ([]models.WUser, error) {
	// 直通 DAO 查询列表。
	return s.dao.ListClientUsers(ctx, userType, limit)
}

// CreateClientUser 新增客户端用户。
func (s *ClientUserService) CreateClientUser(ctx context.Context, item *models.WUser) error {
	// 直通 DAO 创建记录。
	return s.dao.CreateClientUser(ctx, item)
}

// UpdateClientUser 更新客户端用户。
func (s *ClientUserService) UpdateClientUser(ctx context.Context, id uint, updates map[string]interface{}) error {
	// 直通 DAO 更新记录。
	return s.dao.UpdateClientUser(ctx, id, updates)
}

// DeleteClientUser 删除客户端用户。
func (s *ClientUserService) DeleteClientUser(ctx context.Context, id uint) error {
	// 直通 DAO 删除记录。
	return s.dao.DeleteClientUser(ctx, id)
}
