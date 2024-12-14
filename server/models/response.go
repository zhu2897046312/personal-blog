package models

// Response is the standard API response structure
type Response struct {
	Code int         `json:"code"`    // HTTP status code
	Msg  string      `json:"message"` // Response message
	Data interface{} `json:"data"`    // Response data
}

// LoginResponse contains the login response data
type LoginResponse struct {
	Token string `json:"token"`  // JWT token
	User  *User  `json:"user"`   // User information
}

// LoginRequest contains the login request data
type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}
