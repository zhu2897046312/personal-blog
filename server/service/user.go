package service

import (
	"context"
	"errors"
	"time"
	"gorm.io/gorm"

	"github.com/personal-blog/models"
	"github.com/personal-blog/repository/mysql"
	"github.com/personal-blog/repository/redis"
	"github.com/personal-blog/middleware"
	"golang.org/x/crypto/bcrypt"
)

// UserService 用户服务接口
type UserService interface {
	Register(ctx context.Context, user *models.User) error
	Login(ctx context.Context, username, password string) (*models.User, string, error)
	GetUserByID(ctx context.Context, id uint) (*models.User, error)
	UpdateUser(ctx context.Context, user *models.User) error
	ListUsers(ctx context.Context, page, pageSize int) ([]models.User, int64, error)
	ChangePassword(ctx context.Context, userID uint, oldPassword, newPassword string) error
}

type userService struct {
	userRepo  mysql.UserRepository
	userCache redis.UserCache
}

// NewUserService 创建用户服务实例
func NewUserService(userRepo mysql.UserRepository, userCache redis.UserCache) UserService {
	return &userService{
		userRepo:  userRepo,
		userCache: userCache,
	}
}

func (s *userService) Register(ctx context.Context, user *models.User) error {
	// 检查用户名是否已存在
	existingUser, err := s.userRepo.FindByUsername(user.Username)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}
	if existingUser != nil {
		return errors.New("username already exists")
	}

	// 检查邮箱是否已存在
	existingUser, err = s.userRepo.FindByEmail(user.Email)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}
	if existingUser != nil {
		return errors.New("email already exists")
	}

	// 加密密码
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user.Password = string(hashedPassword)

	// 设置默认值
	user.Role = "user"
	user.Status = 1
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()

	return s.userRepo.Create(user)
}

func (s *userService) Login(ctx context.Context, username, password string) (*models.User, string, error) {
	user, err := s.userRepo.FindByUsername(username)
	if err != nil {
		return nil, "", err
	}

	if user.Status != 1 {
		return nil, "", errors.New("user is disabled")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return nil, "", errors.New("invalid password")
	}

	// 生成JWT token
	token, err := middleware.GenerateToken(user.ID, user.Username, user.Role)
	if err != nil {
		return nil, "", err
	}

	// 缓存token
	if err := s.userCache.SetUserToken(ctx, user.ID, token); err != nil {
		return nil, "", err
	}

	return user, token, nil
}

func (s *userService) GetUserByID(ctx context.Context, id uint) (*models.User, error) {
	// 先从缓存获取
	user, err := s.userCache.Get(ctx, id)
	if err != nil {
		return nil, err
	}
	if user != nil {
		return user, nil
	}

	// 缓存未命中，从数据库获取
	user, err = s.userRepo.FindByID(id)
	if err != nil {
		return nil, err
	}

	// 写入缓存
	if err := s.userCache.Set(ctx, user); err != nil {
		return nil, err
	}

	return user, nil
}

func (s *userService) UpdateUser(ctx context.Context, user *models.User) error {
	user.UpdatedAt = time.Now()
	if err := s.userRepo.Update(user); err != nil {
		return err
	}

	// 更新缓存
	return s.userCache.Set(ctx, user)
}

func (s *userService) ListUsers(ctx context.Context, page, pageSize int) ([]models.User, int64, error) {
	return s.userRepo.List(page, pageSize)
}

func (s *userService) ChangePassword(ctx context.Context, userID uint, oldPassword, newPassword string) error {
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(oldPassword)); err != nil {
		return errors.New("invalid old password")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	user.Password = string(hashedPassword)
	user.UpdatedAt = time.Now()

	if err := s.userRepo.Update(user); err != nil {
		return err
	}

	// 更新缓存
	return s.userCache.Set(ctx, user)
}
