package services

import (
	"context"
	"errors"
	"go-admin-full/internal/dao"
	"go-admin-full/internal/models"
	"golang.org/x/crypto/bcrypt"
	"strings"
)

type UserService struct {
	userDao *dao.UserDao
}

func NewUserService(userDao *dao.UserDao) *UserService {
	return &UserService{userDao: userDao}
}

func (s *UserService) CreateUser(ctx context.Context, username, password, email string) (*models.User, error) {
	username = strings.TrimSpace(username)
	email = strings.TrimSpace(email)
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
func (s *UserService) UpdateUser(ctx context.Context, userID uint, email string, status *int, password string) error {
	updates := make(map[string]interface{})

	if strings.TrimSpace(email) != "" {
		updates["email"] = strings.TrimSpace(email)
	}

	if status != nil {
		updates["status"] = *status
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
