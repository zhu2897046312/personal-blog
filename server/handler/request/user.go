package request

// RegisterRequest 用户注册请求
type RegisterRequest struct {
	Username string `json:"username" binding:"required,min=3,max=32"`
	Password string `json:"password" binding:"required,min=6,max=32"`
	Email    string `json:"email" binding:"required,email"`
	Nickname string `json:"nickname" binding:"required,min=2,max=32"`
	Role     string `json:"role" binding:"required,oneof=admin user"`
}

// LoginRequest 用户登录请求
type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// UpdateProfileRequest 更新用户信息请求
type UpdateProfileRequest struct {
	Nickname string `json:"nickname" binding:"required,min=2,max=32"`
	Email    string `json:"email" binding:"required,email"`
	Avatar   string `json:"avatar" binding:"omitempty,max=200"`
}

// ChangePasswordRequest 修改密码请求
type ChangePasswordRequest struct {
	OldPassword string `json:"old_password" binding:"required"`
	NewPassword string `json:"new_password" binding:"required,min=6,max=32"`
}

// ListUsersRequest 用户列表请求
type ListUsersRequest struct {
	SearchRequest
}

// UpdateUserStatusRequest 更新用户状态请求
type UpdateUserStatusRequest struct {
	StatusRequest
}
