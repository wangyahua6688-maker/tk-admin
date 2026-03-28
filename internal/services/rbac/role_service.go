package rbac

import (
	"context"
	"errors"
	rbacdao "go-admin/internal/dao/rbac"
	"go-admin/internal/models"
)

// RoleService 提供角色业务逻辑封装。
type RoleService struct {
	dao rbacdao.RoleDAO // 角色 DAO
}

// NewRoleService 创建角色服务。
func NewRoleService(dao rbacdao.RoleDAO) *RoleService {
	// 返回当前处理结果。
	return &RoleService{dao: dao}
}

// Create 新增角色。
func (s *RoleService) Create(ctx context.Context, r *models.Role) error {
	// 返回当前处理结果。
	return s.dao.Create(ctx, r)
}

// Update 更新角色。
func (s *RoleService) Update(ctx context.Context, r *models.Role) error {
	// 返回当前处理结果。
	return s.dao.Update(ctx, r)
}

// List 查询角色列表。
func (s *RoleService) List(ctx context.Context) ([]models.Role, error) {
	// 返回当前处理结果。
	return s.dao.List(ctx)
}

// Get 获取单个角色。
func (s *RoleService) Get(ctx context.Context, id uint) (*models.Role, error) {
	role, err := s.dao.Get(ctx, id)
	if err != nil {
		return nil, err
	}
	if role == nil {
		return nil, errors.New("角色不存在")
	}
	return role, nil
}

// Delete 删除角色。
func (s *RoleService) Delete(ctx context.Context, id uint) error {
	// 返回当前处理结果。
	return s.dao.Delete(ctx, id)
}
