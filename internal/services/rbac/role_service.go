package rbac

import (
	"context"
	rbacdao "go-admin-full/internal/dao/rbac"
	"go-admin-full/internal/models"
)

// RoleService 提供角色业务逻辑封装。
type RoleService struct {
	dao rbacdao.RoleDAO // 角色 DAO
}

// NewRoleService 创建角色服务。
func NewRoleService(dao rbacdao.RoleDAO) *RoleService {
	return &RoleService{dao: dao}
}

// Create 新增角色。
func (s *RoleService) Create(ctx context.Context, r *models.Role) error {
	return s.dao.Create(ctx, r)
}

// Update 更新角色。
func (s *RoleService) Update(ctx context.Context, r *models.Role) error {
	return s.dao.Update(ctx, r)
}

// List 查询角色列表。
func (s *RoleService) List(ctx context.Context) ([]models.Role, error) {
	return s.dao.List(ctx)
}

// Get 获取单个角色。
func (s *RoleService) Get(ctx context.Context, id uint) (*models.Role, error) {
	return s.dao.Get(ctx, id)
}

// Delete 删除角色。
func (s *RoleService) Delete(ctx context.Context, id uint) error {
	return s.dao.Delete(ctx, id)
}

// GetRolesByUserID 根据用户 ID 查询角色集合。
func (s *RoleService) GetRolesByUserID(ctx context.Context, userID uint) ([]models.Role, error) {
	return s.dao.GetByUserID(ctx, userID)
}
