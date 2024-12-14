package models

import (
	"time"
)

// Tag 标签模型
type Tag struct {
	ID        uint      `gorm:"primarykey" json:"id"`
	Name      string    `gorm:"size:50;not null;unique" json:"name"`
	Posts     []Post    `gorm:"many2many:post_tags;" json:"posts"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
