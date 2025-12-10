package services

import (
	"context"
	"go-admin-full/internal/dao"
	"go-admin-full/internal/models"
)

type PermissionService struct {
	dao dao.PermissionDAO
}

func NewPermissionService(dao dao.PermissionDAO) *PermissionService {
	return &PermissionService{dao: dao}
}

func (s *PermissionService) Create(ctx context.Context, p *models.Permission) error {
	return s.dao.Create(ctx, p)
}

func (s *PermissionService) Update(ctx context.Context, p *models.Permission) error {
	return s.dao.Update(ctx, p)
}

func (s *PermissionService) List(ctx context.Context) ([]models.Permission, error) {
	return s.dao.List(ctx)
}

func (s *PermissionService) Get(ctx context.Context, id uint) (*models.Permission, error) {
	return s.dao.Get(ctx, id)
}

func (s *PermissionService) Delete(ctx context.Context, id uint) error {
	return s.dao.Delete(ctx, id)
}

func (s *PermissionService) GetPermissionsByRoleIDs(ctx context.Context, roleIDs []uint) ([]models.Permission, error) {
	return s.dao.GetByRoleIDs(ctx, roleIDs)
}
