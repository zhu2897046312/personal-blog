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

// CategoryHandler 分类处理器
type CategoryHandler struct {
	categoryService service.CategoryService
}

// NewCategoryHandler 创建分类处理器实例
func NewCategoryHandler(categoryService service.CategoryService) *CategoryHandler {
	return &CategoryHandler{
		categoryService: categoryService,
	}
}

// Create godoc
// @Summary 创建分类
// @Description 创建新分类（管理员）
// @Tags category
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param data body request.CreateCategoryRequest true "分类信息"
// @Success 200 {object} response.Response{data=docs.Category} "创建成功"
// @Failure 400 {object} response.Response "请求参数错误"
// @Failure 401 {object} response.Response "未登录"
// @Failure 403 {object} response.Response "权限不足"
// @Failure 500 {object} response.Response "服务器内部错误"
// @Router /categories [post]
func (h *CategoryHandler) Create(c *gin.Context) {
	var req request.CreateCategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.NewResponse(http.StatusBadRequest, err.Error(), nil))
		return
	}

	category := &models.Category{
		Name:        req.Name,
		Description: req.Description,
	}

	if err := h.categoryService.CreateCategory(c.Request.Context(), category); err != nil {
		c.JSON(http.StatusInternalServerError, response.NewResponse(http.StatusInternalServerError, err.Error(), nil))
		return
	}

	c.JSON(http.StatusOK, response.NewResponse(http.StatusOK, "创建成功", category))
}

// Update godoc
// @Summary 更新分类
// @Description 更新分类信息（管理员）
// @Tags category
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id path int true "分类ID"
// @Param data body request.UpdateCategoryRequest true "分类信息"
// @Success 200 {object} response.Response{data=docs.Category} "更新成功"
// @Failure 400 {object} response.Response "请求参数错误"
// @Failure 401 {object} response.Response "未登录"
// @Failure 403 {object} response.Response "权限不足"
// @Failure 404 {object} response.Response "分类不存在"
// @Failure 500 {object} response.Response "服务器内部错误"
// @Router /categories/{id} [put]
func (h *CategoryHandler) Update(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.NewResponse(http.StatusBadRequest, "无效的分类ID", nil))
		return
	}

	var req request.UpdateCategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.NewResponse(http.StatusBadRequest, err.Error(), nil))
		return
	}

	category, err := h.categoryService.GetCategoryByID(c.Request.Context(), uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, response.NewResponse(http.StatusNotFound, "分类不存在", nil))
		return
	}

	category.Name = req.Name
	category.Description = req.Description

	if err := h.categoryService.UpdateCategory(c.Request.Context(), category); err != nil {
		c.JSON(http.StatusInternalServerError, response.NewResponse(http.StatusInternalServerError, err.Error(), nil))
		return
	}

	c.JSON(http.StatusOK, response.NewResponse(http.StatusOK, "更新成功", category))
}

// Delete godoc
// @Summary 删除分类
// @Description 删除分类（管理员）
// @Tags category
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id path int true "分类ID"
// @Success 200 {object} response.Response "删除成功"
// @Failure 401 {object} response.Response "未登录"
// @Failure 403 {object} response.Response "权限不足"
// @Failure 404 {object} response.Response "分类不存在"
// @Failure 500 {object} response.Response "服务器内部错误"
// @Router /categories/{id} [delete]
func (h *CategoryHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.NewResponse(http.StatusBadRequest, "无效的分类ID", nil))
		return
	}

	if err := h.categoryService.DeleteCategory(c.Request.Context(), uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, response.NewResponse(http.StatusInternalServerError, err.Error(), nil))
		return
	}

	c.JSON(http.StatusOK, response.NewResponse(http.StatusOK, "删除成功", nil))
}

// Get godoc
// @Summary 获取分类详情
// @Description 获取分类详细信息
// @Tags category
// @Accept json
// @Produce json
// @Param id path int true "分类ID"
// @Success 200 {object} response.Response{data=docs.Category} "获取成功"
// @Failure 404 {object} response.Response "分类不存在"
// @Failure 500 {object} response.Response "服务器内部错误"
// @Router /categories/{id} [get]
func (h *CategoryHandler) Get(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.NewResponse(http.StatusBadRequest, "无效的分类ID", nil))
		return
	}

	category, err := h.categoryService.GetCategoryByID(c.Request.Context(), uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, response.NewResponse(http.StatusNotFound, "分类不存在", nil))
		return
	}

	c.JSON(http.StatusOK, response.NewResponse(http.StatusOK, "获取成功", category))
}

// List godoc
// @Summary 获取分类列表
// @Description 获取所有分类列表
// @Tags category
// @Accept json
// @Produce json
// @Success 200 {object} response.Response{data=[]docs.Category} "获取成功"
// @Failure 500 {object} response.Response "服务器内部错误"
// @Router /categories [get]
func (h *CategoryHandler) List(c *gin.Context) {
	categories, err := h.categoryService.ListCategories(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.NewResponse(http.StatusInternalServerError, err.Error(), nil))
		return
	}

	c.JSON(http.StatusOK, response.NewResponse(http.StatusOK, "获取成功", categories))
}
