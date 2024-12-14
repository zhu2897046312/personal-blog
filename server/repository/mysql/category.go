package mysql

import (
	"github.com/personal-blog/models"
	"gorm.io/gorm"
)

// CategoryRepository 分类仓库接口
type CategoryRepository interface {
	Create(category *models.Category) error
	Update(category *models.Category) error
	Delete(id uint) error
	FindByID(id uint) (*models.Category, error)
	List() ([]models.Category, error)
	FindByName(name string) (*models.Category, error)
}

type categoryRepository struct {
	db *gorm.DB
}

// NewCategoryRepository 创建分类仓库实例
func NewCategoryRepository(db *gorm.DB) CategoryRepository {
	return &categoryRepository{db: db}
}

func (r *categoryRepository) Create(category *models.Category) error {
	return r.db.Create(category).Error
}

func (r *categoryRepository) Update(category *models.Category) error {
	return r.db.Save(category).Error
}

func (r *categoryRepository) Delete(id uint) error {
	return r.db.Delete(&models.Category{}, id).Error
}

func (r *categoryRepository) FindByID(id uint) (*models.Category, error) {
	var category models.Category
	err := r.db.First(&category, id).Error
	if err != nil {
		return nil, err
	}
	return &category, nil
}

func (r *categoryRepository) List() ([]models.Category, error) {
	var categories []models.Category
	err := r.db.Find(&categories).Error
	if err != nil {
		return nil, err
	}
	return categories, nil
}

func (r *categoryRepository) FindByName(name string) (*models.Category, error) {
	var category models.Category
	err := r.db.Where("name = ?", name).First(&category).Error
	if err != nil {
		return nil, err
	}
	return &category, nil
}
