package biz

import (
	bizdao "go-admin-full/internal/dao/biz"
	"gorm.io/gorm"
)

// UserOpsService 聚合用户运营相关服务。
type UserOpsService struct {
	dao *bizdao.UserOpsDAO // 用户运营 DAO
	db  *gorm.DB           // 数据库连接
}

// NewUserOpsService 创建用户运营服务。
func NewUserOpsService(db *gorm.DB) *UserOpsService {
	return &UserOpsService{
		dao: bizdao.NewUserOpsDAO(db),
		db:  db,
	}
}
