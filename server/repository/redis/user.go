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
	userKeyPrefix = "user:"
	userExpiration = 24 * time.Hour
)

// UserCache 用户缓存接口
type UserCache interface {
	Set(ctx context.Context, user *models.User) error
	Get(ctx context.Context, id uint) (*models.User, error)
	Delete(ctx context.Context, id uint) error
	SetUserToken(ctx context.Context, userID uint, token string) error
	GetUserToken(ctx context.Context, userID uint) (string, error)
}

type userCache struct {
	client *redis.Client
}

// NewUserCache 创建用户缓存实例
func NewUserCache(client *redis.Client) UserCache {
	return &userCache{client: client}
}

func (c *userCache) Set(ctx context.Context, user *models.User) error {
	key := fmt.Sprintf("%s%d", userKeyPrefix, user.ID)
	data, err := json.Marshal(user)
	if err != nil {
		return err
	}
	return c.client.Set(ctx, key, data, userExpiration).Err()
}

func (c *userCache) Get(ctx context.Context, id uint) (*models.User, error) {
	key := fmt.Sprintf("%s%d", userKeyPrefix, id)
	data, err := c.client.Get(ctx, key).Bytes()
	if err != nil {
		if err == redis.Nil {
			return nil, nil
		}
		return nil, err
	}

	var user models.User
	if err := json.Unmarshal(data, &user); err != nil {
		return nil, err
	}
	return &user, nil
}

func (c *userCache) Delete(ctx context.Context, id uint) error {
	key := fmt.Sprintf("%s%d", userKeyPrefix, id)
	return c.client.Del(ctx, key).Err()
}

func (c *userCache) SetUserToken(ctx context.Context, userID uint, token string) error {
	key := fmt.Sprintf("%stoken:%d", userKeyPrefix, userID)
	return c.client.Set(ctx, key, token, 7*24*time.Hour).Err()
}

func (c *userCache) GetUserToken(ctx context.Context, userID uint) (string, error) {
	key := fmt.Sprintf("%stoken:%d", userKeyPrefix, userID)
	return c.client.Get(ctx, key).Result()
}
