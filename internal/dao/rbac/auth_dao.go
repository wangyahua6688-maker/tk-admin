package rbac

import (
	"context"
	"go-admin-full/internal/models"
	"go-admin-full/internal/utils"
	"gorm.io/gorm"
)

// AuthDao 认证相关数据访问层。
type AuthDao struct {
	db *gorm.DB // 数据库连接
}

// NewAuthDao 创建认证 DAO。
func NewAuthDao(db *gorm.DB) *AuthDao {
	return &AuthDao{db: db}
}

// GetUserByUsername 根据用户名查询用户
func (d *AuthDao) GetUserByUsername(ctx context.Context, username string) (*models.User, error) {
	var user models.User
	db := d.db // 默认使用注入的连接
	if ctx != nil {
		// 优先使用上下文中的数据库连接，确保中间件注入事务/连接时可透传。
		if ctxDB := utils.DBFromContext(ctx); ctxDB != nil {
			db = ctxDB
		}
	}

	// 按用户名查询
	if err := db.WithContext(ctx).Where("username = ?", username).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

// CreateUser 创建用户（注册用）
func (d *AuthDao) CreateUser(ctx context.Context, user *models.User) error {
	return d.db.WithContext(ctx).Create(user).Error
}

// GetUserByID 根据用户ID查询用户。
// 说明：该方法用于 JWT 中间件与 refresh 流程中的账号状态校验。
func (d *AuthDao) GetUserByID(ctx context.Context, userID uint) (*models.User, error) {
	var user models.User
	// 按 ID 查询
	if err := d.db.WithContext(ctx).Where("id = ?", userID).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}
