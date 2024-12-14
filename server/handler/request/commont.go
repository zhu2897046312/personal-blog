package request

type CreateCommentRequest struct {
	PostID   uint   `json:"post_id" binding:"required"`
	Content  string `json:"content" binding:"required,min=1,max=500"`
	ParentID *uint  `json:"parent_id,omitempty"`
}

type ListCommentsRequest struct {
	Page     int `form:"page" binding:"required,min=1"`
	PageSize int `form:"page_size" binding:"required,min=1,max=50"`
}