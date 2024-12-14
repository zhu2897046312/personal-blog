package models

import (
	"time"
)

// Comment 评论模型
type Comment struct {
	ID        uint       `gorm:"primarykey" json:"id"`
	Content   string     `gorm:"size:1000;not null" json:"content"`
	PostID    uint       `json:"post_id"`
	Post      Post       `json:"post"`
	UserID    uint       `json:"user_id"`
	User      User       `json:"user"`
	ParentID  *uint      `json:"parent_id"`           // 父评论ID，用于回复功能
	Parent    *Comment   `json:"parent"`              // 父评论
	Children  []Comment  `gorm:"foreignkey:ParentID"` // 子评论
	Status    int        `gorm:"default:1" json:"status"` // 1:正常 0:待审核 -1:已删除
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `gorm:"index" json:"-"`
}
