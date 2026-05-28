package entity

import "time"

type Follow struct {
	ID          int64     `gorm:"primaryKey;autoIncrement" json:"id"`
	FollowerID  int64     `gorm:"not null;uniqueIndex:idx_follow_pair" json:"followerId"`
	FollowingID int64     `gorm:"not null;uniqueIndex:idx_follow_pair" json:"followingId"`
	CreatedAt   time.Time `json:"createdAt"`
}

func (Follow) TableName() string { return "follows" }
