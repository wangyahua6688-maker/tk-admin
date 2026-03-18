package rbac

import (
	"context"
	rbacdao "go-admin-full/internal/dao/rbac"
	"go-admin-full/internal/models"
)

// PermissionService 提供权限业务逻辑封装。
type PermissionService struct {
	dao rbacdao.PermissionDAO // 权限 DAO
}

// NewPermissionService 创建权限服务。
func NewPermissionService(dao rbacdao.PermissionDAO) *PermissionService {
	// 返回当前处理结果。
	return &PermissionService{dao: dao}
}

// Create 新增权限。
func (s *PermissionService) Create(ctx context.Context, p *models.Permission) error {
	// 返回当前处理结果。
	return s.dao.Create(ctx, p)
}

// Update 更新权限。
func (s *PermissionService) Update(ctx context.Context, p *models.Permission) error {
	// 返回当前处理结果。
	return s.dao.Update(ctx, p)
}

// List 查询权限列表。
func (s *PermissionService) List(ctx context.Context) ([]models.Permission, error) {
	// 返回当前处理结果。
	return s.dao.List(ctx)
}

// Get 获取单个权限。
func (s *PermissionService) Get(ctx context.Context, id uint) (*models.Permission, error) {
	// 返回当前处理结果。
	return s.dao.Get(ctx, id)
}

// Delete 删除权限。
func (s *PermissionService) Delete(ctx context.Context, id uint) error {
	// 返回当前处理结果。
	return s.dao.Delete(ctx, id)
}
