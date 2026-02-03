package repositories

import (
	"fmt"
	"time"

	"realworld-api/internal/models"
	"realworld-api/internal/utils"

	"gorm.io/gorm"
)

type ArticleRepository struct {
	db *gorm.DB
}

func NewArticleRepository(db *gorm.DB) *ArticleRepository {
	return &ArticleRepository{db: db}
}

func (r *ArticleRepository) Create(article *models.Article) error {
	// Generate unique slug
	baseSlug := utils.GenerateSlug(article.Title)
	article.Slug = baseSlug

	// Check if slug exists, append timestamp if needed
	var count int64
	r.db.Model(&models.Article{}).Where("slug = ?", article.Slug).Count(&count)
	if count > 0 {
		article.Slug = fmt.Sprintf("%s-%d", baseSlug, time.Now().UnixNano())
	}

	// Handle tags - find or create each tag
	if len(article.Tags) > 0 {
		var tags []models.Tag
		for _, tag := range article.Tags {
			var existingTag models.Tag
			result := r.db.Where("name = ?", tag.Name).First(&existingTag)
			if result.Error != nil {
				// Tag doesn't exist, create it
				newTag := models.Tag{Name: tag.Name}
				if err := r.db.Create(&newTag).Error; err != nil {
					return err
				}
				tags = append(tags, newTag)
			} else {
				tags = append(tags, existingTag)
			}
		}
		article.Tags = tags
	}

	return r.db.Create(article).Error
}

func (r *ArticleRepository) FindBySlug(slug string) (*models.Article, error) {
	var article models.Article
	if err := r.db.Preload("Author").Preload("Tags").Where("slug = ?", slug).First(&article).Error; err != nil {
		return nil, err
	}
	return &article, nil
}

func (r *ArticleRepository) FindByID(id uint) (*models.Article, error) {
	var article models.Article
	if err := r.db.Preload("Author").Preload("Tags").First(&article, id).Error; err != nil {
		return nil, err
	}
	return &article, nil
}

func (r *ArticleRepository) Update(article *models.Article) error {
	return r.db.Save(article).Error
}

func (r *ArticleRepository) Delete(article *models.Article) error {
	// Delete associated records first
	r.db.Where("article_id = ?", article.ID).Delete(&models.Comment{})
	r.db.Where("article_id = ?", article.ID).Delete(&models.Favorite{})
	r.db.Model(article).Association("Tags").Clear()

	return r.db.Delete(article).Error
}

type ArticleFilter struct {
	Tag       string
	Author    string
	Favorited string
	Limit     int
	Offset    int
}

func (r *ArticleRepository) List(filter ArticleFilter) ([]models.Article, int64, error) {
	var articles []models.Article
	var total int64

	query := r.db.Model(&models.Article{}).Preload("Author").Preload("Tags")

	// Apply filters
	if filter.Tag != "" {
		query = query.Joins("JOIN article_tags ON article_tags.article_id = articles.id").
			Joins("JOIN tags ON tags.id = article_tags.tag_id").
			Where("tags.name = ?", filter.Tag)
	}

	if filter.Author != "" {
		query = query.Joins("JOIN users ON users.id = articles.author_id").
			Where("users.username = ?", filter.Author)
	}

	if filter.Favorited != "" {
		query = query.Joins("JOIN favorites ON favorites.article_id = articles.id").
			Joins("JOIN users AS fav_users ON fav_users.id = favorites.user_id").
			Where("fav_users.username = ?", filter.Favorited)
	}

	// Get total count
	query.Count(&total)

	// Apply pagination
	if filter.Limit <= 0 {
		filter.Limit = 20
	}
	if filter.Limit > 100 {
		filter.Limit = 100
	}

	err := query.Order("articles.created_at DESC").
		Offset(filter.Offset).
		Limit(filter.Limit).
		Find(&articles).Error

	return articles, total, err
}

func (r *ArticleRepository) Feed(followingIDs []uint, limit, offset int) ([]models.Article, int64, error) {
	var articles []models.Article
	var total int64

	if len(followingIDs) == 0 {
		return articles, 0, nil
	}

	query := r.db.Model(&models.Article{}).
		Preload("Author").
		Preload("Tags").
		Where("author_id IN ?", followingIDs)

	query.Count(&total)

	if limit <= 0 {
		limit = 20
	}
	if limit > 100 {
		limit = 100
	}

	err := query.Order("created_at DESC").
		Offset(offset).
		Limit(limit).
		Find(&articles).Error

	return articles, total, err
}

func (r *ArticleRepository) UpdateSlug(article *models.Article, newTitle string) error {
	baseSlug := utils.GenerateSlug(newTitle)
	newSlug := baseSlug

	var count int64
	r.db.Model(&models.Article{}).Where("slug = ? AND id != ?", newSlug, article.ID).Count(&count)
	if count > 0 {
		newSlug = fmt.Sprintf("%s-%d", baseSlug, time.Now().UnixNano())
	}

	article.Slug = newSlug
	article.Title = newTitle
	return nil
}
