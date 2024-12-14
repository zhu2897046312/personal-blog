package models

import (
	"time"
)

// Post 文章模型
type Post struct {
	ID         uint       `gorm:"primarykey" json:"id"`
	Title      string     `gorm:"size:200;not null" json:"title"`
	Content    string     `gorm:"type:text" json:"content"`
	Summary    string     `gorm:"size:500" json:"summary"`
	Cover      string     `gorm:"size:255" json:"cover"`
	Status     int        `gorm:"default:1" json:"status"`     // 1:已发布 0:草稿 -1:已删除
	IsTop      bool       `gorm:"default:false" json:"is_top"` // 是否置顶
	ViewCount  int64      `gorm:"default:0" json:"view_count"` // 浏览量
	UserID     uint       `json:"user_id"`                     // 作者ID
	User       User       `json:"user"`
	CategoryID uint       `json:"category_id"`
	Category   Category   `json:"category"`
	Tags       []Tag      `gorm:"many2many:post_tags;" json:"tags"`
	Comments   []Comment  `json:"comments"`
	CreatedAt  time.Time  `json:"created_at"`
	UpdatedAt  time.Time  `json:"updated_at"`
	DeletedAt  *time.Time `gorm:"index" json:"-"`
}
