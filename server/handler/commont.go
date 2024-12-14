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

type CommentHandler struct {
	commentService service.CommentService
}

func NewCommentHandler(commentService service.CommentService) *CommentHandler {
	return &CommentHandler{
		commentService: commentService,
	}
}

// Create 创建评论
func (h *CommentHandler) Create(c *gin.Context) {
	var req request.CreateCommentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.NewResponse(http.StatusBadRequest, "参数错误", nil))
		return
	}

	userID, _ := c.Get("userID")
	comment := &models.Comment{
		PostID:   req.PostID,
		UserID:   userID.(uint),
		Content:  req.Content,
		ParentID: req.ParentID,
	}

	if err := h.commentService.CreateComment(c, comment); err != nil {
		c.JSON(http.StatusInternalServerError, response.NewResponse(http.StatusInternalServerError, err.Error(), nil))
		return
	}

	c.JSON(http.StatusOK, response.NewResponse(http.StatusOK, "评论成功", comment))
}

// ListByPost 获取文章评论列表
func (h *CommentHandler) ListByPost(c *gin.Context) {
	postID, err := strconv.ParseUint(c.Param("postID"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.NewResponse(http.StatusBadRequest, "参数错误", nil))
		return
	}

	var req request.ListCommentsRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.NewResponse(http.StatusBadRequest, "参数错误", nil))
		return
	}

	comments, total, err := h.commentService.ListCommentsByPost(c, uint(postID), req.Page, req.PageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.NewResponse(http.StatusInternalServerError, err.Error(), nil))
		return
	}

	c.JSON(http.StatusOK, response.NewResponse(http.StatusOK, "获取成功", response.NewPaginationResponse(comments, total, req.Page, req.PageSize)))
}

// Delete 删除评论
func (h *CommentHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.NewResponse(http.StatusBadRequest, "参数错误", nil))
		return
	}

	if err := h.commentService.DeleteComment(c, uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, response.NewResponse(http.StatusInternalServerError, err.Error(), nil))
		return
	}

	c.JSON(http.StatusOK, response.NewResponse(http.StatusOK, "删除成功", nil))
}