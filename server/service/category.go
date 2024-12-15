package service

import (
	"context"
	"errors"
	"time"
	"gorm.io/gorm"

	"github.com/personal-blog/models"
	"github.com/personal-blog/repository/mysql"
	"github.com/personal-blog/repository/redis"
)

// CategoryService 分类服务接口
type CategoryService interface {
	CreateCategory(ctx context.Context, category *models.Category) error
	UpdateCategory(ctx context.Context, category *models.Category) error
	DeleteCategory(ctx context.Context, id uint) error
	GetCategoryByID(ctx context.Context, id uint) (*models.Category, error)
	ListCategories(ctx context.Context, page, pageSize int) ([]models.Category, int64, error)
}

type categoryService struct {
	categoryRepo  mysql.CategoryRepository
	categoryCache redis.CategoryCache
}

// NewCategoryService 创建分类服务实例
func NewCategoryService(categoryRepo mysql.CategoryRepository, categoryCache redis.CategoryCache) CategoryService {
	return &categoryService{
		categoryRepo:  categoryRepo,
		categoryCache: categoryCache,
	}
}

func (s *categoryService) CreateCategory(ctx context.Context, category *models.Category) error {
	// 检查名称是否已存在
	existing, err := s.categoryRepo.FindByName(category.Name)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}
	if existing != nil {
		return errors.New("category name already exists")
	}

	category.CreatedAt = time.Now()
	category.UpdatedAt = time.Now()

	// 创建分类
	if err := s.categoryRepo.Create(category); err != nil {
		return err
	}

	// 写入缓存
	if err := s.categoryCache.Set(ctx, category); err != nil {
		return err
	}

	// 清除分类列表缓存
	return s.categoryCache.Delete(ctx, 0)
}

func (s *categoryService) UpdateCategory(ctx context.Context, category *models.Category) error {
	// 检查名称是否已存在（排除自身）
	existing, err := s.categoryRepo.FindByName(category.Name)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}
	if existing != nil && existing.ID != category.ID {
		return errors.New("category name already exists")
	}

	category.UpdatedAt = time.Now()

	// 更新分类
	if err := s.categoryRepo.Update(category); err != nil {
		return err
	}

	// 更新缓存
	if err := s.categoryCache.Set(ctx, category); err != nil {
		return err
	}

	// 清除分类列表缓存
	return s.categoryCache.Delete(ctx, 0)
}

func (s *categoryService) DeleteCategory(ctx context.Context, id uint) error {
	// 删除分类
	if err := s.categoryRepo.Delete(id); err != nil {
		return err
	}

	// 删除缓存
	if err := s.categoryCache.Delete(ctx, id); err != nil {
		return err
	}

	// 清除分类列表缓存
	return s.categoryCache.Delete(ctx, 0)
}

func (s *categoryService) GetCategoryByID(ctx context.Context, id uint) (*models.Category, error) {
	// 先从缓存获取
	category, err := s.categoryCache.Get(ctx, id)
	if err != nil {
		return nil, err
	}
	if category != nil {
		return category, nil
	}

	// 缓存未命中，从数据库获取
	category, err = s.categoryRepo.FindByID(id)
	if err != nil {
		return nil, err
	}

	// 写入缓存
	if err := s.categoryCache.Set(ctx, category); err != nil {
		return nil, err
	}

	return category, nil
}

func (s *categoryService) ListCategories(ctx context.Context, page, pageSize int) ([]models.Category, int64, error) {
	return s.categoryRepo.List(page, pageSize)
}
