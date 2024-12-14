package mysql

import (
	"github.com/personal-blog/models"
	"gorm.io/gorm"
)

// CommentRepository 评论仓库接口
type CommentRepository interface {
	Create(comment *models.Comment) error
	Update(comment *models.Comment) error
	Delete(id uint) error
	FindByID(id uint) (*models.Comment, error)
	ListByPostID(postID uint, page, pageSize int) ([]models.Comment, int64, error)
	ListByUserID(userID uint, page, pageSize int) ([]models.Comment, int64, error)
}

type commentRepository struct {
	db *gorm.DB
}

// NewCommentRepository 创建评论仓库实例
func NewCommentRepository(db *gorm.DB) CommentRepository {
	return &commentRepository{db: db}
}

func (r *commentRepository) Create(comment *models.Comment) error {
	return r.db.Create(comment).Error
}

func (r *commentRepository) Update(comment *models.Comment) error {
	return r.db.Save(comment).Error
}

func (r *commentRepository) Delete(id uint) error {
	return r.db.Delete(&models.Comment{}, id).Error
}

func (r *commentRepository) FindByID(id uint) (*models.Comment, error) {
	var comment models.Comment
	err := r.db.Preload("User").
		Preload("Post").
		First(&comment, id).Error
	if err != nil {
		return nil, err
	}
	return &comment, nil
}

func (r *commentRepository) ListByPostID(postID uint, page, pageSize int) ([]models.Comment, int64, error) {
	var comments []models.Comment
	var total int64

	err := r.db.Model(&models.Comment{}).
		Where("post_id = ?", postID).
		Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	err = r.db.Where("post_id = ?", postID).
		Preload("User").
		Offset(offset).
		Limit(pageSize).
		Order("created_at DESC").
		Find(&comments).Error
	if err != nil {
		return nil, 0, err
	}

	return comments, total, nil
}

func (r *commentRepository) ListByUserID(userID uint, page, pageSize int) ([]models.Comment, int64, error) {
	var comments []models.Comment
	var total int64

	err := r.db.Model(&models.Comment{}).
		Where("user_id = ?", userID).
		Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	err = r.db.Where("user_id = ?", userID).
		Preload("Post").
		Offset(offset).
		Limit(pageSize).
		Order("created_at DESC").
		Find(&comments).Error
	if err != nil {
		return nil, 0, err
	}

	return comments, total, nil
}
