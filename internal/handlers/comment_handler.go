package handlers

import (
	"net/http"
	"strconv"

	"realworld-api/internal/dtos"
	appErrors "realworld-api/internal/errors"
	"realworld-api/internal/helpers"
	"realworld-api/internal/middleware"
	"realworld-api/internal/services"

	"github.com/gin-gonic/gin"
)

type CommentHandler struct {
	commentService *services.CommentService
}

func NewCommentHandler(commentService *services.CommentService) *CommentHandler {
	return &CommentHandler{
		commentService: commentService,
	}
}

// AddComment handles POST /api/articles/:slug/comments
func (h *CommentHandler) AddComment(c *gin.Context) {
	userID, exists := middleware.GetUserID(c)
	if !exists {
		helpers.Unauthorized(c)
		return
	}

	slug := c.Param("slug")

	var req dtos.CreateCommentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		helpers.ValidationError(c, err)
		return
	}

	comment, err := h.commentService.AddComment(slug, &req, userID)
	if err != nil {
		if appErrors.IsArticleNotFound(err) {
			helpers.NotFound(c, appErrors.ErrArticleNotFound)
		} else {
			helpers.UnprocessableEntity(c, err.Error())
		}
		return
	}

	c.JSON(http.StatusOK, dtos.CommentResponse{
		Comment: h.commentService.ToCommentData(comment, userID),
	})
}

// GetComments handles GET /api/articles/:slug/comments
func (h *CommentHandler) GetComments(c *gin.Context) {
	slug := c.Param("slug")

	comments, err := h.commentService.GetComments(slug)
	if err != nil {
		if appErrors.IsArticleNotFound(err) {
			helpers.NotFound(c, appErrors.ErrArticleNotFound)
		} else {
			helpers.InternalServerError(c, appErrors.ErrFailedToFetch)
		}
		return
	}

	// Get current user ID if authenticated (optional)
	userID, _ := middleware.GetUserID(c)

	c.JSON(http.StatusOK, dtos.CommentsResponse{
		Comments: h.commentService.ToCommentsData(comments, userID),
	})
}

// DeleteComment handles DELETE /api/articles/:slug/comments/:id
func (h *CommentHandler) DeleteComment(c *gin.Context) {
	userID, exists := middleware.GetUserID(c)
	if !exists {
		helpers.Unauthorized(c)
		return
	}

	slug := c.Param("slug")
	commentIDStr := c.Param("id")
	commentID, err := strconv.ParseUint(commentIDStr, 10, 32)
	if err != nil {
		helpers.BadRequest(c, appErrors.ErrInvalidCommentID)
		return
	}

	err = h.commentService.DeleteComment(slug, uint(commentID), userID)
	if err != nil {
		if appErrors.IsArticleNotFound(err) {
			helpers.NotFound(c, appErrors.ErrArticleNotFound)
		} else if appErrors.IsCommentNotFound(err) {
			helpers.NotFound(c, appErrors.ErrCommentNotFound)
		} else if appErrors.IsNotAuthorizedDeleteComment(err) {
			helpers.Forbidden(c, appErrors.ErrNotAuthorizedDeleteComment)
		} else {
			helpers.UnprocessableEntity(c, err.Error())
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "comment deleted successfully",
	})
}
