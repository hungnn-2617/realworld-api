package repositories

import (
	"realworld-api/internal/models"

	"gorm.io/gorm"
)

type TagRepository struct {
	db *gorm.DB
}

func NewTagRepository(db *gorm.DB) *TagRepository {
	return &TagRepository{db: db}
}

// GetAll returns all tags from the database
func (r *TagRepository) GetAll() ([]models.Tag, error) {
	var tags []models.Tag
	if err := r.db.Order("name ASC").Find(&tags).Error; err != nil {
		return nil, err
	}
	return tags, nil
}

// FindByName finds a tag by its name
func (r *TagRepository) FindByName(name string) (*models.Tag, error) {
	var tag models.Tag
	if err := r.db.Where("name = ?", name).First(&tag).Error; err != nil {
		return nil, err
	}
	return &tag, nil
}

// Create creates a new tag
func (r *TagRepository) Create(tag *models.Tag) error {
	return r.db.Create(tag).Error
}

// FindOrCreate finds a tag by name or creates it if it doesn't exist
func (r *TagRepository) FindOrCreate(name string) (*models.Tag, error) {
	var tag models.Tag
	result := r.db.Where("name = ?", name).First(&tag)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			tag = models.Tag{Name: name}
			if err := r.db.Create(&tag).Error; err != nil {
				return nil, err
			}
			return &tag, nil
		}
		return nil, result.Error
	}
	return &tag, nil
}
