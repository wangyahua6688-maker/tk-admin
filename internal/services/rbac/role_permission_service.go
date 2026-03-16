package rbac

import (
	"context"
	"errors"
	rbacdao "go-admin-full/internal/dao/rbac"
	"go-admin-full/internal/models"
)

// RolePermissionService 角色权限关联业务。
type RolePermissionService struct {
	dao *rbacdao.RolePermissionDao // 角色权限 DAO
}

// NewRolePermissionService 创建角色权限服务。
func NewRolePermissionService(d *rbacdao.RolePermissionDao) *RolePermissionService {
	return &RolePermissionService{dao: d}
}

// BindPermissions 全量替换角色权限。
func (s *RolePermissionService) BindPermissions(ctx context.Context, roleID uint, permIDs []uint) error {
	// 查询角色是否存在
	role, err := s.dao.FindRole(ctx, roleID)
	if err != nil {
		return err
	}

	// 允许传空数组：表示清空角色已有权限绑定。
	if len(permIDs) == 0 {
		return s.dao.ReplacePermissions(ctx, role, []models.Permission{})
	}

	// 查询权限集合
	perms, err := s.dao.FindPermissions(ctx, permIDs)
	if err != nil {
		return err
	}
	// 校验权限数量是否匹配，防止无效 ID 混入
	if len(perms) != len(uniquePermissionIDs(permIDs)) {
		return errors.New("存在无效权限ID")
	}

	// 执行权限绑定
	return s.dao.ReplacePermissions(ctx, role, perms)
}

// GetRolePermissions 查询角色已绑定权限。
func (s *RolePermissionService) GetRolePermissions(ctx context.Context, roleID uint) ([]models.Permission, error) {
	return s.dao.GetPermissions(ctx, roleID)
}

// uniquePermissionIDs 对权限 ID 去重。
func uniquePermissionIDs(in []uint) []uint {
	// 使用 map 去重
	set := make(map[uint]struct{}, len(in))
	out := make([]uint, 0, len(in))
	for _, v := range in {
		if _, ok := set[v]; ok {
			continue
		}
		set[v] = struct{}{}
		out = append(out, v)
	}
	return out
}
