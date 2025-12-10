package services

import (
	"context"
	"errors"
	"go-admin-full/internal/dao"
	"go-admin-full/internal/models"
)

type RolePermissionService struct {
	dao *dao.RolePermissionDao
}

func NewRolePermissionService(d *dao.RolePermissionDao) *RolePermissionService {
	return &RolePermissionService{dao: d}
}

func (s *RolePermissionService) BindPermissions(ctx context.Context, roleID uint, permIDs []uint) error {
	role, err := s.dao.FindRole(ctx, roleID)
	if err != nil {
		return err
	}

	// 允许传空数组：表示清空角色已有权限绑定。
	if len(permIDs) == 0 {
		return s.dao.ReplacePermissions(ctx, role, []models.Permission{})
	}

	perms, err := s.dao.FindPermissions(ctx, permIDs)
	if err != nil {
		return err
	}
	if len(perms) != len(uniquePermissionIDs(permIDs)) {
		return errors.New("存在无效权限ID")
	}

	return s.dao.ReplacePermissions(ctx, role, perms)
}

func (s *RolePermissionService) GetRolePermissions(ctx context.Context, roleID uint) ([]models.Permission, error) {
	return s.dao.GetPermissions(ctx, roleID)
}

func uniquePermissionIDs(in []uint) []uint {
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
