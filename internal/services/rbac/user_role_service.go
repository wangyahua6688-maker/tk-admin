package rbac

import (
	"context"
	"errors"
	rbacdao "go-admin-full/internal/dao/rbac"
	"go-admin-full/internal/models"
)

// UserRoleService 用户角色关联业务。
type UserRoleService struct {
	dao *rbacdao.UserRoleDao // 用户角色 DAO
}

// NewUserRoleService 创建用户角色服务。
func NewUserRoleService(d *rbacdao.UserRoleDao) *UserRoleService {
	return &UserRoleService{dao: d}
}

// BindRoles 全量替换用户角色。
func (s *UserRoleService) BindRoles(ctx context.Context, userID uint, roleIDs []uint) error {
	// 校验用户是否存在
	user, err := s.dao.FindUser(ctx, userID)
	if err != nil {
		return err
	}

	// 允许传空数组：表示清空用户角色绑定。
	if len(roleIDs) == 0 {
		return s.dao.ReplaceRoles(ctx, user, []models.Role{})
	}

	// 查询角色集合
	roles, err := s.dao.FindRoles(ctx, roleIDs)
	if err != nil {
		return err
	}
	// 校验角色数量是否匹配，防止无效 ID 混入
	if len(roles) != len(uniqueUint(roleIDs)) {
		return errors.New("存在无效角色ID")
	}

	// 执行绑定
	return s.dao.ReplaceRoles(ctx, user, roles)
}

// AddRoles 追加用户角色。
func (s *UserRoleService) AddRoles(ctx context.Context, userID uint, roleIDs []uint) error {
	// 角色列表为空直接返回错误
	if len(roleIDs) == 0 {
		return errors.New("role_ids不能为空")
	}

	// 校验用户是否存在
	user, err := s.dao.FindUser(ctx, userID)
	if err != nil {
		return err
	}

	// 查询角色集合
	roles, err := s.dao.FindRoles(ctx, roleIDs)
	if err != nil {
		return err
	}
	// 校验角色数量是否匹配
	if len(roles) != len(uniqueUint(roleIDs)) {
		return errors.New("存在无效角色ID")
	}

	// 追加角色
	return s.dao.AppendRoles(ctx, user, roles)
}

// RemoveRoles 移除用户角色。
func (s *UserRoleService) RemoveRoles(ctx context.Context, userID uint, roleIDs []uint) error {
	// 角色列表为空直接返回错误
	if len(roleIDs) == 0 {
		return errors.New("role_ids不能为空")
	}

	// 校验用户是否存在
	user, err := s.dao.FindUser(ctx, userID)
	if err != nil {
		return err
	}

	// 查询角色集合
	roles, err := s.dao.FindRoles(ctx, roleIDs)
	if err != nil {
		return err
	}
	// 校验角色数量是否匹配
	if len(roles) != len(uniqueUint(roleIDs)) {
		return errors.New("存在无效角色ID")
	}

	// 移除角色
	return s.dao.RemoveRoles(ctx, user, roles)
}

// GetUserRoles 查询用户角色并预加载权限。
func (s *UserRoleService) GetUserRoles(ctx context.Context, userID uint) ([]models.Role, error) {
	return s.dao.GetUserRolesWithPermissions(ctx, userID)
}

// uniqueUint 对 uint 切片去重。
func uniqueUint(in []uint) []uint {
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
