package docs

import (
	"github.com/swaggo/swag"
	"time"
)

// Comment 评论模型
type Comment struct {
	ID        uint      `json:"id" example:"1"`
	PostID    uint      `json:"post_id" example:"1"`
	UserID    uint      `json:"user_id" example:"1"`
	Content   string    `json:"content" example:"这是一条评论"`
	ParentID  *uint     `json:"parent_id,omitempty" example:"0"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	User      *User     `json:"user,omitempty"`
}

// @title Personal Blog API
// @version 1.0
// @description 个人博客系统API文档
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:8080
// @BasePath /api/v1
// @schemes http https

// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization

// @x-enum-varnames Status
// @x-enum-descriptions 文章状态
// @x-enum-values 1:公开 2:草稿

// @x-enum-varnames Role
// @x-enum-descriptions 用户角色
// @x-enum-values admin:管理员 user:普通用户

// @tag.name user
// @tag.description 用户相关接口

// @tag.name post
// @tag.description 文章相关接口

// @tag.name category
// @tag.description 分类相关接口

// @tag.name tag
// @tag.description 标签相关接口

// Response 通用响应
type Response struct {
	Code int         `json:"code" example:"200"`    // HTTP状态码
	Msg  string      `json:"msg" example:"success"` // 响应消息
	Data interface{} `json:"data"`                  // 响应数据
}

// PaginationData 分页数据
type PaginationData struct {
	Total    int64       `json:"total" example:"100"`     // 总数
	Page     int         `json:"page" example:"1"`        // 当前页码
	PageSize int         `json:"page_size" example:"10"`  // 每页数量
	Items    interface{} `json:"items"`                   // 数据列表
}

// RegisterRequest 用户注册请求
type RegisterRequest struct {
	Username string `json:"username" binding:"required,min=3,max=32" example:"johndoe"`     // 用户名
	Password string `json:"password" binding:"required,min=6,max=32" example:"password123"` // 密码
	Email    string `json:"email" binding:"required,email" example:"john@example.com"`      // 邮箱
	Nickname string `json:"nickname" binding:"required,min=2,max=32" example:"John Doe"`    // 昵称
}

// LoginRequest 用户登录请求
type LoginRequest struct {
	Username string `json:"username" binding:"required" example:"johndoe"`     // 用户名
	Password string `json:"password" binding:"required" example:"password123"` // 密码
}

// LoginResponse 登录响应
type LoginResponse struct {
	User  *User  `json:"user"`                                             // 用户信息
	Token string `json:"token" example:"eyJhbGciOiJIUzI1NiIsInR5cCI..."` // JWT Token
}

// UpdateProfileRequest 更新用户信息请求
type UpdateProfileRequest struct {
	Nickname string `json:"nickname" binding:"required,min=2,max=32" example:"John Doe"`      // 昵称
	Email    string `json:"email" binding:"required,email" example:"john@example.com"`        // 邮箱
	Avatar   string `json:"avatar" example:"https://example.com/avatar/johndoe.jpg"`         // 头像
}

// ChangePasswordRequest 修改密码请求
type ChangePasswordRequest struct {
	OldPassword string `json:"old_password" binding:"required" example:"oldpass123"`           // 旧密码
	NewPassword string `json:"new_password" binding:"required,min=6,max=32" example:"newpass123"` // 新密码
}

// ListUsersRequest 用户列表请求
type ListUsersRequest struct {
	Page     int `form:"page" binding:"required,min=1" example:"1"`       // 页码
	PageSize int `form:"page_size" binding:"required,min=1,max=100" example:"10"` // 每页数量
}

// CreatePostRequest 创建文章请求
type CreatePostRequest struct {
	Title      string   `json:"title" binding:"required,min=1,max=100" example:"My First Blog Post"` // 标题
	Content    string   `json:"content" binding:"required,min=1" example:"This is my first blog post content..."` // 内容
	CategoryID uint     `json:"category_id" binding:"required" example:"1"` // 分类ID
	Tags       []string `json:"tags" example:"['tech','golang']"`          // 标签列表
	Status     int      `json:"status" binding:"required,oneof=1 2" example:"1"` // 状态：1公开 2草稿
}

// UpdatePostRequest 更新文章请求
type UpdatePostRequest struct {
	Title      string   `json:"title" binding:"required,min=1,max=100" example:"Updated Blog Post"` // 标题
	Content    string   `json:"content" binding:"required,min=1" example:"Updated content..."` // 内容
	CategoryID uint     `json:"category_id" binding:"required" example:"1"` // 分类ID
	Tags       []string `json:"tags" example:"['tech','golang']"`          // 标签列表
	Status     int      `json:"status" binding:"required,oneof=1 2" example:"1"` // 状态：1公开 2草稿
}

// ListPostsRequest 文章列表请求
type ListPostsRequest struct {
	Page       int    `form:"page" binding:"required,min=1" example:"1"`       // 页码
	PageSize   int    `form:"page_size" binding:"required,min=1,max=100" example:"10"` // 每页数量
	CategoryID uint   `form:"category_id" example:"1"`                         // 分类ID
	TagID      uint   `form:"tag_id" example:"1"`                             // 标签ID
	UserID     uint   `form:"user_id" example:"1"`                            // 作者ID
	Keyword    string `form:"keyword" example:"golang"`                       // 关键词
	Status     int    `form:"status" example:"1"`                             // 状态：1公开 2草稿
}

// CreateCategoryRequest 创建分类请求
type CreateCategoryRequest struct {
	Name        string `json:"name" binding:"required,min=1,max=50" example:"Technology"` // 名称
	Description string `json:"description" binding:"omitempty,max=200" example:"Tech related posts"` // 描述
}

// UpdateCategoryRequest 更新分类请求
type UpdateCategoryRequest struct {
	Name        string `json:"name" binding:"required,min=1,max=50" example:"Technology"` // 名称
	Description string `json:"description" binding:"omitempty,max=200" example:"Tech related posts"` // 描述
}

// CreateTagRequest 创建标签请求
type CreateTagRequest struct {
	Name string `json:"name" binding:"required,min=1,max=50" example:"golang"` // 名称
}

// UpdateTagRequest 更新标签请求
type UpdateTagRequest struct {
	Name string `json:"name" binding:"required,min=1,max=50" example:"golang"` // 名称
}

// CreateTagsRequest 批量创建标签请求
type CreateTagsRequest struct {
	Names []string `json:"names" binding:"required,dive,min=1,max=50" example:"['golang','tech']"` // 标签名称列表
}

// User 用户信息
type User struct {
	ID        uint   `json:"id" example:"1"`                     // 用户ID
	Username  string `json:"username" example:"johndoe"`         // 用户名
	Nickname  string `json:"nickname" example:"John Doe"`        // 昵称
	Email     string `json:"email" example:"john@example.com"`   // 邮箱
	Avatar    string `json:"avatar" example:"https://example.com/avatar/johndoe.jpg"` // 头像
	Role      string `json:"role" example:"user"`                // 角色：admin管理员 user普通用户
	CreatedAt string `json:"created_at" example:"2024-01-01 00:00:00"` // 创建时间
	UpdatedAt string `json:"updated_at" example:"2024-01-01 00:00:00"` // 更新时间
}

// Post 文章信息
type Post struct {
	ID          uint      `json:"id" example:"1"`                    // 文章ID
	Title       string    `json:"title" example:"My First Blog Post"` // 标题
	Content     string    `json:"content" example:"This is my first blog post content..."` // 内容
	UserID      uint      `json:"user_id" example:"1"`              // 作者ID
	CategoryID  uint      `json:"category_id" example:"1"`          // 分类ID
	Status      int       `json:"status" example:"1"`               // 状态：1公开 2草稿
	ViewCount   int       `json:"view_count" example:"100"`         // 浏览次数
	Tags        []Tag     `json:"tags"`                             // 标签列表
	CreatedAt   string    `json:"created_at" example:"2024-01-01 00:00:00"` // 创建时间
	UpdatedAt   string    `json:"updated_at" example:"2024-01-01 00:00:00"` // 更新时间
}

// Category 分类信息
// @Description Category model
// @Description 分类信息
// @Success 200 {object} docs.Category "成功"
// @Failure 404 {object} response.Response "未找到"
// @Failure 500 {object} response.Response "服务器错误"
// @Router /categories/{id} [get]
type Category struct {
	ID          uint   `json:"id" example:"1"`
	Name        string `json:"name" example:"Technology"`
	Description string `json:"description" example:"Tech related posts"`
	CreatedAt   string `json:"created_at" example:"2024-01-01 00:00:00"`
	UpdatedAt   string `json:"updated_at" example:"2024-01-01 00:00:00"`
}

// Tag 标签信息
type Tag struct {
	ID        uint   `json:"id" example:"1"`
	Name      string `json:"name" example:"golang"`
	CreatedAt string `json:"created_at" example:"2024-01-01 00:00:00"`
	UpdatedAt string `json:"updated_at" example:"2024-01-01 00:00:00"`
}

var docTemplate = `{
    "swagger": "2.0",
    "info": {
        "title": "Personal Blog API",
        "version": "1.0",
        "description": "API documentation for the Personal Blog system."
    },
    "paths": {}
}`

func SwaggerInfo() {
	swag.Register(swag.Name, &swag.Spec{
		InfoInstanceName: "swagger",
		SwaggerTemplate: docTemplate,
	})
}
