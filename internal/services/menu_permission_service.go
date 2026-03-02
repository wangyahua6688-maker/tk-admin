package services

import (
	"context"
	"errors"

	"go-admin-full/internal/dao"
	"go-admin-full/internal/models"
)

// MenuPermissionService 菜单权限关联业务。
type MenuPermissionService struct {
	dao *dao.MenuPermissionDao
}

func NewMenuPermissionService(d *dao.MenuPermissionDao) *MenuPermissionService {
	return &MenuPermissionService{dao: d}
}

// BindPermissions 全量替换菜单权限。
func (s *MenuPermissionService) BindPermissions(ctx context.Context, menuID uint, permIDs []uint) error {
	menu, err := s.dao.FindMenu(ctx, menuID)
	if err != nil {
		return err
	}

	// 允许传空数组：表示清空菜单已有权限绑定。
	if len(permIDs) == 0 {
		return s.dao.ReplacePermissions(ctx, menu, []models.Permission{})
	}

	perms, err := s.dao.FindPermissions(ctx, permIDs)
	if err != nil {
		return err
	}
	if len(perms) != len(uniqueMenuPermissionIDs(permIDs)) {
		return errors.New("存在无效权限ID")
	}

	return s.dao.ReplacePermissions(ctx, menu, perms)
}

// GetMenuPermissions 查询菜单已绑定权限。
func (s *MenuPermissionService) GetMenuPermissions(ctx context.Context, menuID uint) ([]models.Permission, error) {
	return s.dao.GetPermissions(ctx, menuID)
}

func uniqueMenuPermissionIDs(in []uint) []uint {
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
