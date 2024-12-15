package request

// CreateCommentRequest 创建评论请求
type CreateCommentRequest struct {
	PostID  uint   `json:"post_id" binding:"required,min=1"`
	Content string `json:"content" binding:"required,min=1,max=500"`
	ParentID uint  `json:"parent_id" binding:"omitempty,min=1"`
}

// ListCommentsRequest 评论列表请求
type ListCommentsRequest struct {
	PostID uint `form:"post_id" binding:"required,min=1"`
	SearchRequest
}

// UpdateCommentStatusRequest 更新评论状态请求
type UpdateCommentStatusRequest struct {
	StatusRequest
}
