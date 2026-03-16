package biz

import "gorm.io/gorm"

// UserOpsDAO 聚合用户运营相关 DAO。
// 说明：具体方法拆分在同包不同文件中。
type UserOpsDAO struct {
	db *gorm.DB // 数据库连接
}

// NewUserOpsDAO 创建用户运营 DAO。
func NewUserOpsDAO(db *gorm.DB) *UserOpsDAO {
	return &UserOpsDAO{db: db}
}
