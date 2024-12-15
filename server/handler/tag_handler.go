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

// TagHandler 标签处理器
type TagHandler struct {
	tagService service.TagService
}

// NewTagHandler 创建标签处理器实例
func NewTagHandler(tagService service.TagService) *TagHandler {
	return &TagHandler{
		tagService: tagService,
	}
}

// Create godoc
// @Summary 创建标签
// @Description 创建新标签（管理员）
// @Tags tag
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param data body request.CreateTagRequest true "标签信息"
// @Success 200 {object} response.Response{data=models.Tag} "创建成功"
// @Failure 400 {object} response.Response "请求参数错误"
// @Failure 401 {object} response.Response "未登录"
// @Failure 403 {object} response.Response "权限不足"
// @Failure 500 {object} response.Response "服务器内部错误"
// @Router /tags [post]
func (h *TagHandler) Create(c *gin.Context) {
	var req request.CreateTagRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.NewResponse(http.StatusBadRequest, err.Error(), nil))
		return
	}

	tag := &models.Tag{
		Name:        req.Name,
	}

	if err := h.tagService.CreateTag(c.Request.Context(), tag); err != nil {
		c.JSON(http.StatusInternalServerError, response.NewResponse(http.StatusInternalServerError, err.Error(), nil))
		return
	}

	c.JSON(http.StatusOK, response.NewResponse(http.StatusOK, "创建成功", tag))
}

// CreateBatch godoc
// @Summary 批量创建标签
// @Description 批量创建标签（管理员）
// @Tags tag
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param data body request.CreateTagsRequest true "标签列表"
// @Success 200 {object} response.Response{data=[]models.Tag} "创建成功"
// @Failure 400 {object} response.Response "请求参数错误"
// @Failure 401 {object} response.Response "未登录"
// @Failure 403 {object} response.Response "权限不足"
// @Failure 500 {object} response.Response "服务器内部错误"
// @Router /tags/batch [post]
func (h *TagHandler) CreateBatch(c *gin.Context) {
	var req request.CreateTagsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.NewResponse(http.StatusBadRequest, err.Error(), nil))
		return
	}

	tags, err := h.tagService.CreateTagsIfNotExist(c.Request.Context(), req.Names)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.NewResponse(http.StatusInternalServerError, err.Error(), nil))
		return
	}

	c.JSON(http.StatusOK, response.NewResponse(http.StatusOK, "创建成功", tags))
}

// Update godoc
// @Summary 更新标签
// @Description 更新标签信息（管理员）
// @Tags tag
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id path int true "标签ID"
// @Param data body request.UpdateTagRequest true "标签信息"
// @Success 200 {object} response.Response{data=models.Tag} "更新成功"
// @Failure 400 {object} response.Response "请求参数错误"
// @Failure 401 {object} response.Response "未登录"
// @Failure 403 {object} response.Response "权限不足"
// @Failure 404 {object} response.Response "标签不存在"
// @Failure 500 {object} response.Response "服务器内部错误"
// @Router /tags/{id} [put]
func (h *TagHandler) Update(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.NewResponse(http.StatusBadRequest, "无效的标签ID", nil))
		return
	}

	var req request.UpdateTagRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.NewResponse(http.StatusBadRequest, err.Error(), nil))
		return
	}

	tag, err := h.tagService.GetTagByID(c.Request.Context(), uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, response.NewResponse(http.StatusNotFound, "标签不存在", nil))
		return
	}

	tag.Name = req.Name

	if err := h.tagService.UpdateTag(c.Request.Context(), tag); err != nil {
		c.JSON(http.StatusInternalServerError, response.NewResponse(http.StatusInternalServerError, err.Error(), nil))
		return
	}

	c.JSON(http.StatusOK, response.NewResponse(http.StatusOK, "更新成功", tag))
}

// Delete godoc
// @Summary 删除标签
// @Description 删除标签（管理员）
// @Tags tag
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id path int true "标签ID"
// @Success 200 {object} response.Response "删除成功"
// @Failure 401 {object} response.Response "未登录"
// @Failure 403 {object} response.Response "权限不足"
// @Failure 404 {object} response.Response "标签不存在"
// @Failure 500 {object} response.Response "服务器内部错误"
// @Router /tags/{id} [delete]
func (h *TagHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.NewResponse(http.StatusBadRequest, "无效的标签ID", nil))
		return
	}

	if err := h.tagService.DeleteTag(c.Request.Context(), uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, response.NewResponse(http.StatusInternalServerError, err.Error(), nil))
		return
	}

	c.JSON(http.StatusOK, response.NewResponse(http.StatusOK, "删除成功", nil))
}

// Get godoc
// @Summary 获取标签详情
// @Description 获取标签详细信息
// @Tags tag
// @Accept json
// @Produce json
// @Param id path int true "标签ID"
// @Success 200 {object} response.Response{data=models.Tag} "获取成功"
// @Failure 404 {object} response.Response "标签不存在"
// @Failure 500 {object} response.Response "服务器内部错误"
// @Router /tags/{id} [get]
func (h *TagHandler) Get(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.NewResponse(http.StatusBadRequest, "无效的标签ID", nil))
		return
	}

	tag, err := h.tagService.GetTagByID(c.Request.Context(), uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, response.NewResponse(http.StatusNotFound, "标签不存在", nil))
		return
	}

	c.JSON(http.StatusOK, response.NewResponse(http.StatusOK, "获取成功", tag))
}

// List godoc
// @Summary 获取标签列表
// @Description 获取所有标签列表
// @Tags tag
// @Accept json
// @Produce json
// @Success 200 {object} response.Response{data=[]models.Tag} "获取成功"
// @Failure 500 {object} response.Response "服务器内部错误"
// @Router /tags [get]
func (h *TagHandler) List(c *gin.Context) {
	tags, err := h.tagService.ListTags(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.NewResponse(http.StatusInternalServerError, err.Error(), nil))
		return
	}

	c.JSON(http.StatusOK, response.NewResponse(http.StatusOK, "获取成功", tags))
}

// GetPostTags godoc
// @Summary 获取文章标签
// @Description 获取指定文章的标签列表
// @Tags tag
// @Accept json
// @Produce json
// @Param post_id path int true "文章ID"
// @Success 200 {object} response.Response{data=[]models.Tag} "获取成功"
// @Failure 404 {object} response.Response "文章不存在"
// @Failure 500 {object} response.Response "服务器内部错误"
// @Router /posts/{post_id}/tags [get]
func (h *TagHandler) GetPostTags(c *gin.Context) {
	postID, err := strconv.ParseUint(c.Param("post_id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.NewResponse(http.StatusBadRequest, "无效的文章ID", nil))
		return
	}

	tags, err := h.tagService.GetPostTags(c.Request.Context(), uint(postID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.NewResponse(http.StatusInternalServerError, err.Error(), nil))
		return
	}

	c.JSON(http.StatusOK, response.NewResponse(http.StatusOK, "获取成功", tags))
}
