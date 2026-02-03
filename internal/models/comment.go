package models

import (
	"time"

	"gorm.io/gorm"
)

type Comment struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	Body      string         `gorm:"type:text;not null" json:"body"`
	ArticleID uint           `gorm:"not null;index" json:"articleId"`
	AuthorID  uint           `gorm:"not null;index" json:"authorId"`
	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	Author  User    `gorm:"foreignKey:AuthorID" json:"author"`
	Article Article `gorm:"foreignKey:ArticleID" json:"-"`
}
