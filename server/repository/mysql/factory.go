package mysql

import (
	"sync"

	"gorm.io/gorm"
)

var (
	factoryInstance *factory
	once           sync.Once
)

// Factory 定义仓库工厂接口
type Factory interface {
	GetUserRepository() UserRepository
	GetPostRepository() PostRepository
	GetCategoryRepository() CategoryRepository
	GetTagRepository() TagRepository
	GetCommentRepository() CommentRepository
}

// factory 实现Factory接口
type factory struct {
	db          *gorm.DB
	userRepo    UserRepository
	postRepo    PostRepository
	categoryRepo CategoryRepository
	tagRepo     TagRepository
	commentRepo CommentRepository
	mu          sync.RWMutex
}

// NewFactory 创建工厂实例（单例）
func NewFactory(db *gorm.DB) Factory {
	once.Do(func() {
		factoryInstance = &factory{
			db: db,
		}
	})
	return factoryInstance
}

func (f *factory) GetUserRepository() UserRepository {
	f.mu.RLock()
	if f.userRepo != nil {
		defer f.mu.RUnlock()
		return f.userRepo
	}
	f.mu.RUnlock()

	f.mu.Lock()
	defer f.mu.Unlock()
	if f.userRepo == nil {
		f.userRepo = NewUserRepository(f.db)
	}
	return f.userRepo
}

func (f *factory) GetPostRepository() PostRepository {
	f.mu.RLock()
	if f.postRepo != nil {
		defer f.mu.RUnlock()
		return f.postRepo
	}
	f.mu.RUnlock()

	f.mu.Lock()
	defer f.mu.Unlock()
	if f.postRepo == nil {
		f.postRepo = NewPostRepository(f.db)
	}
	return f.postRepo
}

func (f *factory) GetCategoryRepository() CategoryRepository {
	f.mu.RLock()
	if f.categoryRepo != nil {
		defer f.mu.RUnlock()
		return f.categoryRepo
	}
	f.mu.RUnlock()

	f.mu.Lock()
	defer f.mu.Unlock()
	if f.categoryRepo == nil {
		f.categoryRepo = NewCategoryRepository(f.db)
	}
	return f.categoryRepo
}

func (f *factory) GetTagRepository() TagRepository {
	f.mu.RLock()
	if f.tagRepo != nil {
		defer f.mu.RUnlock()
		return f.tagRepo
	}
	f.mu.RUnlock()

	f.mu.Lock()
	defer f.mu.Unlock()
	if f.tagRepo == nil {
		f.tagRepo = NewTagRepository(f.db)
	}
	return f.tagRepo
}

func (f *factory) GetCommentRepository() CommentRepository {
	f.mu.RLock()
	if f.commentRepo != nil {
		defer f.mu.RUnlock()
		return f.commentRepo
	}
	f.mu.RUnlock()

	f.mu.Lock()
	defer f.mu.Unlock()
	if f.commentRepo == nil {
		f.commentRepo = NewCommentRepository(f.db)
	}
	return f.commentRepo
}
