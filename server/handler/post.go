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

// PostHandler 文章处理器
type PostHandler struct {
	postService service.PostService
}

// NewPostHandler 创建文章处理器实例
func NewPostHandler(postService service.PostService) *PostHandler {
	return &PostHandler{
		postService: postService,
	}
}

// Create godoc
// @Summary 创建文章
// @Description 创建新文章
// @Tags post
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param data body request.CreatePostRequest true "文章信息"
// @Success 200 {object} response.Response{data=docs.Post} "创建成功"
// @Failure 400 {object} response.Response "请求参数错误"
// @Failure 401 {object} response.Response "未登录"
// @Failure 500 {object} response.Response "服务器内部错误"
// @Router /posts [post]
func (h *PostHandler) Create(c *gin.Context) {
	var req request.CreatePostRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.NewResponse(http.StatusBadRequest, err.Error(), nil))
		return
	}

	userID := c.GetUint("user_id")
	post := &models.Post{
		Title:      req.Title,
		Content:    req.Content,
		UserID:     userID,
		CategoryID: req.CategoryID,
		Status:     req.Status,
	}

	if err := h.postService.CreatePost(c.Request.Context(), post, req.Tags); err != nil {
		c.JSON(http.StatusInternalServerError, response.NewResponse(http.StatusInternalServerError, err.Error(), nil))
		return
	}

	c.JSON(http.StatusOK, response.NewResponse(http.StatusOK, "创建成功", post))
}

// Update godoc
// @Summary 更新文章
// @Description 更新文章信息
// @Tags post
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id path int true "文章ID"
// @Param data body request.UpdatePostRequest true "文章信息"
// @Success 200 {object} response.Response{data=docs.Post} "更新成功"
// @Failure 400 {object} response.Response "请求参数错误"
// @Failure 401 {object} response.Response "未登录"
// @Failure 403 {object} response.Response "权限不足"
// @Failure 404 {object} response.Response "文章不存在"
// @Failure 500 {object} response.Response "服务器内部错误"
// @Router /posts/{id} [put]
func (h *PostHandler) Update(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.NewResponse(http.StatusBadRequest, "无效的文章ID", nil))
		return
	}

	var req request.UpdatePostRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.NewResponse(http.StatusBadRequest, err.Error(), nil))
		return
	}

	post, err := h.postService.GetPostByID(c.Request.Context(), uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, response.NewResponse(http.StatusNotFound, "文章不存在", nil))
		return
	}

	userID := c.GetUint("user_id")
	if post.UserID != userID {
		c.JSON(http.StatusForbidden, response.NewResponse(http.StatusForbidden, "无权操作", nil))
		return
	}

	post.Title = req.Title
	post.Content = req.Content
	post.CategoryID = req.CategoryID
	post.Status = req.Status

	if err := h.postService.UpdatePost(c.Request.Context(), post, req.Tags); err != nil {
		c.JSON(http.StatusInternalServerError, response.NewResponse(http.StatusInternalServerError, err.Error(), nil))
		return
	}

	c.JSON(http.StatusOK, response.NewResponse(http.StatusOK, "更新成功", post))
}

// Delete godoc
// @Summary 删除文章
// @Description 删除文章
// @Tags post
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id path int true "文章ID"
// @Success 200 {object} response.Response "删除成功"
// @Failure 401 {object} response.Response "未登录"
// @Failure 403 {object} response.Response "权限不足"
// @Failure 404 {object} response.Response "文章不存在"
// @Failure 500 {object} response.Response "服务器内部错误"
// @Router /posts/{id} [delete]
func (h *PostHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.NewResponse(http.StatusBadRequest, "无效的文章ID", nil))
		return
	}

	post, err := h.postService.GetPostByID(c.Request.Context(), uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, response.NewResponse(http.StatusNotFound, "文章不存在", nil))
		return
	}

	userID := c.GetUint("user_id")
	if post.UserID != userID {
		c.JSON(http.StatusForbidden, response.NewResponse(http.StatusForbidden, "无权操作", nil))
		return
	}

	if err := h.postService.DeletePost(c.Request.Context(), uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, response.NewResponse(http.StatusInternalServerError, err.Error(), nil))
		return
	}

	c.JSON(http.StatusOK, response.NewResponse(http.StatusOK, "删除成功", nil))
}

// Get godoc
// @Summary 获取文章详情
// @Description 获取文章详细信息
// @Tags post
// @Accept json
// @Produce json
// @Param id path int true "文章ID"
// @Success 200 {object} response.Response{data=docs.Post} "获取成功"
// @Failure 404 {object} response.Response "文章不存在"
// @Failure 500 {object} response.Response "服务器内部错误"
// @Router /posts/{id} [get]
func (h *PostHandler) Get(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.NewResponse(http.StatusBadRequest, "无效的文章ID", nil))
		return
	}

	post, err := h.postService.GetPostByID(c.Request.Context(), uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, response.NewResponse(http.StatusNotFound, "文章不存在", nil))
		return
	}

	// 增加浏览次数
	go h.postService.IncrementViewCount(c.Request.Context(), uint(id))

	c.JSON(http.StatusOK, response.NewResponse(http.StatusOK, "获取成功", post))
}

// List godoc
// @Summary 获取文章列表
// @Description 获取文章列表，支持分页和筛选
// @Tags post
// @Accept json
// @Produce json
// @Param page query int false "页码" default(1)
// @Param page_size query int false "每页数量" default(10)
// @Param category_id query int false "分类ID"
// @Param tag_id query int false "标签ID"
// @Param user_id query int false "作者ID"
// @Param status query int false "状态" Enums(1,2)
// @Success 200 {object} response.Response{data=docs.PaginationData{items=[]docs.Post}} "获取成功"
// @Failure 400 {object} response.Response "请求参数错误"
// @Failure 500 {object} response.Response "服务器内部错误"
// @Router /posts [get]
func (h *PostHandler) List(c *gin.Context) {
	var req request.ListPostsRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.NewResponse(http.StatusBadRequest, err.Error(), nil))
		return
	}

	var posts []models.Post
	var total int64
	var err error

	// 根据不同条件调用不同的列表接口
	switch {
	case req.CategoryID > 0:
		posts, total, err = h.postService.ListPostsByCategory(c.Request.Context(), req.CategoryID, req.Page, req.PageSize)
	case req.TagID > 0:
		posts, total, err = h.postService.ListPostsByTag(c.Request.Context(), req.TagID, req.Page, req.PageSize)
	case req.UserID > 0:
		posts, total, err = h.postService.ListPostsByUser(c.Request.Context(), req.UserID, req.Page, req.PageSize)
	default:
		conditions := make(map[string]interface{})
		if req.Status > 0 {
			conditions["status"] = req.Status
		}
		if req.Keyword != "" {
			conditions["keyword"] = req.Keyword
		}
		posts, total, err = h.postService.ListPosts(c.Request.Context(), req.Page, req.PageSize, conditions)
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, response.NewResponse(http.StatusInternalServerError, err.Error(), nil))
		return
	}

	c.JSON(http.StatusOK, response.NewResponse(http.StatusOK, "获取成功", response.NewPaginationResponse(posts, total, req.Page, req.PageSize)))
}
