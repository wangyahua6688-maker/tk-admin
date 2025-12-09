package services

import (
	"context"
	"errors"
	"go-admin-full/internal/dao"
	"go-admin-full/internal/models"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	userDao *dao.UserDao
}

func NewUserService(userDao *dao.UserDao) *UserService {
	return &UserService{userDao: userDao}
}

func (s *UserService) CreateUser(ctx context.Context, username, password, email string) (*models.User, error) {
	// 检查用户名是否存在
	existing, err := s.userDao.GetByUsername(ctx, username)
	if err == nil && existing != nil {
		return nil, errors.New("用户名已存在")
	}

	// 密码加密
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
