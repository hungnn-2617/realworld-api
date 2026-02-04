package repositories

import (
	"realworld-api/internal/models"

	"gorm.io/gorm"
)

type CommentRepository struct {
	db *gorm.DB
}

func NewCommentRepository(db *gorm.DB) *CommentRepository {
	return &CommentRepository{db: db}
}

func (r *CommentRepository) Create(comment *models.Comment) error {
	return r.db.Create(comment).Error
}

func (r *CommentRepository) FindByID(id uint) (*models.Comment, error) {
	var comment models.Comment
	if err := r.db.Preload("Author").First(&comment, id).Error; err != nil {
		return nil, err
	}
	return &comment, nil
}

func (r *CommentRepository) FindByArticleID(articleID uint) ([]models.Comment, error) {
	var comments []models.Comment
	if err := r.db.Preload("Author").
		Where("article_id = ?", articleID).
		Order("created_at DESC").
		Find(&comments).Error; err != nil {
		return nil, err
	}
	return comments, nil
}

func (r *CommentRepository) Delete(comment *models.Comment) error {
	return r.db.Delete(comment).Error
}

func (r *CommentRepository) DeleteByID(id uint) error {
	return r.db.Delete(&models.Comment{}, id).Error
}
