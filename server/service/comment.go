package service

import (
	"context"
	"time"

	"github.com/personal-blog/models"
	"github.com/personal-blog/repository/mysql"
	"github.com/personal-blog/repository/redis"
)

// CommentService 评论服务接口
type CommentService interface {
	CreateComment(ctx context.Context, comment *models.Comment) error
	UpdateComment(ctx context.Context, comment *models.Comment) error
	DeleteComment(ctx context.Context, id uint) error
	GetCommentByID(ctx context.Context, id uint) (*models.Comment, error)
	ListCommentsByPost(ctx context.Context, postID uint, page, pageSize int) ([]models.Comment, int64, error)
	ListCommentsByUser(ctx context.Context, userID uint, page, pageSize int) ([]models.Comment, int64, error)
}

type commentService struct {
	commentRepo  mysql.CommentRepository
	commentCache redis.CommentCache
}

// NewCommentService 创建评论服务实例
func NewCommentService(commentRepo mysql.CommentRepository, commentCache redis.CommentCache) CommentService {
	return &commentService{
		commentRepo:  commentRepo,
		commentCache: commentCache,
	}
}

func (s *commentService) CreateComment(ctx context.Context, comment *models.Comment) error {
	comment.CreatedAt = time.Now()
	comment.UpdatedAt = time.Now()

	// 创建评论
	if err := s.commentRepo.Create(comment); err != nil {
		return err
	}

	// 写入缓存
	if err := s.commentCache.Set(ctx, comment); err != nil {
		return err
	}

	// 清除文章评论列表缓存
	return s.commentCache.Delete(ctx, comment.PostID)
}

func (s *commentService) UpdateComment(ctx context.Context, comment *models.Comment) error {
	comment.UpdatedAt = time.Now()

	// 更新评论
	if err := s.commentRepo.Update(comment); err != nil {
		return err
	}

	// 更新缓存
	if err := s.commentCache.Set(ctx, comment); err != nil {
		return err
	}

	// 清除文章评论列表缓存
	return s.commentCache.Delete(ctx, comment.PostID)
}

func (s *commentService) DeleteComment(ctx context.Context, id uint) error {
	// 获取评论信息（用于后续清除缓存）
	comment, err := s.commentRepo.FindByID(id)
	if err != nil {
		return err
	}

	// 删除评论
	if err := s.commentRepo.Delete(id); err != nil {
		return err
	}

	// 删除缓存
	if err := s.commentCache.Delete(ctx, id); err != nil {
		return err
	}

	// 清除文章评论列表缓存
	return s.commentCache.Delete(ctx, comment.PostID)
}

func (s *commentService) GetCommentByID(ctx context.Context, id uint) (*models.Comment, error) {
	// 先从缓存获取
	comment, err := s.commentCache.Get(ctx, id)
	if err != nil {
		return nil, err
	}
	if comment != nil {
		return comment, nil
	}

	// 缓存未命中，从数据库获取
	comment, err = s.commentRepo.FindByID(id)
	if err != nil {
		return nil, err
	}

	// 写入缓存
	if err := s.commentCache.Set(ctx, comment); err != nil {
		return nil, err
	}

	return comment, nil
}

func (s *commentService) ListCommentsByPost(ctx context.Context, postID uint, page, pageSize int) ([]models.Comment, int64, error) {
	// 先从缓存获取
	comments, err := s.commentCache.GetPostComments(ctx, postID)
	if err != nil {
		return nil, 0, err
	}
	if len(comments) > 0 {
		// 处理分页
		start := (page - 1) * pageSize
		end := start + pageSize
		if start >= len(comments) {
			return []models.Comment{}, int64(len(comments)), nil
		}
		if end > len(comments) {
			end = len(comments)
		}
		return comments[start:end], int64(len(comments)), nil
	}

	// 从数据库获取
	comments, total, err := s.commentRepo.ListByPostID(postID, page, pageSize)
	if err != nil {
		return nil, 0, err
	}

	// 写入缓存
	if err := s.commentCache.SetPostComments(ctx, postID, comments); err != nil {
		return nil, 0, err
	}

	return comments, total, nil
}

func (s *commentService) ListCommentsByUser(ctx context.Context, userID uint, page, pageSize int) ([]models.Comment, int64, error) {
	// 先从缓存获取
	comments, err := s.commentCache.GetUserComments(ctx, userID)
	if err != nil {
		return nil, 0, err
	}
	if len(comments) > 0 {
		// 处理分页
		start := (page - 1) * pageSize
		end := start + pageSize
		if start >= len(comments) {
			return []models.Comment{}, int64(len(comments)), nil
		}
		if end > len(comments) {
			end = len(comments)
		}
		return comments[start:end], int64(len(comments)), nil
	}

	// 从数据库获取
	comments, total, err := s.commentRepo.ListByUserID(userID, page, pageSize)
	if err != nil {
		return nil, 0, err
	}

	// 写入缓存
	if err := s.commentCache.SetUserComments(ctx, userID, comments); err != nil {
		return nil, 0, err
	}

	return comments, total, nil
}
