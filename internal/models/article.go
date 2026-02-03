package models

import (
	"time"

	"gorm.io/gorm"
)

type Article struct {
	ID             uint           `gorm:"primaryKey" json:"id"`
	Slug           string         `gorm:"uniqueIndex;size:255;not null" json:"slug"`
	Title          string         `gorm:"size:255;not null" json:"title"`
	Description    string         `gorm:"type:text" json:"description"`
	Body           string         `gorm:"type:text;not null" json:"body"`
	AuthorID       uint           `gorm:"not null;index" json:"authorId"`
	FavoritesCount int            `gorm:"default:0" json:"favoritesCount"`
	CreatedAt      time.Time      `json:"createdAt"`
	UpdatedAt      time.Time      `json:"updatedAt"`
	DeletedAt      gorm.DeletedAt `gorm:"index" json:"-"`

	// Relations
	Author    User       `gorm:"foreignKey:AuthorID" json:"author"`
	Tags      []Tag      `gorm:"many2many:article_tags;" json:"tagList"`
	Comments  []Comment  `gorm:"foreignKey:ArticleID" json:"-"`
	Favorites []Favorite `gorm:"foreignKey:ArticleID" json:"-"`
}
