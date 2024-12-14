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
	tagKeyPrefix = "tag:"
	tagListKey = "tag:list"
	tagExpiration = 12 * time.Hour
)

// TagCache 标签缓存接口
type TagCache interface {
	Set(ctx context.Context, tag *models.Tag) error
	Get(ctx context.Context, id uint) (*models.Tag, error)
	Delete(ctx context.Context, id uint) error
	SetList(ctx context.Context, tags []models.Tag) error
	GetList(ctx context.Context) ([]models.Tag, error)
	SetPostTags(ctx context.Context, postID uint, tags []models.Tag) error
	GetPostTags(ctx context.Context, postID uint) ([]models.Tag, error)
}

type tagCache struct {
	client *redis.Client
}

// NewTagCache 创建标签缓存实例
func NewTagCache(client *redis.Client) TagCache {
	return &tagCache{client: client}
}

func (c *tagCache) Set(ctx context.Context, tag *models.Tag) error {
	key := fmt.Sprintf("%s%d", tagKeyPrefix, tag.ID)
	data, err := json.Marshal(tag)
	if err != nil {
		return err
	}
	return c.client.Set(ctx, key, data, tagExpiration).Err()
}

func (c *tagCache) Get(ctx context.Context, id uint) (*models.Tag, error) {
	key := fmt.Sprintf("%s%d", tagKeyPrefix, id)
	data, err := c.client.Get(ctx, key).Bytes()
	if err != nil {
		if err == redis.Nil {
			return nil, nil
		}
		return nil, err
	}

	var tag models.Tag
	if err := json.Unmarshal(data, &tag); err != nil {
		return nil, err
	}
	return &tag, nil
}

func (c *tagCache) Delete(ctx context.Context, id uint) error {
	key := fmt.Sprintf("%s%d", tagKeyPrefix, id)
	return c.client.Del(ctx, key).Err()
}

func (c *tagCache) SetList(ctx context.Context, tags []models.Tag) error {
	data, err := json.Marshal(tags)
	if err != nil {
		return err
	}
	return c.client.Set(ctx, tagListKey, data, tagExpiration).Err()
}

func (c *tagCache) GetList(ctx context.Context) ([]models.Tag, error) {
	data, err := c.client.Get(ctx, tagListKey).Bytes()
	if err != nil {
		if err == redis.Nil {
			return nil, nil
		}
		return nil, err
	}

	var tags []models.Tag
	if err := json.Unmarshal(data, &tags); err != nil {
		return nil, err
	}
	return tags, nil
}

func (c *tagCache) SetPostTags(ctx context.Context, postID uint, tags []models.Tag) error {
	key := fmt.Sprintf("%spost:%d", tagKeyPrefix, postID)
	data, err := json.Marshal(tags)
	if err != nil {
		return err
	}
	return c.client.Set(ctx, key, data, tagExpiration).Err()
}

func (c *tagCache) GetPostTags(ctx context.Context, postID uint) ([]models.Tag, error) {
	key := fmt.Sprintf("%spost:%d", tagKeyPrefix, postID)
	data, err := c.client.Get(ctx, key).Bytes()
	if err != nil {
		if err == redis.Nil {
			return nil, nil
		}
		return nil, err
	}

	var tags []models.Tag
	if err := json.Unmarshal(data, &tags); err != nil {
		return nil, err
	}
	return tags, nil
}
