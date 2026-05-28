package entity

import "time"

// VideoLike 视频点赞记录
type VideoLike struct {
	ID        int64     `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID    int64     `gorm:"not null;uniqueIndex:idx_video_like_user_video" json:"userId"`
	VideoID   int64     `gorm:"not null;uniqueIndex:idx_video_like_user_video;index" json:"videoId"`
	CreatedAt time.Time `json:"createdAt"`
}

func (VideoLike) TableName() string { return "video_likes" }

// VideoFavorite 视频收藏记录
type VideoFavorite struct {
	ID        int64     `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID    int64     `gorm:"not null;uniqueIndex:idx_video_fav_user_video" json:"userId"`
	VideoID   int64     `gorm:"not null;uniqueIndex:idx_video_fav_user_video;index" json:"videoId"`
	CreatedAt time.Time `json:"createdAt"`
}

func (VideoFavorite) TableName() string { return "video_favorites" }

// VideoShare 视频分享记录（每次分享一条）
type VideoShare struct {
	ID        int64     `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID    int64     `gorm:"not null;index" json:"userId"`
	VideoID   int64     `gorm:"not null;index" json:"videoId"`
	CreatedAt time.Time `json:"createdAt"`
}

func (VideoShare) TableName() string { return "video_shares" }

// VideoComment 视频评论
type VideoComment struct {
	ID        int64     `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID    int64     `gorm:"not null;index" json:"userId"`
	VideoID   int64     `gorm:"not null;index" json:"videoId"`
	Content   string    `gorm:"type:text;not null" json:"content"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`

	User User `gorm:"foreignKey:UserID" json:"-"`
}

func (VideoComment) TableName() string { return "video_comments" }

type CommentView struct {
	ID        int64  `json:"id"`
	UserID    int64  `json:"userId"`
	VideoID   int64  `json:"videoId"`
	Username  string `json:"username"`
	Avatar    string `json:"avatar"`
	Content   string `json:"content"`
	CreatedAt string `json:"createdAt"`
}

func ToCommentView(c VideoComment) CommentView {
	username := c.User.Username
	if c.User.Nickname != "" {
		username = c.User.Nickname
	}
	if username == "" {
		username = "user"
	}
	return CommentView{
		ID:        c.ID,
		UserID:    c.UserID,
		VideoID:   c.VideoID,
		Username:  username,
		Avatar:    c.User.Avatar,
		Content:   c.Content,
		CreatedAt: c.CreatedAt.Format("2006-01-02 15:04"),
	}
}

type CreateCommentInput struct {
	UserID  int64  `json:"userId"`
	Content string `json:"content"`
}

type InteractionStatus struct {
	Liked     bool `json:"liked"`
	Favorited bool `json:"favorited"`
}
