package models

import (
	"time"
)

type Tag struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Name      string    `gorm:"uniqueIndex;size:100;not null" json:"name"`
	CreatedAt time.Time `json:"createdAt"`

	Articles []Article `gorm:"many2many:article_tags;" json:"-"`
}
