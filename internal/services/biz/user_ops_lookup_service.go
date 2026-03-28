package biz

import (
	"context"

	bizdao "go-admin/internal/dao/biz"
	"go-admin/internal/models"
	"gorm.io/gorm"
)

// UserOpsLookupService 用户运营公共查询服务。
type UserOpsLookupService struct {
	dao *bizdao.UserOpsLookupDAO
}

// NewUserOpsLookupService 创建用户运营公共查询服务。
func NewUserOpsLookupService(db *gorm.DB) *UserOpsLookupService {
	return &UserOpsLookupService{dao: bizdao.NewUserOpsLookupDAO(db)}
}

// IsUserTypes 判断用户是否为指定类型。
func (s *UserOpsLookupService) IsUserTypes(ctx context.Context, userID uint, expectedTypes ...string) bool {
	if userID == 0 {
		return false
	}
	current, err := s.dao.GetActiveUserType(ctx, userID)
	if err != nil {
		return false
	}
	for _, t := range expectedTypes {
		if current == t {
			return true
		}
	}
	return false
}

// GetUsersByIDs 批量获取用户基础信息。
func (s *UserOpsLookupService) GetUsersByIDs(ctx context.Context, userIDs []uint) ([]models.WUser, error) {
	return s.dao.GetUsersByIDs(ctx, userIDs)
}
