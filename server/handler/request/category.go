package request

// CreateCategoryRequest 创建分类请求
type CreateCategoryRequest struct {
	Name        string `json:"name" binding:"required,min=1,max=50"`
	Description string `json:"description" binding:"omitempty,max=200"`
}

// UpdateCategoryRequest 更新分类请求
type UpdateCategoryRequest struct {
	Name        string `json:"name" binding:"required,min=1,max=50"`
	Description string `json:"description" binding:"omitempty,max=200"`
}
