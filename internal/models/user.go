package models

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID           uint           `gorm:"primaryKey" json:"id"`
	Username     string         `gorm:"uniqueIndex;size:255;not null" json:"username"`
	Email        string         `gorm:"uniqueIndex;size:255;not null" json:"email"`
	PasswordHash string         `gorm:"size:255;not null" json:"-"`
	Bio          string         `gorm:"type:text" json:"bio"`
	Image        string         `gorm:"size:500" json:"image"`
	CreatedAt    time.Time      `json:"createdAt"`
	UpdatedAt    time.Time      `json:"updatedAt"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"-"`

	Followers []Follow `gorm:"foreignKey:FollowingID" json:"-"`
	Following []Follow `gorm:"foreignKey:FollowerID" json:"-"`
}
