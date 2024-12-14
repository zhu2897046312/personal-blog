package models

import (
	"time"
)

// User 用户模型
type User struct {
	ID        uint      `gorm:"primarykey" json:"id"`
	Username  string    `gorm:"size:50;not null;unique" json:"username"`
	Password  string    `gorm:"size:100;not null" json:"-"`  // 密码不返回给前端
	Email     string    `gorm:"size:100;not null;unique" json:"email"`
	Nickname  string    `gorm:"size:50" json:"nickname"`
	Avatar    string    `gorm:"size:255" json:"avatar"`
	Role      string    `gorm:"size:20;default:'user'" json:"role"` // admin/user
	Bio       string    `gorm:"size:500" json:"bio"`
	Status    int       `gorm:"default:1" json:"status"`           // 1:正常 0:禁用
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
