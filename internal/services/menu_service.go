package services

import (
	"context"

	"go-admin-full/internal/models"
	"gorm.io/gorm"
)

type MenuService struct {
	db *gorm.DB
}

func NewMenuService(db *gorm.DB) *MenuService { return &MenuService{db: db} }

func (s *MenuService) Create(ctx context.Context, m *models.Menu) error {
	return s.db.WithContext(ctx).Create(m).Error
}
func (s *MenuService) List(ctx context.Context) ([]models.Menu, error) {
	var items []models.Menu
	err := s.db.WithContext(ctx).Find(&items).Error
	return items, err
}
func (s *MenuService) Get(ctx context.Context, id uint) (*models.Menu, error) {
	var m models.Menu
	if err := s.db.WithContext(ctx).First(&m, id).Error; err != nil {
		return nil, err
	}
	return &m, nil
}
func (s *MenuService) Delete(ctx context.Context, id uint) error {
	return s.db.WithContext(ctx).Delete(&models.Menu{}, id).Error
}

func (s *MenuService) BuildMenuTree(menus []models.Menu) []*models.Menu {
	idMap := make(map[uint]*models.Menu)

	// 复制并建立索引
	for i := range menus {
		m := menus[i]
		m.Children = []*models.Menu{}
		idMap[m.ID] = &m
	}

	var roots []*models.Menu

	for _, m := range idMap {
		if m.ParentID == nil {
			roots = append(roots, m)
		} else {
			parent := idMap[*m.ParentID]
			if parent != nil {
				parent.Children = append(parent.Children, m)
			}
		}
	}

	return roots
}

func (s *MenuService) ListForRoles(ctx context.Context, roleCodes []string) ([]*models.Menu, error) {
	var menus []models.Menu

	// 查询 role → permission → menu
	err := s.db.WithContext(ctx).
		Preload("Permissions").
		Order("`order` asc").
		Find(&menus).Error

	if err != nil {
		return nil, err
	}

	if len(roleCodes) == 0 {
		return nil, nil
	}

	// 1）查询角色所有权限
	var perms []models.Permission
	err = s.db.WithContext(ctx).
		Joins("JOIN role_permissions rp ON rp.permission_id = permissions.id").
		Joins("JOIN roles r ON r.id = rp.role_id AND r.code IN ?", roleCodes).
		Find(&perms).Error

	if err != nil {
		return nil, err
	}

	permSet := map[string]bool{}
	for _, p := range perms {
		permSet[p.Code] = true
	}

	// 2）过滤菜单（菜单没有要求权限 → 所有人可以访问）
	var filtered []models.Menu
	for _, m := range menus {
		if len(m.Permissions) == 0 {
			filtered = append(filtered, m)
			continue
		}
		allowed := false
		for _, p := range m.Permissions {
			if permSet[p.Code] {
				allowed = true
				break
			}
		}
		if allowed {
			filtered = append(filtered, m)
		}
	}

	// 3）构建树结构
	return s.BuildMenuTree(filtered), nil
}
