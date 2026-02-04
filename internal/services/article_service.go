package services

import (
	"errors"
	"time"

	"realworld-api/internal/dtos"
	"realworld-api/internal/models"
	"realworld-api/internal/repositories"

	"gorm.io/gorm"
)

type ArticleService struct {
	articleRepo  *repositories.ArticleRepository
	favoriteRepo *repositories.FavoriteRepository
	userRepo     *repositories.UserRepository
}

func NewArticleService(db *gorm.DB) *ArticleService {
	return &ArticleService{
		articleRepo:  repositories.NewArticleRepository(db),
		favoriteRepo: repositories.NewFavoriteRepository(db),
		userRepo:     repositories.NewUserRepository(db),
	}
}

func (s *ArticleService) Create(req *dtos.CreateArticleRequest, authorID uint) (*models.Article, error) {
	// Create tags first if provided
	var tags []models.Tag
	if len(req.Article.TagList) > 0 {
		for _, tagName := range req.Article.TagList {
			tag := models.Tag{Name: tagName}
			tags = append(tags, tag)
		}
	}

	article := &models.Article{
		Title:       req.Article.Title,
		Description: req.Article.Description,
		Body:        req.Article.Body,
		AuthorID:    authorID,
		Tags:        tags,
	}

	if err := s.articleRepo.Create(article); err != nil {
		return nil, err
	}

	// Reload article with relations
	return s.articleRepo.FindByID(article.ID)
}

func (s *ArticleService) GetBySlug(slug string) (*models.Article, error) {
	article, err := s.articleRepo.FindBySlug(slug)
	if err != nil {
		return nil, errors.New("article not found")
	}
	return article, nil
}

func (s *ArticleService) Update(slug string, req *dtos.UpdateArticleRequest, userID uint) (*models.Article, error) {
	article, err := s.articleRepo.FindBySlug(slug)
	if err != nil {
		return nil, errors.New("article not found")
	}

	if article.AuthorID != userID {
		return nil, errors.New("not authorized to update this article")
	}

	if req.Article.Title != nil {
		if err := s.articleRepo.UpdateSlug(article, *req.Article.Title); err != nil {
			return nil, errors.New("failed to update slug")
		}
	}

	if req.Article.Description != nil {
		article.Description = *req.Article.Description
	}

	if req.Article.Body != nil {
		article.Body = *req.Article.Body
	}

	if err := s.articleRepo.Update(article); err != nil {
		return nil, errors.New("failed to update article")
	}

	return s.articleRepo.FindByID(article.ID)
}

func (s *ArticleService) Delete(slug string, userID uint) error {
	article, err := s.articleRepo.FindBySlug(slug)
	if err != nil {
		return errors.New("article not found")
	}

	if article.AuthorID != userID {
		return errors.New("not authorized to delete this article")
	}

	return s.articleRepo.Delete(article)
}

func (s *ArticleService) List(filter repositories.ArticleFilter) ([]models.Article, int64, error) {
	return s.articleRepo.List(filter)
}

func (s *ArticleService) Feed(userID uint, limit, offset int) ([]models.Article, int64, error) {
	followingIDs, err := s.userRepo.GetFollowingIDs(userID)
	if err != nil {
		return nil, 0, err
	}

	return s.articleRepo.Feed(followingIDs, limit, offset)
}

func (s *ArticleService) Favorite(slug string, userID uint) (*models.Article, error) {
	article, err := s.articleRepo.FindBySlug(slug)
	if err != nil {
		return nil, errors.New("article not found")
	}

	if err := s.favoriteRepo.Favorite(userID, article.ID); err != nil {
		return nil, errors.New("failed to favorite article")
	}

	return s.articleRepo.FindByID(article.ID)
}

func (s *ArticleService) Unfavorite(slug string, userID uint) (*models.Article, error) {
	article, err := s.articleRepo.FindBySlug(slug)
	if err != nil {
		return nil, errors.New("article not found")
	}

	if err := s.favoriteRepo.Unfavorite(userID, article.ID); err != nil {
		return nil, errors.New("failed to unfavorite article")
	}

	return s.articleRepo.FindByID(article.ID)
}

func (s *ArticleService) IsFavorited(articleID, userID uint) bool {
	if userID == 0 {
		return false
	}
	return s.favoriteRepo.IsFavorited(userID, articleID)
}

func (s *ArticleService) IsFollowing(currentUserID, authorID uint) bool {
	if currentUserID == 0 {
		return false
	}
	return s.userRepo.IsFollowing(currentUserID, authorID)
}

// ToArticleData converts Article model to DTO
func (s *ArticleService) ToArticleData(article *models.Article, currentUserID uint) dtos.ArticleData {
	tagList := make([]string, len(article.Tags))
	for i, tag := range article.Tags {
		tagList[i] = tag.Name
	}

	return dtos.ArticleData{
		Slug:           article.Slug,
		Title:          article.Title,
		Description:    article.Description,
		Body:           article.Body,
		TagList:        tagList,
		CreatedAt:      article.CreatedAt.Format(time.RFC3339),
		UpdatedAt:      article.UpdatedAt.Format(time.RFC3339),
		Favorited:      s.IsFavorited(article.ID, currentUserID),
		FavoritesCount: article.FavoritesCount,
		Author: dtos.ProfileData{
			Username:  article.Author.Username,
			Bio:       article.Author.Bio,
			Image:     article.Author.Image,
			Following: s.IsFollowing(currentUserID, article.AuthorID),
		},
	}
}

func (s *ArticleService) ToArticlesData(articles []models.Article, currentUserID uint) []dtos.ArticleData {
	result := make([]dtos.ArticleData, len(articles))
	for i, article := range articles {
		result[i] = s.ToArticleData(&article, currentUserID)
	}
	return result
}
