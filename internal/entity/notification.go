package entity

import "time"

const (
	NotificationTypeFollow  = "follow"
	NotificationTypeActivity  = "activity"
	NotificationTypeSystem    = "system"
)

type Notification struct {
	ID        int64     `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID    int64     `gorm:"not null;index" json:"userId"`
	Type      string    `gorm:"size:32;not null;index" json:"type"`
	Title     string    `gorm:"size:128;not null" json:"title"`
	Body      string    `gorm:"type:text;not null" json:"body"`
	IsRead    bool      `gorm:"not null;default:false" json:"isRead"`
	CreatedAt time.Time `json:"createdAt"`
}

func (Notification) TableName() string { return "notifications" }
