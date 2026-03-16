package rbac

import (
	"context"
	"errors"
	rbacdao "go-admin-full/internal/dao/rbac"
	"go-admin-full/internal/models"
	"golang.org/x/crypto/bcrypt"
	"strings"
)

type UserService struct {
	userDao *rbacdao.UserDao // 用户 DAO
}

// NewUserService 创建用户服务。
func NewUserService(userDao *rbacdao.UserDao) *UserService {
	return &UserService{userDao: userDao}
}

func (s *UserService) CreateUser(ctx context.Context, username, password, email, avatar string) (*models.User, error) {
	username = strings.TrimSpace(username) // 规范化用户名
	email = strings.TrimSpace(email)       // 规范化邮箱
	avatar = strings.TrimSpace(avatar)     // 规范化头像
	if username == "" {
		return nil, errors.New("用户名不能为空")
	}

	// 检查用户名是否存在
	existing, err := s.userDao.GetByUsername(ctx, username)
	if err == nil && existing != nil {
		return nil, errors.New("用户名已存在")
	}

	// 密码加密
	if err := ValidatePasswordStrength(password); err != nil {
		return nil, err
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	// 创建用户
	user := &models.User{
		Username:     username,
		PasswordHash: string(hashedPassword),
		Email:        email,
		Avatar:       avatar,
		Status:       1,
	}

	if err := s.userDao.Create(ctx, user); err != nil {
		return nil, err
	}

	return user, nil
}

func (s *UserService) ListAllUsers(ctx context.Context) ([]models.User, error) {
	return s.userDao.ListAll(ctx)
}

// GetUserByID 查询用户基础信息。
func (s *UserService) GetUserByID(ctx context.Context, userID uint) (*models.User, error) {
	return s.userDao.GetByID(ctx, userID)
}

// UpdateUser 更新用户资料。
// 安全规则：
// 1. 用户名不可通过该接口修改；
// 2. 密码非空时会重新哈希后存储。
func (s *UserService) UpdateUser(ctx context.Context, userID uint, email string, status *int, password string, avatar *string) error {
	updates := make(map[string]interface{})

	if strings.TrimSpace(email) != "" {
		updates["email"] = strings.TrimSpace(email)
	}

	if status != nil {
		updates["status"] = *status
	}

	// avatar 允许显式置空，因此以指针判断“是否传入”。
	if avatar != nil {
		updates["avatar"] = strings.TrimSpace(*avatar)
	}

	if strings.TrimSpace(password) != "" {
		if err := ValidatePasswordStrength(password); err != nil {
			return err
		}
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		if err != nil {
			return err
		}
		updates["password_hash"] = string(hashedPassword)
	}

	if len(updates) == 0 {
		return errors.New("无可更新字段")
	}

	return s.userDao.UpdateByID(ctx, userID, updates)
}

// DeleteUser 删除用户。
func (s *UserService) DeleteUser(ctx context.Context, userID uint) error {
	return s.userDao.DeleteByID(ctx, userID)
}
