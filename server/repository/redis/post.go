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
	postKeyPrefix = "post:"
	postExpiration = 1 * time.Hour
	postViewCountPrefix = "post:view:"
)

// PostCache 文章缓存接口
type PostCache interface {
	Set(ctx context.Context, post *models.Post) error
	Get(ctx context.Context, id uint) (*models.Post, error)
	Delete(ctx context.Context, id uint) error
	IncrViewCount(ctx context.Context, id uint) (int64, error)
	GetViewCount(ctx context.Context, id uint) (int64, error)
	SetPostList(ctx context.Context, key string, posts []models.Post) error
	GetPostList(ctx context.Context, key string) ([]models.Post, error)
}

type postCache struct {
	client *redis.Client
}

// NewPostCache 创建文章缓存实例
func NewPostCache(client *redis.Client) PostCache {
	return &postCache{client: client}
}

func (c *postCache) Set(ctx context.Context, post *models.Post) error {
	key := fmt.Sprintf("%s%d", postKeyPrefix, post.ID)
	data, err := json.Marshal(post)
	if err != nil {
		return err
	}
	return c.client.Set(ctx, key, data, postExpiration).Err()
}

func (c *postCache) Get(ctx context.Context, id uint) (*models.Post, error) {
	key := fmt.Sprintf("%s%d", postKeyPrefix, id)
	data, err := c.client.Get(ctx, key).Bytes()
	if err != nil {
		if err == redis.Nil {
			return nil, nil
		}
		return nil, err
	}

	var post models.Post
	if err := json.Unmarshal(data, &post); err != nil {
		return nil, err
	}
	return &post, nil
}

func (c *postCache) Delete(ctx context.Context, id uint) error {
	key := fmt.Sprintf("%s%d", postKeyPrefix, id)
	return c.client.Del(ctx, key).Err()
}

func (c *postCache) IncrViewCount(ctx context.Context, id uint) (int64, error) {
	key := fmt.Sprintf("%s%d", postViewCountPrefix, id)
	return c.client.Incr(ctx, key).Result()
}

func (c *postCache) GetViewCount(ctx context.Context, id uint) (int64, error) {
	key := fmt.Sprintf("%s%d", postViewCountPrefix, id)
	count, err := c.client.Get(ctx, key).Int64()
	if err == redis.Nil {
		return 0, nil
	}
	return count, err
}

func (c *postCache) SetPostList(ctx context.Context, key string, posts []models.Post) error {
	data, err := json.Marshal(posts)
	if err != nil {
		return err
	}
	return c.client.Set(ctx, key, data, postExpiration).Err()
}

func (c *postCache) GetPostList(ctx context.Context, key string) ([]models.Post, error) {
	data, err := c.client.Get(ctx, key).Bytes()
	if err != nil {
		if err == redis.Nil {
			return nil, nil
		}
		return nil, err
	}

	var posts []models.Post
	if err := json.Unmarshal(data, &posts); err != nil {
		return nil, err
	}
	return posts, nil
}
