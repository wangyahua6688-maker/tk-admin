package services

import (
	"context"
	"go-admin-full/internal/dao"
	"go-admin-full/internal/models"
)

type UserRoleService struct {
	dao *dao.UserRoleDao
}

func NewUserRoleService(d *dao.UserRoleDao) *UserRoleService {
	return &UserRoleService{dao: d}
}

func (s *UserRoleService) BindRoles(ctx context.Context, userID uint, roleIDs []uint) error {
	user, err := s.dao.FindUser(ctx, userID)
	if err != nil {
		return err
	}

	roles, err := s.dao.FindRoles(ctx, roleIDs)
	if err != nil {
		return err
	}

	return s.dao.ReplaceRoles(ctx, user, roles)
}

func (s *UserRoleService) AddRoles(ctx context.Context, userID uint, roleIDs []uint) error {
	user, err := s.dao.FindUser(ctx, userID)
	if err != nil {
		return err
	}

	roles, err := s.dao.FindRoles(ctx, roleIDs)
	if err != nil {
		return err
	}

	return s.dao.AppendRoles(ctx, user, roles)
}

func (s *UserRoleService) RemoveRoles(ctx context.Context, userID uint, roleIDs []uint) error {
	user, err := s.dao.FindUser(ctx, userID)
	if err != nil {
		return err
	}

	roles, err := s.dao.FindRoles(ctx, roleIDs)
	if err != nil {
		return err
	}

	return s.dao.RemoveRoles(ctx, user, roles)
}

func (s *UserRoleService) GetUserRoles(ctx context.Context, userID uint) ([]models.Role, error) {
	return s.dao.GetUserRolesWithPermissions(ctx, userID)
}
