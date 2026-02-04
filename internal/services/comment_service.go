package services

import (
	"errors"
	"time"

	"realworld-api/internal/dtos"
	appErrors "realworld-api/internal/errors"
	"realworld-api/internal/models"
	"realworld-api/internal/repositories"

	"gorm.io/gorm"
)

type CommentService struct {
	commentRepo *repositories.CommentRepository
	articleRepo *repositories.ArticleRepository
	userRepo    *repositories.UserRepository
}

func NewCommentService(db *gorm.DB) *CommentService {
	return &CommentService{
		commentRepo: repositories.NewCommentRepository(db),
		articleRepo: repositories.NewArticleRepository(db),
		userRepo:    repositories.NewUserRepository(db),
	}
}

func (s *CommentService) AddComment(slug string, req *dtos.CreateCommentRequest, authorID uint) (*models.Comment, error) {
	// Find article by slug
	article, err := s.articleRepo.FindBySlug(slug)
	if err != nil {
		return nil, errors.New(appErrors.ErrArticleNotFound)
	}

	// Create comment
	comment := &models.Comment{
		Body:      req.Comment.Body,
		ArticleID: article.ID,
		AuthorID:  authorID,
	}

	if err := s.commentRepo.Create(comment); err != nil {
		return nil, errors.New("failed to create comment")
	}

	// Reload comment with author
	return s.commentRepo.FindByID(comment.ID)
}

func (s *CommentService) GetComments(slug string) ([]models.Comment, error) {
	// Find article by slug
	article, err := s.articleRepo.FindBySlug(slug)
	if err != nil {
		return nil, errors.New(appErrors.ErrArticleNotFound)
	}

	return s.commentRepo.FindByArticleID(article.ID)
}

func (s *CommentService) DeleteComment(slug string, commentID uint, userID uint) error {
	// Find article by slug
	_, err := s.articleRepo.FindBySlug(slug)
	if err != nil {
		return errors.New(appErrors.ErrArticleNotFound)
	}

	// Find comment
	comment, err := s.commentRepo.FindByID(commentID)
	if err != nil {
		return errors.New(appErrors.ErrCommentNotFound)
	}

	// Check if user is the author of the comment
	if comment.AuthorID != userID {
		return errors.New(appErrors.ErrNotAuthorizedDeleteComment)
	}

	return s.commentRepo.Delete(comment)
}

func (s *CommentService) IsFollowing(currentUserID, authorID uint) bool {
	if currentUserID == 0 {
		return false
	}
	return s.userRepo.IsFollowing(currentUserID, authorID)
}

// ToCommentData converts Comment model to DTO
func (s *CommentService) ToCommentData(comment *models.Comment, currentUserID uint) dtos.CommentData {
	return dtos.CommentData{
		ID:        comment.ID,
		Body:      comment.Body,
		CreatedAt: comment.CreatedAt.Format(time.RFC3339),
		UpdatedAt: comment.UpdatedAt.Format(time.RFC3339),
		Author: dtos.ProfileData{
			Username:  comment.Author.Username,
			Bio:       comment.Author.Bio,
			Image:     comment.Author.Image,
			Following: s.IsFollowing(currentUserID, comment.AuthorID),
		},
	}
}

func (s *CommentService) ToCommentsData(comments []models.Comment, currentUserID uint) []dtos.CommentData {
	result := make([]dtos.CommentData, len(comments))
	for i, comment := range comments {
		result[i] = s.ToCommentData(&comment, currentUserID)
	}
	return result
}
