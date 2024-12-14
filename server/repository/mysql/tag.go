package mysql

import (
	"github.com/personal-blog/models"
	"gorm.io/gorm"
)

// TagRepository 标签仓库接口
type TagRepository interface {
	Create(tag *models.Tag) error
	Update(tag *models.Tag) error
	Delete(id uint) error
	FindByID(id uint) (*models.Tag, error)
	List() ([]models.Tag, error)
	FindByName(name string) (*models.Tag, error)
	BatchCreate(tags []models.Tag) error
	FindOrCreateByNames(names []string) ([]models.Tag, error)
}

type tagRepository struct {
	db *gorm.DB
}

// NewTagRepository 创建标签仓库实例
func NewTagRepository(db *gorm.DB) TagRepository {
	return &tagRepository{db: db}
}

func (r *tagRepository) Create(tag *models.Tag) error {
	return r.db.Create(tag).Error
}

func (r *tagRepository) Update(tag *models.Tag) error {
	return r.db.Save(tag).Error
}

func (r *tagRepository) Delete(id uint) error {
	return r.db.Delete(&models.Tag{}, id).Error
}

func (r *tagRepository) FindByID(id uint) (*models.Tag, error) {
	var tag models.Tag
	err := r.db.First(&tag, id).Error
	if err != nil {
		return nil, err
	}
	return &tag, nil
}

func (r *tagRepository) List() ([]models.Tag, error) {
	var tags []models.Tag
	err := r.db.Find(&tags).Error
	if err != nil {
		return nil, err
	}
	return tags, nil
}

func (r *tagRepository) FindByName(name string) (*models.Tag, error) {
	var tag models.Tag
	err := r.db.Where("name = ?", name).First(&tag).Error
	if err != nil {
		return nil, err
	}
	return &tag, nil
}

func (r *tagRepository) BatchCreate(tags []models.Tag) error {
	return r.db.Create(&tags).Error
}

func (r *tagRepository) FindOrCreateByNames(names []string) ([]models.Tag, error) {
	var tags []models.Tag
	for _, name := range names {
		var tag models.Tag
		err := r.db.Where("name = ?", name).FirstOrCreate(&tag, models.Tag{Name: name}).Error
		if err != nil {
			return nil, err
		}
		tags = append(tags, tag)
	}
	return tags, nil
}
