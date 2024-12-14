package models

import (
	"time"
)

// Category 分类模型
type Category struct {
	ID          uint      `gorm:"primarykey" json:"id"`
	Name        string    `gorm:"size:50;not null;unique" json:"name"`
	Description string    `gorm:"size:200" json:"description"`
	Posts       []Post    `json:"posts"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
