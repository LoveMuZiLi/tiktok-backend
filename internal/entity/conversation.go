package entity

import "time"

type Conversation struct {
	ID        int64     `gorm:"primaryKey;autoIncrement" json:"id"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

func (Conversation) TableName() string { return "conversations" }

type ConversationMember struct {
	ID             int64     `gorm:"primaryKey;autoIncrement" json:"id"`
	ConversationID int64     `gorm:"not null;uniqueIndex:idx_conv_user" json:"conversationId"`
	UserID         int64     `gorm:"not null;uniqueIndex:idx_conv_user" json:"userId"`
	UnreadCount    int       `gorm:"not null;default:0" json:"unreadCount"`
	CreatedAt      time.Time `json:"createdAt"`
}

func (ConversationMember) TableName() string { return "conversation_members" }

// ChatSummary 收件箱会话列表项
type ChatSummary struct {
	ID          int64  `json:"id"`
	PeerID      int64  `json:"peerId"`
	Name        string `json:"name"`
	Avatar      string `json:"avatar"`
	LastMessage string `json:"message"`
	Time        string `json:"time"`
	Unread      int    `json:"unread"`
	HasCamera   bool   `json:"hasCamera"`
}
