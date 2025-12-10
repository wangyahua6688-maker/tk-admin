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

func (s *MenuService) BuildMenuTree(ctx context.Context) ([]*models.Menu, error) {
	menus, err := s.dao.List(ctx)
	if err != nil {
		return nil, err
	}

	// 转成 map
	m := make(map[int]*models.Menu)
	for _, menu := range menus {
		node := &models.Menu{
			ID:        menu.ID,
			ParentID:  menu.ParentID,
			Title:     menu.Title,
			Path:      menu.Path,
			Icon:      menu.Icon,
			Component: menu.Component,
			Children:  []*models.Menu{},
		}
		m[menu.ID] = node
	}

	// 构建树
	tree := make([]*models.Menu, 0)
	for _, item := range m {
		if item.ParentID == 0 {
			tree = append(tree, item)
			continue
		}
		if parent, ok := m[item.ParentID]; ok {
			parent.Children = append(parent.Children, item)
		}
	}

	return tree, nil
}
