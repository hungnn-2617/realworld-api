package services

import (
	"realworld-api/internal/repositories"

	"gorm.io/gorm"
)

type TagService struct {
	tagRepo *repositories.TagRepository
}

func NewTagService(db *gorm.DB) *TagService {
	return &TagService{
		tagRepo: repositories.NewTagRepository(db),
	}
}

// GetAllTags returns all tag names
func (s *TagService) GetAllTags() ([]string, error) {
	tags, err := s.tagRepo.GetAll()
	if err != nil {
		return nil, err
	}

	tagNames := make([]string, len(tags))
	for i, tag := range tags {
		tagNames[i] = tag.Name
	}

	return tagNames, nil
}
