package handler

import (
	"net/http"
	"strconv"

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

// Register 用户注册
// @Summary 用户注册
// @Description 用户注册接口
// @Tags user
// @Accept json
// @Produce json
// @Param data body request.RegisterRequest true "注册信息"
// @Success 200 {object} response.Response{data=models.User} "注册成功"
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
		Role:     req.Role,
	}

	if err := h.userService.Register(c, user); err != nil {
		c.JSON(http.StatusInternalServerError, response.NewResponse(http.StatusInternalServerError, err.Error(), nil))
		return
	}

	c.JSON(http.StatusOK, response.NewResponse(http.StatusOK, "注册成功", user))
}

// Login 用户登录
// @Summary 用户登录
// @Description 用户登录接口
// @Tags user
// @Accept json
// @Produce json
// @Param data body request.LoginRequest true "登录信息"
// @Success 200 {object} response.Response{data=models.User} "登录成功"
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

	user, token, err := h.userService.Login(c, req.Username, req.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, response.NewResponse(http.StatusUnauthorized, err.Error(), nil))
		return
	}

	c.JSON(http.StatusOK, response.NewResponse(http.StatusOK, "登录成功", gin.H{
		"token": token,
		"user":  user,
	}))
}

// GetProfile 获取用户信息
// @Summary 获取用户信息
// @Description 获取当前登录用户信息
// @Tags user
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {object} response.Response{data=models.User} "获取成功"
// @Failure 401 {object} response.Response "未登录"
// @Failure 500 {object} response.Response "服务器内部错误"
// @Router /users/profile [get]
func (h *UserHandler) GetProfile(c *gin.Context) {
	userID, _ := c.Get("userID")
	user, err := h.userService.GetUserByID(c, userID.(uint))
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.NewResponse(http.StatusInternalServerError, err.Error(), nil))
		return
	}

	c.JSON(http.StatusOK, response.NewResponse(http.StatusOK, "获取成功", user))
}

// UpdateProfile 更新用户信息
// @Summary 更新用户信息
// @Description 更新当前登录用户信息
// @Tags user
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param data body request.UpdateProfileRequest true "用户信息"
// @Success 200 {object} response.Response{data=models.User} "更新成功"
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

	userID, _ := c.Get("userID")
	user := &models.User{
		ID:       userID.(uint),
		Nickname: req.Nickname,
		Email:    req.Email,
		Avatar:   req.Avatar,
	}

	if err := h.userService.UpdateUser(c, user); err != nil {
		c.JSON(http.StatusInternalServerError, response.NewResponse(http.StatusInternalServerError, err.Error(), nil))
		return
	}

	c.JSON(http.StatusOK, response.NewResponse(http.StatusOK, "更新成功", user))
}

// ChangePassword 修改密码
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

	userID, _ := c.Get("userID")
	if err := h.userService.ChangePassword(c, userID.(uint), req.OldPassword, req.NewPassword); err != nil {
		c.JSON(http.StatusInternalServerError, response.NewResponse(http.StatusInternalServerError, err.Error(), nil))
		return
	}

	c.JSON(http.StatusOK, response.NewResponse(http.StatusOK, "修改成功", nil))
}

// List 获取用户列表
// @Summary 获取用户列表
// @Description 获取用户列表（管理员）
// @Tags user
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param page query int false "页码" default(1)
// @Param page_size query int false "每页数量" default(10)
// @Param keyword query string false "搜索关键词"
// @Success 200 {object} response.Response{data=response.PaginationData{items=[]models.User}} "获取成功"
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

	users, total, err := h.userService.ListUsers(c, req.Page, req.PageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.NewResponse(http.StatusInternalServerError, err.Error(), nil))
		return
	}

	c.JSON(http.StatusOK, response.NewResponse(http.StatusOK, "获取成功", response.NewPaginationResponse(users, total, req.Page, req.PageSize)))
}

// UpdateStatus 更新用户状态
// @Summary 更新用户状态
// @Description 更新用户状态（管理员）
// @Tags user
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id path int true "用户ID"
// @Param data body request.UpdateUserStatusRequest true "状态信息"
// @Success 200 {object} response.Response "更新成功"
// @Failure 400 {object} response.Response "请求参数错误"
// @Failure 401 {object} response.Response "未登录"
// @Failure 403 {object} response.Response "权限不足"
// @Failure 500 {object} response.Response "服务器内部错误"
// @Router /users/{id}/status [put]
func (h *UserHandler) UpdateStatus(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.NewResponse(http.StatusBadRequest, "参数错误", nil))
		return
	}

	var req request.UpdateUserStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.NewResponse(http.StatusBadRequest, err.Error(), nil))
		return
	}

	user := &models.User{
		ID:     uint(id),
		Status: req.Status,
	}

	if err := h.userService.UpdateUser(c, user); err != nil {
		c.JSON(http.StatusInternalServerError, response.NewResponse(http.StatusInternalServerError, err.Error(), nil))
		return
	}

	c.JSON(http.StatusOK, response.NewResponse(http.StatusOK, "更新成功", nil))
}
