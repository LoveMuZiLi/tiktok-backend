package entity

import "time"

type Video struct {
	ID          int64     `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID      int64     `gorm:"not null;index" json:"userId"`
	Description string    `gorm:"type:text;not null" json:"description"`
	Music       string    `gorm:"size:256;not null" json:"music"`
	Image       string    `gorm:"size:512;not null" json:"image"`
	VideoURL    string    `gorm:"size:1024;not null;default:''" json:"videoUrl"`
	Likes       int64     `gorm:"not null;default:0" json:"likes"`
	Favorites   int64     `gorm:"not null;default:0" json:"favorites"`
	Comments    int64     `gorm:"not null;default:0" json:"comments"`
	Shares      int64     `gorm:"not null;default:0" json:"shares"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`

	User User `gorm:"foreignKey:UserID" json:"-"`
}

func (Video) TableName() string { return "videos" }

// VideoView API 响应（含作者信息，兼容前端）
type VideoView struct {
	ID          int64  `json:"id"`
	UserID      int64  `json:"userId"`
	Username    string `json:"username"`
	Avatar      string `json:"avatar"`
	Description string `json:"description"`
	Music       string `json:"music"`
	Image       string `json:"image"`
	VideoURL    string `json:"videoUrl"`
	Likes       int64  `json:"likes"`
	Favorites   int64  `json:"favorites"`
	Comments    int64  `json:"comments"`
	Shares      int64  `json:"shares"`
}

func ToVideoView(v Video) VideoView {
	username := v.User.Username
	avatar := v.User.Avatar
	if username == "" {
		username = "unknown"
	}
	return VideoView{
		ID:          v.ID,
		UserID:      v.UserID,
		Username:    username,
		Avatar:      avatar,
		Description: v.Description,
		Music:       v.Music,
		Image:       v.Image,
		VideoURL:    v.VideoURL,
		Likes:       v.Likes,
		Favorites:   v.Favorites,
		Comments:    v.Comments,
		Shares:      v.Shares,
	}
}

func ToVideoViews(items []Video) []VideoView {
	out := make([]VideoView, len(items))
	for i, v := range items {
		out[i] = ToVideoView(v)
	}
	return out
}

type VideoFeedResponse struct {
	Items  []VideoView `json:"items"`
	Total  int64       `json:"total"`
	Offset int         `json:"offset"`
	Limit  int         `json:"limit"`
}

type CreateVideoInput struct {
	UserID      int64  `json:"userId"`
	Description string `json:"description"`
	Music       string `json:"music"`
	Image       string `json:"image"`
}

type UpdateVideoInput struct {
	Description *string `json:"description"`
	Music       *string `json:"music"`
	Image       *string `json:"image"`
}
