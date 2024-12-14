package service

import (
	"sync"

	"github.com/personal-blog/repository/mysql"
	"github.com/personal-blog/repository/redis"
)

var (
	factoryInstance *factory
	once           sync.Once
)

// Factory 定义服务工厂接口
type Factory interface {
	GetUserService() UserService
	GetPostService() PostService
	GetCategoryService() CategoryService
	GetTagService() TagService
	GetCommentService() CommentService
}

// factory 实现Factory接口
type factory struct {
	mysqlFactory mysql.Factory
	redisFactory redis.Factory
	userSrv      UserService
	postSrv      PostService
	categorySrv  CategoryService
	tagSrv       TagService
	commentSrv   CommentService
	mu           sync.RWMutex
}

// NewFactory 创建服务工厂实例（单例）
func NewFactory(mysqlFactory mysql.Factory, redisFactory redis.Factory) Factory {
	once.Do(func() {
		factoryInstance = &factory{
			mysqlFactory: mysqlFactory,
			redisFactory: redisFactory,
		}
	})
	return factoryInstance
}

func (f *factory) GetUserService() UserService {
	f.mu.RLock()
	if f.userSrv != nil {
		defer f.mu.RUnlock()
		return f.userSrv
	}
	f.mu.RUnlock()

	f.mu.Lock()
	defer f.mu.Unlock()
	if f.userSrv == nil {
		f.userSrv = NewUserService(f.mysqlFactory.GetUserRepository(), f.redisFactory.GetUserCache())
	}
	return f.userSrv
}

func (f *factory) GetPostService() PostService {
	f.mu.RLock()
	if f.postSrv != nil {
		defer f.mu.RUnlock()
		return f.postSrv
	}
	f.mu.RUnlock()

	f.mu.Lock()
	defer f.mu.Unlock()
	if f.postSrv == nil {
		f.postSrv = NewPostService(
			f.mysqlFactory.GetPostRepository(),
			f.redisFactory.GetPostCache(),
			f.mysqlFactory.GetTagRepository(),
			f.mysqlFactory.GetCategoryRepository(),
		)
	}
	return f.postSrv
}

func (f *factory) GetCategoryService() CategoryService {
	f.mu.RLock()
	if f.categorySrv != nil {
		defer f.mu.RUnlock()
		return f.categorySrv
	}
	f.mu.RUnlock()

	f.mu.Lock()
	defer f.mu.Unlock()
	if f.categorySrv == nil {
		f.categorySrv = NewCategoryService(f.mysqlFactory.GetCategoryRepository(), f.redisFactory.GetCategoryCache())
	}
	return f.categorySrv
}

func (f *factory) GetTagService() TagService {
	f.mu.RLock()
	if f.tagSrv != nil {
		defer f.mu.RUnlock()
		return f.tagSrv
	}
	f.mu.RUnlock()

	f.mu.Lock()
	defer f.mu.Unlock()
	if f.tagSrv == nil {
		f.tagSrv = NewTagService(
			f.mysqlFactory.GetTagRepository(),
			f.redisFactory.GetTagCache(),
			f.mysqlFactory.GetPostRepository(),
		)
	}
	return f.tagSrv
}

func (f *factory) GetCommentService() CommentService {
	f.mu.RLock()
	if f.commentSrv != nil {
		defer f.mu.RUnlock()
		return f.commentSrv
	}
	f.mu.RUnlock()

	f.mu.Lock()
	defer f.mu.Unlock()
	if f.commentSrv == nil {
		f.commentSrv = NewCommentService(f.mysqlFactory.GetCommentRepository(), f.redisFactory.GetCommentCache())
	}
	return f.commentSrv
}
