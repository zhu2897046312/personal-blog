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

// TagService 标签服务接口
type TagService interface {
	CreateTag(ctx context.Context, tag *models.Tag) error
	UpdateTag(ctx context.Context, tag *models.Tag) error
	DeleteTag(ctx context.Context, id uint) error
	GetTagByID(ctx context.Context, id uint) (*models.Tag, error)
	ListTags(ctx context.Context) ([]models.Tag, error)
	GetPostTags(ctx context.Context, postID uint) ([]models.Tag, error)
	CreateTagsIfNotExist(ctx context.Context, names []string) ([]models.Tag, error)
}

type tagService struct {
	tagRepo   mysql.TagRepository
	tagCache  redis.TagCache
	postRepo  mysql.PostRepository
}

// NewTagService 创建标签服务实例
func NewTagService(tagRepo mysql.TagRepository, tagCache redis.TagCache, postRepo mysql.PostRepository) TagService {
	return &tagService{
		tagRepo:  tagRepo,
		tagCache: tagCache,
		postRepo: postRepo,
	}
}

func (s *tagService) CreateTag(ctx context.Context, tag *models.Tag) error {
	// 检查名称是否已存在
	existing, err := s.tagRepo.FindByName(tag.Name)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}
	if existing != nil {
		return errors.New("tag name already exists")
	}

	tag.CreatedAt = time.Now()
	tag.UpdatedAt = time.Now()

	// 创建标签
	if err := s.tagRepo.Create(tag); err != nil {
		return err
	}

	// 写入缓存
	if err := s.tagCache.Set(ctx, tag); err != nil {
		return err
	}

	// 清除标签列表缓存
	return s.tagCache.Delete(ctx, 0)
}

func (s *tagService) UpdateTag(ctx context.Context, tag *models.Tag) error {
	// 检查名称是否已存在（排除自身）
	existing, err := s.tagRepo.FindByName(tag.Name)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}
	if existing != nil && existing.ID != tag.ID {
		return errors.New("tag name already exists")
	}

	tag.UpdatedAt = time.Now()

	// 更新标签
	if err := s.tagRepo.Update(tag); err != nil {
		return err
	}

	// 更新缓存
	if err := s.tagCache.Set(ctx, tag); err != nil {
		return err
	}

	// 清除标签列表缓存
	return s.tagCache.Delete(ctx, 0)
}

func (s *tagService) DeleteTag(ctx context.Context, id uint) error {
	// 删除标签
	if err := s.tagRepo.Delete(id); err != nil {
		return err
	}

	// 删除缓存
	if err := s.tagCache.Delete(ctx, id); err != nil {
		return err
	}

	// 清除标签列表缓存
	return s.tagCache.Delete(ctx, 0)
}

func (s *tagService) GetTagByID(ctx context.Context, id uint) (*models.Tag, error) {
	// 先从缓存获取
	tag, err := s.tagCache.Get(ctx, id)
	if err != nil {
		return nil, err
	}
	if tag != nil {
		return tag, nil
	}

	// 缓存未命中，从数据库获取
	tag, err = s.tagRepo.FindByID(id)
	if err != nil {
		return nil, err
	}

	// 写入缓存
	if err := s.tagCache.Set(ctx, tag); err != nil {
		return nil, err
	}

	return tag, nil
}

func (s *tagService) ListTags(ctx context.Context) ([]models.Tag, error) {
	// 先从缓存获取
	tags, err := s.tagCache.GetList(ctx)
	if err != nil {
		return nil, err
	}
	if len(tags) > 0 {
		return tags, nil
	}

	// 缓存未命中，从数据库获取
	tags, err = s.tagRepo.List()
	if err != nil {
		return nil, err
	}

	// 写入缓存
	if err := s.tagCache.SetList(ctx, tags); err != nil {
		return nil, err
	}

	return tags, nil
}

func (s *tagService) GetPostTags(ctx context.Context, postID uint) ([]models.Tag, error) {
	// 先从缓存获取
	tags, err := s.tagCache.GetPostTags(ctx, postID)
	if err != nil {
		return nil, err
	}
	if len(tags) > 0 {
		return tags, nil
	}

	// 从数据库获取
	post, err := s.postRepo.FindByID(postID)
	if err != nil {
		return nil, err
	}

	// 写入缓存
	if err := s.tagCache.SetPostTags(ctx, postID, post.Tags); err != nil {
		return nil, err
	}

	return post.Tags, nil
}

func (s *tagService) CreateTagsIfNotExist(ctx context.Context, names []string) ([]models.Tag, error) {
	return s.tagRepo.FindOrCreateByNames(names)
}
