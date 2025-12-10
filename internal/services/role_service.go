package services

import (
	"context"
	"go-admin-full/internal/dao"
	"go-admin-full/internal/models"
)

type RoleService struct {
	dao dao.RoleDAO
}

func NewRoleService(dao dao.RoleDAO) *RoleService {
	return &RoleService{dao: dao}
}

func (s *RoleService) Create(ctx context.Context, r *models.Role) error {
	return s.dao.Create(ctx, r)
}

func (s *RoleService) Update(ctx context.Context, r *models.Role) error {
	return s.dao.Update(ctx, r)
}

func (s *RoleService) List(ctx context.Context) ([]models.Role, error) {
	return s.dao.List(ctx)
}

func (s *RoleService) Get(ctx context.Context, id uint) (*models.Role, error) {
	return s.dao.Get(ctx, id)
}

func (s *RoleService) Delete(ctx context.Context, id uint) error {
	return s.dao.Delete(ctx, id)
}

func (s *RoleService) GetRolesByUserID(ctx context.Context, userID uint) ([]models.Role, error) {
	return s.dao.GetByUserID(ctx, userID)
}
