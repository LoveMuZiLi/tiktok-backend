package model

import "time"

type Video struct {
	ID          int64     `gorm:"primaryKey;autoIncrement" json:"id"`
	Username    string    `gorm:"size:128;not null;index" json:"username"`
	Avatar      string    `gorm:"size:512;not null" json:"avatar"`
	Description string    `gorm:"type:text;not null" json:"description"`
	Music       string    `gorm:"size:256;not null" json:"music"`
	Image       string    `gorm:"size:512;not null" json:"image"`
	Likes       int64     `gorm:"not null;default:0" json:"likes"`
	Comments    int64     `gorm:"not null;default:0" json:"comments"`
	Shares      int64     `gorm:"not null;default:0" json:"shares"`
	CreatedAt   time.Time `json:"-"`
	UpdatedAt   time.Time `json:"-"`
}

func (Video) TableName() string {
	return "videos"
}

type VideoFeedResponse struct {
	Items  []Video `json:"items"`
	Total  int64   `json:"total"`
	Offset int     `json:"offset"`
	Limit  int     `json:"limit"`
}
