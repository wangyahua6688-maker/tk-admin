package rbac

import (
	"context"
	rbacdao "go-admin/internal/dao/rbac"
	"go-admin/internal/models"
)

// MenuService 提供菜单业务逻辑封装。
type MenuService struct {
	dao rbacdao.MenuDAO // 菜单 DAO
}

// NewMenuService 创建菜单服务。
func NewMenuService(dao rbacdao.MenuDAO) *MenuService {
	// 返回当前处理结果。
	return &MenuService{dao: dao}
}

// Create 新增菜单。
func (s *MenuService) Create(ctx context.Context, m *models.Menu) error {
	// 返回当前处理结果。
	return s.dao.Create(ctx, m)
}

// Update 更新菜单。
func (s *MenuService) Update(ctx context.Context, m *models.Menu) error {
	// 返回当前处理结果。
	return s.dao.Update(ctx, m)
}

// List 查询全部菜单。
func (s *MenuService) List(ctx context.Context) ([]models.Menu, error) {
	// 返回当前处理结果。
	return s.dao.List(ctx)
}

// Get 获取单个菜单。
func (s *MenuService) Get(ctx context.Context, id uint) (*models.Menu, error) {
	// 返回当前处理结果。
	return s.dao.Get(ctx, id)
}

// Delete 删除菜单。
func (s *MenuService) Delete(ctx context.Context, id uint) error {
	// 返回当前处理结果。
	return s.dao.Delete(ctx, id)
}

// ListForUser 按用户可访问权限返回菜单列表。
func (s *MenuService) ListForUser(ctx context.Context, userID uint) ([]models.Menu, error) {
	// 返回当前处理结果。
	return s.dao.ListByUserID(ctx, userID)
}

// BuildMenuTreeFromList 将指定菜单集合构建成树结构（用于按权限返回前端路由）。
func (s *MenuService) BuildMenuTreeFromList(menus []models.Menu) []*models.Menu {
	// 返回当前处理结果。
	return s.buildTree(menus)
}

// buildTree 将扁平菜单转换为树结构。
func (s *MenuService) buildTree(menus []models.Menu) []*models.Menu {
	// 转成 map
	m := make(map[uint]*models.Menu, len(menus))
	// 循环处理当前数据集合。
	for _, menu := range menus {
		// 显式拷贝字段，避免外部复用导致的 Children 引用污染
		node := &models.Menu{
			// 处理当前语句逻辑。
			ID: menu.ID,
			// 处理当前语句逻辑。
			ParentID: menu.ParentID,
			// 处理当前语句逻辑。
			Title: menu.Title,
			// 处理当前语句逻辑。
			Path: menu.Path,
			// 处理当前语句逻辑。
			Icon: menu.Icon,
			// 处理当前语句逻辑。
			Component: menu.Component,
			// 处理当前语句逻辑。
			OrderNum: menu.OrderNum,
			// 处理当前语句逻辑。
			Permissions: menu.Permissions,
			// 处理当前语句逻辑。
			Children: []*models.Menu{},
		}
		// 更新当前变量或字段值。
		m[menu.ID] = node
	}

	// 构建树
	tree := make([]*models.Menu, 0)
	// 循环处理当前数据集合。
	for _, raw := range menus {
		// 定义并初始化当前变量。
		item := m[raw.ID]
		// 父级为 0 的视作根节点
		if item.ParentID == 0 {
			// 更新当前变量或字段值。
			tree = append(tree, item)
			// 处理当前语句逻辑。
			continue
		}
		// 若父节点存在则挂载到父级 Children
		if parent, ok := m[item.ParentID]; ok {
			// 更新当前变量或字段值。
			parent.Children = append(parent.Children, item)
			// 处理当前语句逻辑。
			continue
		}

		// 容错：若父节点未在当前权限集合中，避免子节点丢失，降级为根节点返回。
		tree = append(tree, item)
	}

	// 返回当前处理结果。
	return tree
}
