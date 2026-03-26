package rbac

import (
	"context"
	"errors"

	rbacdao "go-admin/internal/dao/rbac"
	"go-admin/internal/models"
)

// MenuPermissionService 菜单权限关联业务。
type MenuPermissionService struct {
	dao *rbacdao.MenuPermissionDao // 菜单权限 DAO
}

// NewMenuPermissionService 创建菜单权限服务。
func NewMenuPermissionService(d *rbacdao.MenuPermissionDao) *MenuPermissionService {
	// 返回当前处理结果。
	return &MenuPermissionService{dao: d}
}

// BindPermissions 全量替换菜单权限。
func (s *MenuPermissionService) BindPermissions(ctx context.Context, menuID uint, permIDs []uint) error {
	// 查询菜单是否存在
	menu, err := s.dao.FindMenu(ctx, menuID)
	// 判断条件并进入对应分支逻辑。
	if err != nil {
		// 返回当前处理结果。
		return err
	}

	// 允许传空数组：表示清空菜单已有权限绑定。
	if len(permIDs) == 0 {
		// 返回当前处理结果。
		return s.dao.ReplacePermissions(ctx, menu, []models.Permission{})
	}

	// 查询权限集合
	perms, err := s.dao.FindPermissions(ctx, permIDs)
	// 判断条件并进入对应分支逻辑。
	if err != nil {
		// 返回当前处理结果。
		return err
	}
	// 校验权限数量是否匹配，防止无效 ID 混入
	if len(perms) != len(uniqueMenuPermissionIDs(permIDs)) {
		// 返回当前处理结果。
		return errors.New("存在无效权限ID")
	}

	// 执行权限绑定
	return s.dao.ReplacePermissions(ctx, menu, perms)
}

// GetMenuPermissions 查询菜单已绑定权限。
func (s *MenuPermissionService) GetMenuPermissions(ctx context.Context, menuID uint) ([]models.Permission, error) {
	// 返回当前处理结果。
	return s.dao.GetPermissions(ctx, menuID)
}

// uniqueMenuPermissionIDs 对权限 ID 去重。
func uniqueMenuPermissionIDs(in []uint) []uint {
	// 使用 map 去重
	set := make(map[uint]struct{}, len(in))
	// 定义并初始化当前变量。
	out := make([]uint, 0, len(in))
	// 循环处理当前数据集合。
	for _, v := range in {
		// 判断条件并进入对应分支逻辑。
		if _, ok := set[v]; ok {
			// 处理当前语句逻辑。
			continue
		}
		// 更新当前变量或字段值。
		set[v] = struct{}{}
		// 更新当前变量或字段值。
		out = append(out, v)
	}
	// 返回当前处理结果。
	return out
}
