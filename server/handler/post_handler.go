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

// Create 创建文章
func (h *PostHandler) Create(c *gin.Context) {
	var req request.CreatePostRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.NewResponse(http.StatusBadRequest, err.Error(), nil))
		return
	}

	userID, _ := c.Get("userID")
	post := &models.Post{
		Title:      req.Title,
		Content:    req.Content,
		CategoryID: req.CategoryID,
		UserID:     userID.(uint),
	}

	if err := h.postService.CreatePost(c, post, req.Tags); err != nil {
		c.JSON(http.StatusInternalServerError, response.NewResponse(http.StatusInternalServerError, err.Error(), nil))
		return
	}

	c.JSON(http.StatusOK, response.NewResponse(http.StatusOK, "创建成功", post))
}

// Update 更新文章
func (h *PostHandler) Update(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.NewResponse(http.StatusBadRequest, "参数错误", nil))
		return
	}

	var req request.UpdatePostRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.NewResponse(http.StatusBadRequest, err.Error(), nil))
		return
	}

	post := &models.Post{
		ID:         uint(id),
		Title:      req.Title,
		Content:    req.Content,
		CategoryID: req.CategoryID,
	}

	if err := h.postService.UpdatePost(c, post, req.Tags); err != nil {
		c.JSON(http.StatusInternalServerError, response.NewResponse(http.StatusInternalServerError, err.Error(), nil))
		return
	}

	c.JSON(http.StatusOK, response.NewResponse(http.StatusOK, "更新成功", post))
}

// Delete 删除文章
func (h *PostHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.NewResponse(http.StatusBadRequest, "参数错误", nil))
		return
	}

	if err := h.postService.DeletePost(c, uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, response.NewResponse(http.StatusInternalServerError, err.Error(), nil))
		return
	}

	c.JSON(http.StatusOK, response.NewResponse(http.StatusOK, "删除成功", nil))
}

// Get 获取文章详情
func (h *PostHandler) Get(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.NewResponse(http.StatusBadRequest, "参数错误", nil))
		return
	}

	post, err := h.postService.GetPostByID(c, uint(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.NewResponse(http.StatusInternalServerError, err.Error(), nil))
		return
	}

	c.JSON(http.StatusOK, response.NewResponse(http.StatusOK, "获取成功", post))
}

// List 获取文章列表
// List 获取文章列表
func (h *PostHandler) List(c *gin.Context) {
	var req request.ListPostsRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.NewResponse(http.StatusBadRequest, err.Error(), nil))
		return
	}

	conditions := make(map[string]interface{})
	if req.Keyword != "" {
		conditions["keyword"] = req.Keyword
	}
	if req.CategoryID > 0 {
		conditions["category_id"] = req.CategoryID
	}
	if req.Tag != "" {
		conditions["tag"] = req.Tag
	}
	if req.Status > 0 {
		conditions["status"] = req.Status
	}

	posts, total, err := h.postService.ListPosts(c, req.Page, req.PageSize, conditions)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.NewResponse(http.StatusInternalServerError, err.Error(), nil))
		return
	}

	c.JSON(http.StatusOK, response.NewResponse(http.StatusOK, "获取成功", response.NewPaginationResponse(posts, total, req.Page, req.PageSize)))
}

// UpdateStatus 更新文章状态
func (h *PostHandler) UpdateStatus(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.NewResponse(http.StatusBadRequest, "参数错误", nil))
		return
	}

	var req request.UpdatePostStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.NewResponse(http.StatusBadRequest, err.Error(), nil))
		return
	}

	post := &models.Post{
		ID:     uint(id),
		Status: req.Status,
	}

	if err := h.postService.UpdatePost(c, post, nil); err != nil {
		c.JSON(http.StatusInternalServerError, response.NewResponse(http.StatusInternalServerError, err.Error(), nil))
		return
	}

	c.JSON(http.StatusOK, response.NewResponse(http.StatusOK, "更新成功", nil))
}
