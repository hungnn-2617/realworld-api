package models

import (
	"time"
)

type Favorite struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	UserID    uint      `gorm:"not null;index" json:"userId"`
	ArticleID uint      `gorm:"not null;index" json:"articleId"`
	CreatedAt time.Time `json:"createdAt"`

	User    User    `gorm:"foreignKey:UserID" json:"-"`
	Article Article `gorm:"foreignKey:ArticleID" json:"-"`
}

func (Favorite) TableName() string {
	return "favorites"
}
