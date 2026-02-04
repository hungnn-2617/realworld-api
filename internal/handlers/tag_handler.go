package handlers

import (
	"net/http"

	"realworld-api/internal/services"

	"github.com/gin-gonic/gin"
)

type TagHandler struct {
	tagService *services.TagService
}

func NewTagHandler(tagService *services.TagService) *TagHandler {
	return &TagHandler{
		tagService: tagService,
	}
}

func (h *TagHandler) GetTags(c *gin.Context) {
	tags, err := h.tagService.GetAllTags()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"errors": gin.H{
				"body": []string{"failed to fetch tags"},
			},
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"tags": tags,
	})
}
