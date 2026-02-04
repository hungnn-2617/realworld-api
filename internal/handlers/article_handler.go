package handlers

import (
	"net/http"
	"strconv"

	"realworld-api/internal/dtos"
	"realworld-api/internal/helpers"
	"realworld-api/internal/middleware"
	"realworld-api/internal/repositories"
	"realworld-api/internal/services"

	"github.com/gin-gonic/gin"
)

type ArticleHandler struct {
	articleService *services.ArticleService
}

func NewArticleHandler(articleService *services.ArticleService) *ArticleHandler {
	return &ArticleHandler{
		articleService: articleService,
	}
}

// CreateArticle handles POST /api/articles
func (h *ArticleHandler) CreateArticle(c *gin.Context) {
	userID, exists := middleware.GetUserID(c)
	if !exists {
		helpers.Unauthorized(c)
		return
	}

	var req dtos.CreateArticleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		helpers.ValidationError(c, err)
		return
	}

	article, err := h.articleService.Create(&req, userID)
	if err != nil {
		helpers.UnprocessableEntity(c, err.Error())
		return
	}

	c.JSON(http.StatusCreated, dtos.ArticleResponse{
		Article: h.articleService.ToArticleData(article, userID),
	})
}

// GetArticle handles GET /api/articles/:slug
func (h *ArticleHandler) GetArticle(c *gin.Context) {
	slug := c.Param("slug")

	article, err := h.articleService.GetBySlug(slug)
	if err != nil {
		helpers.NotFound(c, "article not found")
		return
	}

	// Get current user ID if authenticated (optional)
	userID, _ := middleware.GetUserID(c)

	c.JSON(http.StatusOK, dtos.ArticleResponse{
		Article: h.articleService.ToArticleData(article, userID),
	})
}

// UpdateArticle handles PUT /api/articles/:slug
func (h *ArticleHandler) UpdateArticle(c *gin.Context) {
	userID, exists := middleware.GetUserID(c)
	if !exists {
		helpers.Unauthorized(c)
		return
	}

	slug := c.Param("slug")

	var req dtos.UpdateArticleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		helpers.ValidationError(c, err)
		return
	}

	article, err := h.articleService.Update(slug, &req, userID)
	if err != nil {
		if err.Error() == "article not found" {
			helpers.NotFound(c, err.Error())
		} else if err.Error() == "not authorized to update this article" {
			helpers.Forbidden(c, err.Error())
		} else {
			helpers.UnprocessableEntity(c, err.Error())
		}
		return
	}

	c.JSON(http.StatusOK, dtos.ArticleResponse{
		Article: h.articleService.ToArticleData(article, userID),
	})
}

// DeleteArticle handles DELETE /api/articles/:slug
func (h *ArticleHandler) DeleteArticle(c *gin.Context) {
	userID, exists := middleware.GetUserID(c)
	if !exists {
		helpers.Unauthorized(c)
		return
	}

	slug := c.Param("slug")

	err := h.articleService.Delete(slug, userID)
	if err != nil {
		if err.Error() == "article not found" {
			helpers.NotFound(c, err.Error())
		} else if err.Error() == "not authorized to delete this article" {
			helpers.Forbidden(c, err.Error())
		} else {
			helpers.UnprocessableEntity(c, err.Error())
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "article deleted successfully",
	})
}

// ListArticles handles GET /api/articles
func (h *ArticleHandler) ListArticles(c *gin.Context) {
	var filter repositories.ArticleFilter

	// Parse query parameters
	filter.Tag = c.Query("tag")
	filter.Author = c.Query("author")
	filter.Favorited = c.Query("favorited")

	// Parse limit (default: 20)
	limit := 20
	if limitStr := c.Query("limit"); limitStr != "" {
		if parsedLimit, err := strconv.Atoi(limitStr); err == nil && parsedLimit > 0 {
			limit = parsedLimit
		}
	}
	filter.Limit = limit

	// Parse offset (default: 0)
	offset := 0
	if offsetStr := c.Query("offset"); offsetStr != "" {
		if parsedOffset, err := strconv.Atoi(offsetStr); err == nil && parsedOffset >= 0 {
			offset = parsedOffset
		}
	}
	filter.Offset = offset

	articles, count, err := h.articleService.List(filter)
	if err != nil {
		helpers.InternalServerError(c, "failed to fetch articles")
		return
	}

	// Get current user ID if authenticated (optional)
	userID, _ := middleware.GetUserID(c)

	c.JSON(http.StatusOK, dtos.ArticlesResponse{
		Articles:      h.articleService.ToArticlesData(articles, userID),
		ArticlesCount: int(count),
	})
}

// FeedArticles handles GET /api/articles/feed
func (h *ArticleHandler) FeedArticles(c *gin.Context) {
	userID, exists := middleware.GetUserID(c)
	if !exists {
		helpers.Unauthorized(c)
		return
	}

	// Parse limit (default: 20)
	limit := 20
	if limitStr := c.Query("limit"); limitStr != "" {
		if parsedLimit, err := strconv.Atoi(limitStr); err == nil && parsedLimit > 0 {
			limit = parsedLimit
		}
	}

	// Parse offset (default: 0)
	offset := 0
	if offsetStr := c.Query("offset"); offsetStr != "" {
		if parsedOffset, err := strconv.Atoi(offsetStr); err == nil && parsedOffset >= 0 {
			offset = parsedOffset
		}
	}

	articles, count, err := h.articleService.Feed(userID, limit, offset)
	if err != nil {
		helpers.InternalServerError(c, "failed to fetch feed")
		return
	}

	c.JSON(http.StatusOK, dtos.ArticlesResponse{
		Articles:      h.articleService.ToArticlesData(articles, userID),
		ArticlesCount: int(count),
	})
}
