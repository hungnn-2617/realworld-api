package handlers

import (
	"net/http"

	"realworld-api/internal/dtos"
	"realworld-api/internal/helpers"
	"realworld-api/internal/middleware"
	"realworld-api/internal/services"

	"github.com/gin-gonic/gin"
)

type FavoriteHandler struct {
	articleService *services.ArticleService
}

func NewFavoriteHandler(articleService *services.ArticleService) *FavoriteHandler {
	return &FavoriteHandler{
		articleService: articleService,
	}
}

// FavoriteArticle handles POST /api/articles/:slug/favorite
func (h *FavoriteHandler) FavoriteArticle(c *gin.Context) {
	userID, exists := middleware.GetUserID(c)
	if !exists {
		helpers.Unauthorized(c)
		return
	}

	slug := c.Param("slug")

	article, err := h.articleService.Favorite(slug, userID)
	if err != nil {
		if err.Error() == "article not found" {
			helpers.NotFound(c, err.Error())
		} else {
			helpers.UnprocessableEntity(c, err.Error())
		}
		return
	}

	c.JSON(http.StatusOK, dtos.ArticleResponse{
		Article: h.articleService.ToArticleData(article, userID),
	})
}

// UnfavoriteArticle handles DELETE /api/articles/:slug/favorite
func (h *FavoriteHandler) UnfavoriteArticle(c *gin.Context) {
	userID, exists := middleware.GetUserID(c)
	if !exists {
		helpers.Unauthorized(c)
		return
	}

	slug := c.Param("slug")

	article, err := h.articleService.Unfavorite(slug, userID)
	if err != nil {
		if err.Error() == "article not found" {
			helpers.NotFound(c, err.Error())
		} else {
			helpers.UnprocessableEntity(c, err.Error())
		}
		return
	}

	c.JSON(http.StatusOK, dtos.ArticleResponse{
		Article: h.articleService.ToArticleData(article, userID),
	})
}
