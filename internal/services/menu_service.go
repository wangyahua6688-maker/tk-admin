package services

import (
	"context"
	"go-admin-full/internal/dao"
	"go-admin-full/internal/models"
)

type MenuService struct {
	dao dao.MenuDAO
}

func NewMenuService(dao dao.MenuDAO) *MenuService {
	return &MenuService{dao: dao}
}

func (s *MenuService) Create(ctx context.Context, m *models.Menu) error {
	return s.dao.Create(ctx, m)
}

func (s *MenuService) Update(ctx context.Context, m *models.Menu) error {
	return s.dao.Update(ctx, m)
}

func (s *MenuService) List(ctx context.Context) ([]models.Menu, error) {
	return s.dao.List(ctx)
}

func (s *MenuService) Get(ctx context.Context, id uint) (*models.Menu, error) {
	return s.dao.Get(ctx, id)
}

func (s *MenuService) Delete(ctx context.Context, id uint) error {
	return s.dao.Delete(ctx, id)
}

// ListForUser 按用户可访问权限返回菜单列表。
func (s *MenuService) ListForUser(ctx context.Context, userID uint) ([]models.Menu, error) {
	return s.dao.ListByUserID(ctx, userID)
}

func (s *MenuService) BuildMenuTree(ctx context.Context) ([]*models.Menu, error) {
	menus, err := s.dao.List(ctx)
	if err != nil {
		return nil, err
	}
	return s.buildTree(menus), nil
}

// BuildMenuTreeFromList 将指定菜单集合构建成树结构（用于按权限返回前端路由）。
func (s *MenuService) BuildMenuTreeFromList(menus []models.Menu) []*models.Menu {
	return s.buildTree(menus)
}

// buildTree 将扁平菜单转换为树结构。
func (s *MenuService) buildTree(menus []models.Menu) []*models.Menu {
	// 转成 map
	m := make(map[uint]*models.Menu, len(menus))
	for _, menu := range menus {
		node := &models.Menu{
			ID:          menu.ID,
			ParentID:    menu.ParentID,
			Title:       menu.Title,
			Path:        menu.Path,
			Icon:        menu.Icon,
			Component:   menu.Component,
			OrderNum:    menu.OrderNum,
			Permissions: menu.Permissions,
			Children:    []*models.Menu{},
		}
		m[menu.ID] = node
	}

	// 构建树
	tree := make([]*models.Menu, 0)
	for _, raw := range menus {
		item := m[raw.ID]
		if item.ParentID == 0 {
			tree = append(tree, item)
			continue
		}
		if parent, ok := m[item.ParentID]; ok {
			parent.Children = append(parent.Children, item)
			continue
		}

		// 容错：若父节点未在当前权限集合中，避免子节点丢失，降级为根节点返回。
		tree = append(tree, item)
	}

	return tree
}
