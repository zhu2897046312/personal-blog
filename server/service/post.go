package service

import (
	"context"
	"time"

	"github.com/personal-blog/models"
	"github.com/personal-blog/pkg/utils"
	"github.com/personal-blog/repository/mysql"
	"github.com/personal-blog/repository/redis"
)

// PostService 文章服务接口
type PostService interface {
	CreatePost(ctx context.Context, post *models.Post, tagNames []string) error
	UpdatePost(ctx context.Context, post *models.Post, tagNames []string) error
	DeletePost(ctx context.Context, id uint) error
	GetPostByID(ctx context.Context, id uint) (*models.Post, error)
	ListPosts(ctx context.Context, page, pageSize int, conditions map[string]interface{}) ([]models.Post, int64, error)
	IncrementViewCount(ctx context.Context, id uint) error
	ListPostsByCategory(ctx context.Context, categoryID uint, page, pageSize int) ([]models.Post, int64, error)
	ListPostsByTag(ctx context.Context, tagID uint, page, pageSize int) ([]models.Post, int64, error)
	ListPostsByUser(ctx context.Context, userID uint, page, pageSize int) ([]models.Post, int64, error)
}

type postService struct {
	postRepo     mysql.PostRepository
	tagRepo      mysql.TagRepository
	categoryRepo mysql.CategoryRepository
	postCache    redis.PostCache
}

// NewPostService 创建文章服务实例
func NewPostService(
	postRepo mysql.PostRepository,
	postCache redis.PostCache,
	tagRepo mysql.TagRepository,
	categoryRepo mysql.CategoryRepository,
) PostService {
	return &postService{
		postRepo:     postRepo,
		postCache:    postCache,
		tagRepo:      tagRepo,
		categoryRepo: categoryRepo,
	}
}

func (s *postService) CreatePost(ctx context.Context, post *models.Post, tagNames []string) error {
	// 处理标签
	if len(tagNames) > 0 {
		tags, err := s.tagRepo.FindOrCreateByNames(tagNames)
		if err != nil {
			return err
		}
		post.Tags = tags
	}

	// 设置时间
	post.CreatedAt = time.Now()
	post.UpdatedAt = time.Now()

	// 创建文章
	if err := s.postRepo.Create(post); err != nil {
		return err
	}

	// 写入缓存
	return s.postCache.Set(ctx, post)
}

func (s *postService) UpdatePost(ctx context.Context, post *models.Post, tagNames []string) error {
	// 处理标签
	if len(tagNames) > 0 {
		tags, err := s.tagRepo.FindOrCreateByNames(tagNames)
		if err != nil {
			return err
		}
		post.Tags = tags
	}

	post.UpdatedAt = time.Now()

	// 更新文章
	if err := s.postRepo.Update(post); err != nil {
		return err
	}

	// 更新缓存
	return s.postCache.Set(ctx, post)
}

func (s *postService) DeletePost(ctx context.Context, id uint) error {
	// 删除文章
	if err := s.postRepo.Delete(id); err != nil {
		return err
	}

	// 删除缓存
	return s.postCache.Delete(ctx, id)
}

func (s *postService) GetPostByID(ctx context.Context, id uint) (*models.Post, error) {
	// 先从缓存获取
	post, err := s.postCache.Get(ctx, id)
	if err != nil {
		return nil, err
	}
	if post != nil {
		return post, nil
	}

	// 缓存未命中，从数据库获取
	post, err = s.postRepo.FindByID(id)
	if err != nil {
		return nil, err
	}

	// 写入缓存
	if err := s.postCache.Set(ctx, post); err != nil {
		return nil, err
	}

	return post, nil
}

func (s *postService) ListPosts(ctx context.Context, page, pageSize int, conditions map[string]interface{}) ([]models.Post, int64, error) {
	// 生成缓存key
	cacheKey := "posts:list:" + utils.GenerateCacheKey(conditions, page, pageSize)

	// 尝试从缓存获取
	posts, err := s.postCache.GetPostList(ctx, cacheKey)
	if err != nil {
		return nil, 0, err
	}
	if len(posts) > 0 {
		return posts, int64(len(posts)), nil
	}

	// 从数据库获取
	posts, total, err := s.postRepo.List(page, pageSize, conditions)
	if err != nil {
		return nil, 0, err
	}

	// 写入缓存
	if err := s.postCache.SetPostList(ctx, cacheKey, posts); err != nil {
		return nil, 0, err
	}

	return posts, total, nil
}

func (s *postService) IncrementViewCount(ctx context.Context, id uint) error {
	// 增加缓存中的计数
	count, err := s.postCache.IncrViewCount(ctx, id)
	if err != nil {
		return err
	}

	// 定期同步到数据库
	if count%10 == 0 { // 每10次访问同步一次
		if err := s.postRepo.IncrementViewCount(id); err != nil {
			return err
		}
	}

	return nil
}

func (s *postService) ListPostsByCategory(ctx context.Context, categoryID uint, page, pageSize int) ([]models.Post, int64, error) {
	return s.postRepo.ListByCategoryID(categoryID, page, pageSize)
}

func (s *postService) ListPostsByTag(ctx context.Context, tagID uint, page, pageSize int) ([]models.Post, int64, error) {
	return s.postRepo.ListByTagID(tagID, page, pageSize)
}

func (s *postService) ListPostsByUser(ctx context.Context, userID uint, page, pageSize int) ([]models.Post, int64, error) {
	return s.postRepo.ListByUserID(userID, page, pageSize)
}
