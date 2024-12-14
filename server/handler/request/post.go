package request

// CreatePostRequest 创建文章请求
type CreatePostRequest struct {
	Title      string   `json:"title" binding:"required,min=1,max=100"`
	Content    string   `json:"content" binding:"required,min=1"`
	CategoryID uint     `json:"category_id" binding:"required"`
	Tags       []string `json:"tags" binding:"omitempty,dive,min=1"`
	Status     int      `json:"status" binding:"required,oneof=1 2"` // 1:公开 2:草稿
}

// UpdatePostRequest 更新文章请求
type UpdatePostRequest struct {
	Title      string   `json:"title" binding:"required,min=1,max=100"`
	Content    string   `json:"content" binding:"required,min=1"`
	CategoryID uint     `json:"category_id" binding:"required"`
	Tags       []string `json:"tags" binding:"omitempty,dive,min=1"`
	Status     int      `json:"status" binding:"required,oneof=1 2"` // 1:公开 2:草稿
}

// ListPostsRequest 文章列表请求
type ListPostsRequest struct {
	Page       int    `form:"page" binding:"required,min=1"`
	PageSize   int    `form:"page_size" binding:"required,min=1,max=100"`
	CategoryID uint   `form:"category_id"`
	TagID      uint   `form:"tag_id"`
	UserID     uint   `form:"user_id"`
	Keyword    string `form:"keyword"`
	Status     int    `form:"status" binding:"omitempty,oneof=1 2"` // 1:公开 2:草稿
}
