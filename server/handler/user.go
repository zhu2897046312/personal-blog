package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/personal-blog/handler/request"
	"github.com/personal-blog/handler/response"
	"github.com/personal-blog/models"
	"github.com/personal-blog/service"
)

// UserHandler 用户处理器
type UserHandler struct {
	userService service.UserService
}

// NewUserHandler 创建用户处理器实例
func NewUserHandler(userService service.UserService) *UserHandler {
	return &UserHandler{
		userService: userService,
	}
}

// Register godoc
// @Summary 用户注册
// @Description 用户注册接口
// @Tags user
// @Accept json
// @Produce json
// @Param data body request.RegisterRequest true "注册信息"
// @Success 200 {object} response.Response{data=docs.User} "注册成功"
// @Failure 400 {object} response.Response "请求参数错误"
// @Failure 500 {object} response.Response "服务器内部错误"
// @Router /users/register [post]
func (h *UserHandler) Register(c *gin.Context) {
	var req request.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.NewResponse(http.StatusBadRequest, err.Error(), nil))
		return
	}

	user := &models.User{
		Username: req.Username,
		Password: req.Password,
		Email:    req.Email,
		Nickname: req.Nickname,
	}

	if err := h.userService.Register(c.Request.Context(), user); err != nil {
		c.JSON(http.StatusBadRequest, response.NewResponse(http.StatusBadRequest, err.Error(), nil))
		return
	}

	c.JSON(http.StatusOK, response.NewResponse(http.StatusOK, "注册成功", user))
}

// Login godoc
// @Summary 用户登录
// @Description 用户登录接口
// @Tags user
// @Accept json
// @Produce json
// @Param data body request.LoginRequest true "登录信息"
// @Success 200 {object} response.Response{data=docs.LoginResponse} "登录成功"
// @Failure 400 {object} response.Response "请求参数错误"
// @Failure 401 {object} response.Response "用户名或密码错误"
// @Failure 500 {object} response.Response "服务器内部错误"
// @Router /users/login [post]
func (h *UserHandler) Login(c *gin.Context) {
	var req request.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.NewResponse(http.StatusBadRequest, err.Error(), nil))
		return
	}

	user, token, err := h.userService.Login(c.Request.Context(), req.Username, req.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, response.NewResponse(http.StatusUnauthorized, err.Error(), nil))
		return
	}

	c.JSON(http.StatusOK, response.NewResponse(http.StatusOK, "登录成功", gin.H{
		"user":  user,
		"token": token,
	}))
}

// GetProfile godoc
// @Summary 获取用户信息
// @Description 获取当前登录用户信息
// @Tags user
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {object} response.Response{data=docs.User} "获取成功"
// @Failure 401 {object} response.Response "未登录"
// @Failure 500 {object} response.Response "服务器内部错误"
// @Router /users/profile [get]
func (h *UserHandler) GetProfile(c *gin.Context) {
	userID := c.GetUint("user_id")
	user, err := h.userService.GetUserByID(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.NewResponse(http.StatusInternalServerError, err.Error(), nil))
		return
	}

	c.JSON(http.StatusOK, response.NewResponse(http.StatusOK, "获取成功", user))
}

// UpdateProfile godoc
// @Summary 更新用户信息
// @Description 更新当前登录用户信息
// @Tags user
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param data body request.UpdateProfileRequest true "用户信息"
// @Success 200 {object} response.Response{data=docs.User} "更新成功"
// @Failure 400 {object} response.Response "请求参数错误"
// @Failure 401 {object} response.Response "未登录"
// @Failure 500 {object} response.Response "服务器内部错误"
// @Router /users/profile [put]
func (h *UserHandler) UpdateProfile(c *gin.Context) {
	var req request.UpdateProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.NewResponse(http.StatusBadRequest, err.Error(), nil))
		return
	}

	userID := c.GetUint("user_id")
	user, err := h.userService.GetUserByID(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.NewResponse(http.StatusInternalServerError, err.Error(), nil))
		return
	}

	user.Nickname = req.Nickname
	user.Email = req.Email
	if req.Avatar != "" {
		user.Avatar = req.Avatar
	}

	if err := h.userService.UpdateUser(c.Request.Context(), user); err != nil {
		c.JSON(http.StatusInternalServerError, response.NewResponse(http.StatusInternalServerError, err.Error(), nil))
		return
	}

	c.JSON(http.StatusOK, response.NewResponse(http.StatusOK, "更新成功", user))
}

// ChangePassword godoc
// @Summary 修改密码
// @Description 修改当前登录用户密码
// @Tags user
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param data body request.ChangePasswordRequest true "密码信息"
// @Success 200 {object} response.Response "修改成功"
// @Failure 400 {object} response.Response "请求参数错误"
// @Failure 401 {object} response.Response "未登录"
// @Failure 500 {object} response.Response "服务器内部错误"
// @Router /users/password [put]
func (h *UserHandler) ChangePassword(c *gin.Context) {
	var req request.ChangePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.NewResponse(http.StatusBadRequest, err.Error(), nil))
		return
	}

	userID := c.GetUint("user_id")
	if err := h.userService.ChangePassword(c.Request.Context(), userID, req.OldPassword, req.NewPassword); err != nil {
		c.JSON(http.StatusBadRequest, response.NewResponse(http.StatusBadRequest, err.Error(), nil))
		return
	}

	c.JSON(http.StatusOK, response.NewResponse(http.StatusOK, "修改成功", nil))
}

// ListUsers godoc
// @Summary 获取用户列表
// @Description 获取用户列表（管理员）
// @Tags user
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param page query int false "页码" default(1)
// @Param page_size query int false "每页数量" default(10)
// @Success 200 {object} response.Response{data=docs.PaginationData{items=[]docs.User}} "获取成功"
// @Failure 401 {object} response.Response "未登录"
// @Failure 403 {object} response.Response "权限不足"
// @Failure 500 {object} response.Response "服务器内部错误"
// @Router /users [get]
func (h *UserHandler) ListUsers(c *gin.Context) {
	var req request.ListUsersRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.NewResponse(http.StatusBadRequest, err.Error(), nil))
		return
	}

	users, total, err := h.userService.ListUsers(c.Request.Context(), req.Page, req.PageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.NewResponse(http.StatusInternalServerError, err.Error(), nil))
		return
	}

	c.JSON(http.StatusOK, response.NewResponse(http.StatusOK, "获取成功", response.NewPaginationResponse(users, total, req.Page, req.PageSize)))
}
