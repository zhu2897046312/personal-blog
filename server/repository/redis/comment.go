package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/personal-blog/models"
)

const (
	commentKeyPrefix = "comment:"
	commentExpiration = 6 * time.Hour
)

// CommentCache 评论缓存接口
type CommentCache interface {
	Set(ctx context.Context, comment *models.Comment) error
	Get(ctx context.Context, id uint) (*models.Comment, error)
	Delete(ctx context.Context, id uint) error
	SetPostComments(ctx context.Context, postID uint, comments []models.Comment) error
	GetPostComments(ctx context.Context, postID uint) ([]models.Comment, error)
	SetUserComments(ctx context.Context, userID uint, comments []models.Comment) error
	GetUserComments(ctx context.Context, userID uint) ([]models.Comment, error)
}

type commentCache struct {
	client *redis.Client
}

// NewCommentCache 创建评论缓存实例
func NewCommentCache(client *redis.Client) CommentCache {
	return &commentCache{client: client}
}

func (c *commentCache) Set(ctx context.Context, comment *models.Comment) error {
	key := fmt.Sprintf("%s%d", commentKeyPrefix, comment.ID)
	data, err := json.Marshal(comment)
	if err != nil {
		return err
	}
	return c.client.Set(ctx, key, data, commentExpiration).Err()
}

func (c *commentCache) Get(ctx context.Context, id uint) (*models.Comment, error) {
	key := fmt.Sprintf("%s%d", commentKeyPrefix, id)
	data, err := c.client.Get(ctx, key).Bytes()
	if err != nil {
		if err == redis.Nil {
			return nil, nil
		}
		return nil, err
	}

	var comment models.Comment
	if err := json.Unmarshal(data, &comment); err != nil {
		return nil, err
	}
	return &comment, nil
}

func (c *commentCache) Delete(ctx context.Context, id uint) error {
	key := fmt.Sprintf("%s%d", commentKeyPrefix, id)
	return c.client.Del(ctx, key).Err()
}

func (c *commentCache) SetPostComments(ctx context.Context, postID uint, comments []models.Comment) error {
	key := fmt.Sprintf("%spost:%d", commentKeyPrefix, postID)
	data, err := json.Marshal(comments)
	if err != nil {
		return err
	}
	return c.client.Set(ctx, key, data, commentExpiration).Err()
}

func (c *commentCache) GetPostComments(ctx context.Context, postID uint) ([]models.Comment, error) {
	key := fmt.Sprintf("%spost:%d", commentKeyPrefix, postID)
	data, err := c.client.Get(ctx, key).Bytes()
	if err != nil {
		if err == redis.Nil {
			return nil, nil
		}
		return nil, err
	}

	var comments []models.Comment
	if err := json.Unmarshal(data, &comments); err != nil {
		return nil, err
	}
	return comments, nil
}

func (c *commentCache) SetUserComments(ctx context.Context, userID uint, comments []models.Comment) error {
	key := fmt.Sprintf("%suser:%d", commentKeyPrefix, userID)
	data, err := json.Marshal(comments)
	if err != nil {
		return err
	}
	return c.client.Set(ctx, key, data, commentExpiration).Err()
}

func (c *commentCache) GetUserComments(ctx context.Context, userID uint) ([]models.Comment, error) {
	key := fmt.Sprintf("%suser:%d", commentKeyPrefix, userID)
	data, err := c.client.Get(ctx, key).Bytes()
	if err != nil {
		if err == redis.Nil {
			return nil, nil
		}
		return nil, err
	}

	var comments []models.Comment
	if err := json.Unmarshal(data, &comments); err != nil {
		return nil, err
	}
	return comments, nil
}
