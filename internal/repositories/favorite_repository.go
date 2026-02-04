package repositories

import (
	"realworld-api/internal/models"

	"gorm.io/gorm"
)

type FavoriteRepository struct {
	db *gorm.DB
}

func NewFavoriteRepository(db *gorm.DB) *FavoriteRepository {
	return &FavoriteRepository{db: db}
}

func (r *FavoriteRepository) Favorite(userID, articleID uint) error {
	// Check if already favorited
	var count int64
	r.db.Model(&models.Favorite{}).
		Where("user_id = ? AND article_id = ?", userID, articleID).
		Count(&count)

	if count > 0 {
		return nil // Already favorited
	}

	// Create favorite
	favorite := &models.Favorite{
		UserID:    userID,
		ArticleID: articleID,
	}

	if err := r.db.Create(favorite).Error; err != nil {
		return err
	}

	// Update article favorites count
	return r.db.Model(&models.Article{}).
		Where("id = ?", articleID).
		UpdateColumn("favorites_count", gorm.Expr("favorites_count + ?", 1)).Error
}

func (r *FavoriteRepository) Unfavorite(userID, articleID uint) error {
	// Check if favorited
	var favorite models.Favorite
	result := r.db.Where("user_id = ? AND article_id = ?", userID, articleID).First(&favorite)
	if result.Error != nil {
		return nil // Not favorited, nothing to do
	}

	// Delete favorite
	if err := r.db.Delete(&favorite).Error; err != nil {
		return err
	}

	// Update article favorites count
	return r.db.Model(&models.Article{}).
		Where("id = ?", articleID).
		UpdateColumn("favorites_count", gorm.Expr("GREATEST(favorites_count - 1, 0)")).Error
}

func (r *FavoriteRepository) IsFavorited(userID, articleID uint) bool {
	var count int64
	r.db.Model(&models.Favorite{}).
		Where("user_id = ? AND article_id = ?", userID, articleID).
		Count(&count)
	return count > 0
}

func (r *FavoriteRepository) GetFavoritesByUserID(userID uint) ([]models.Favorite, error) {
	var favorites []models.Favorite
	if err := r.db.Where("user_id = ?", userID).Find(&favorites).Error; err != nil {
		return nil, err
	}
	return favorites, nil
}
