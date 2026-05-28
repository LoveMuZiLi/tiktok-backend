package entity

import "time"

type User struct {
	ID        int64     `gorm:"primaryKey;autoIncrement" json:"id"`
	Username  string    `gorm:"size:64;not null;uniqueIndex" json:"username"`
	Nickname  string    `gorm:"size:128;not null" json:"nickname"`
	Avatar    string    `gorm:"size:512;not null" json:"avatar"`
	DouyinID  string    `gorm:"size:32;not null;uniqueIndex" json:"douyinId"`
	Bio       string    `gorm:"type:text" json:"bio"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

func (User) TableName() string { return "users" }

type UserProfile struct {
	User            User    `json:"user"`
	FollowingCount  int64   `json:"followingCount"`
	FollowerCount   int64   `json:"followerCount"`
	TotalLikes      int64   `json:"totalLikes"`
	VideoCount      int64   `json:"videoCount"`
}
