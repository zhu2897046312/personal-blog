package mysql

import (
	"github.com/personal-blog/models"
	"gorm.io/gorm"
)

// PostRepository 文章仓库接口
type PostRepository interface {
	Create(post *models.Post) error
	Update(post *models.Post) error
	Delete(id uint) error
	FindByID(id uint) (*models.Post, error)
	List(page, pageSize int, conditions map[string]interface{}) ([]models.Post, int64, error)
	IncrementViewCount(id uint) error
	ListByUserID(userID uint, page, pageSize int) ([]models.Post, int64, error)
	ListByCategoryID(categoryID uint, page, pageSize int) ([]models.Post, int64, error)
	ListByTagID(tagID uint, page, pageSize int) ([]models.Post, int64, error)
}

type postRepository struct {
	db *gorm.DB
}

// NewPostRepository 创建文章仓库实例
func NewPostRepository(db *gorm.DB) PostRepository {
	return &postRepository{db: db}
}

func (r *postRepository) Create(post *models.Post) error {
	return r.db.Create(post).Error
}

func (r *postRepository) Update(post *models.Post) error {
	return r.db.Save(post).Error
}

func (r *postRepository) Delete(id uint) error {
	return r.db.Delete(&models.Post{}, id).Error
}

func (r *postRepository) FindByID(id uint) (*models.Post, error) {
	var post models.Post
	err := r.db.Preload("User").
		Preload("Category").
		Preload("Tags").
		Preload("Comments").
		First(&post, id).Error
	if err != nil {
		return nil, err
	}
	return &post, nil
}

func (r *postRepository) List(page, pageSize int, conditions map[string]interface{}) ([]models.Post, int64, error) {
	var posts []models.Post
	var total int64

	query := r.db.Model(&models.Post{})
	
	// 应用查询条件
	for key, value := range conditions {
		query = query.Where(key, value)
	}

	err := query.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	err = query.Preload("User").
		Preload("Category").
		Preload("Tags").
		Offset(offset).
		Limit(pageSize).
		Order("is_top DESC, created_at DESC").
		Find(&posts).Error
	if err != nil {
		return nil, 0, err
	}

	return posts, total, nil
}

func (r *postRepository) IncrementViewCount(id uint) error {
	return r.db.Model(&models.Post{}).Where("id = ?", id).
		UpdateColumn("view_count", gorm.Expr("view_count + ?", 1)).Error
}

func (r *postRepository) ListByUserID(userID uint, page, pageSize int) ([]models.Post, int64, error) {
	return r.List(page, pageSize, map[string]interface{}{"user_id": userID})
}

func (r *postRepository) ListByCategoryID(categoryID uint, page, pageSize int) ([]models.Post, int64, error) {
	return r.List(page, pageSize, map[string]interface{}{"category_id": categoryID})
}

func (r *postRepository) ListByTagID(tagID uint, page, pageSize int) ([]models.Post, int64, error) {
	var posts []models.Post
	var total int64

	subQuery := r.db.Table("post_tags").Select("post_id").Where("tag_id = ?", tagID)
	
	err := r.db.Model(&models.Post{}).Where("id IN (?)", subQuery).Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	err = r.db.Preload("User").
		Preload("Category").
		Preload("Tags").
		Where("id IN (?)", subQuery).
		Offset(offset).
		Limit(pageSize).
		Order("created_at DESC").
		Find(&posts).Error
	if err != nil {
		return nil, 0, err
	}

	return posts, total, nil
}
