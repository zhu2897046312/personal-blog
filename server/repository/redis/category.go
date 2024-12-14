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
	categoryKeyPrefix = "category:"
	categoryListKey = "category:list"
	categoryExpiration = 12 * time.Hour
)

// CategoryCache 分类缓存接口
type CategoryCache interface {
	Set(ctx context.Context, category *models.Category) error
	Get(ctx context.Context, id uint) (*models.Category, error)
	Delete(ctx context.Context, id uint) error
	SetList(ctx context.Context, categories []models.Category) error
	GetList(ctx context.Context) ([]models.Category, error)
}

type categoryCache struct {
	client *redis.Client
}

// NewCategoryCache 创建分类缓存实例
func NewCategoryCache(client *redis.Client) CategoryCache {
	return &categoryCache{client: client}
}

func (c *categoryCache) Set(ctx context.Context, category *models.Category) error {
	key := fmt.Sprintf("%s%d", categoryKeyPrefix, category.ID)
	data, err := json.Marshal(category)
	if err != nil {
		return err
	}
	return c.client.Set(ctx, key, data, categoryExpiration).Err()
}

func (c *categoryCache) Get(ctx context.Context, id uint) (*models.Category, error) {
	key := fmt.Sprintf("%s%d", categoryKeyPrefix, id)
	data, err := c.client.Get(ctx, key).Bytes()
	if err != nil {
		if err == redis.Nil {
			return nil, nil
		}
		return nil, err
	}

	var category models.Category
	if err := json.Unmarshal(data, &category); err != nil {
		return nil, err
	}
	return &category, nil
}

func (c *categoryCache) Delete(ctx context.Context, id uint) error {
	key := fmt.Sprintf("%s%d", categoryKeyPrefix, id)
	return c.client.Del(ctx, key).Err()
}

func (c *categoryCache) SetList(ctx context.Context, categories []models.Category) error {
	data, err := json.Marshal(categories)
	if err != nil {
		return err
	}
	return c.client.Set(ctx, categoryListKey, data, categoryExpiration).Err()
}

func (c *categoryCache) GetList(ctx context.Context) ([]models.Category, error) {
	data, err := c.client.Get(ctx, categoryListKey).Bytes()
	if err != nil {
		if err == redis.Nil {
			return nil, nil
		}
		return nil, err
	}

	var categories []models.Category
	if err := json.Unmarshal(data, &categories); err != nil {
		return nil, err
	}
	return categories, nil
}
