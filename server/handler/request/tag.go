package request

// CreateTagRequest 创建标签请求
type CreateTagRequest struct {
	Name        string `json:"name" binding:"required,min=1,max=50"`
}

// UpdateTagRequest 更新标签请求
type UpdateTagRequest struct {
	Name        string `json:"name" binding:"required,min=1,max=50"`
}

// CreateTagsRequest 批量创建标签请求
type CreateTagsRequest struct {
	Names []string `json:"names" binding:"required,dive,min=1,max=50"`
}
