package models

import (
	"time"
)

type Follow struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	FollowerID  uint      `gorm:"not null;index" json:"followerId"`
	FollowingID uint      `gorm:"not null;index" json:"followingId"`
	CreatedAt   time.Time `json:"createdAt"`

	Follower  User `gorm:"foreignKey:FollowerID" json:"-"`
	Following User `gorm:"foreignKey:FollowingID" json:"-"`
}

func (Follow) TableName() string {
	return "follows"
}
