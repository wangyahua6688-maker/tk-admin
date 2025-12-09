package services

import (
	"context"
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

	perms, err := s.dao.FindPermissions(ctx, permIDs)
	if err != nil {
		return err
	}

	return s.dao.ReplacePermissions(ctx, role, perms)
}

func (s *RolePermissionService) GetRolePermissions(ctx context.Context, roleID uint) ([]models.Permission, error) {
	return s.dao.GetPermissions(ctx, roleID)
}
