package redis

import (
	"sync"

	"github.com/redis/go-redis/v9"
)

var (
	factoryInstance *factory
	once           sync.Once
)

// Factory 定义Redis仓库工厂接口
type Factory interface {
	GetUserCache() UserCache
	GetPostCache() PostCache
	GetCategoryCache() CategoryCache
	GetTagCache() TagCache
	GetCommentCache() CommentCache
}

// factory 实现Factory接口
type factory struct {
	client       *redis.Client
	userCache    UserCache
	postCache    PostCache
	categoryCache CategoryCache
	tagCache     TagCache
	commentCache CommentCache
	mu           sync.RWMutex
}

// NewFactory 创建Redis工厂实例（单例）
func NewFactory(client *redis.Client) Factory {
	once.Do(func() {
		factoryInstance = &factory{
			client: client,
		}
	})
	return factoryInstance
}

func (f *factory) GetUserCache() UserCache {
	f.mu.RLock()
	if f.userCache != nil {
		defer f.mu.RUnlock()
		return f.userCache
	}
	f.mu.RUnlock()

	f.mu.Lock()
	defer f.mu.Unlock()
	if f.userCache == nil {
		f.userCache = NewUserCache(f.client)
	}
	return f.userCache
}

func (f *factory) GetPostCache() PostCache {
	f.mu.RLock()
	if f.postCache != nil {
		defer f.mu.RUnlock()
		return f.postCache
	}
	f.mu.RUnlock()

	f.mu.Lock()
	defer f.mu.Unlock()
	if f.postCache == nil {
		f.postCache = NewPostCache(f.client)
	}
	return f.postCache
}

func (f *factory) GetCategoryCache() CategoryCache {
	f.mu.RLock()
	if f.categoryCache != nil {
		defer f.mu.RUnlock()
		return f.categoryCache
	}
	f.mu.RUnlock()

	f.mu.Lock()
	defer f.mu.Unlock()
	if f.categoryCache == nil {
		f.categoryCache = NewCategoryCache(f.client)
	}
	return f.categoryCache
}

func (f *factory) GetTagCache() TagCache {
	f.mu.RLock()
	if f.tagCache != nil {
		defer f.mu.RUnlock()
		return f.tagCache
	}
	f.mu.RUnlock()

	f.mu.Lock()
	defer f.mu.Unlock()
	if f.tagCache == nil {
		f.tagCache = NewTagCache(f.client)
	}
	return f.tagCache
}

func (f *factory) GetCommentCache() CommentCache {
	f.mu.RLock()
	if f.commentCache != nil {
		defer f.mu.RUnlock()
		return f.commentCache
	}
	f.mu.RUnlock()

	f.mu.Lock()
	defer f.mu.Unlock()
	if f.commentCache == nil {
		f.commentCache = NewCommentCache(f.client)
	}
	return f.commentCache
}
