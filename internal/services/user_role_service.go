package services

import (
	"context"
	"errors"
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

	// 允许传空数组：表示清空用户角色绑定。
	if len(roleIDs) == 0 {
		return s.dao.ReplaceRoles(ctx, user, []models.Role{})
	}

	roles, err := s.dao.FindRoles(ctx, roleIDs)
	if err != nil {
		return err
	}
	if len(roles) != len(uniqueUint(roleIDs)) {
		return errors.New("存在无效角色ID")
	}

	return s.dao.ReplaceRoles(ctx, user, roles)
}

func (s *UserRoleService) AddRoles(ctx context.Context, userID uint, roleIDs []uint) error {
	if len(roleIDs) == 0 {
		return errors.New("role_ids不能为空")
	}

	user, err := s.dao.FindUser(ctx, userID)
	if err != nil {
		return err
	}

	roles, err := s.dao.FindRoles(ctx, roleIDs)
	if err != nil {
		return err
	}
	if len(roles) != len(uniqueUint(roleIDs)) {
		return errors.New("存在无效角色ID")
	}

	return s.dao.AppendRoles(ctx, user, roles)
}

func (s *UserRoleService) RemoveRoles(ctx context.Context, userID uint, roleIDs []uint) error {
	if len(roleIDs) == 0 {
		return errors.New("role_ids不能为空")
	}

	user, err := s.dao.FindUser(ctx, userID)
	if err != nil {
		return err
	}

	roles, err := s.dao.FindRoles(ctx, roleIDs)
	if err != nil {
		return err
	}
	if len(roles) != len(uniqueUint(roleIDs)) {
		return errors.New("存在无效角色ID")
	}

	return s.dao.RemoveRoles(ctx, user, roles)
}

func (s *UserRoleService) GetUserRoles(ctx context.Context, userID uint) ([]models.Role, error) {
	return s.dao.GetUserRolesWithPermissions(ctx, userID)
}

func uniqueUint(in []uint) []uint {
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
