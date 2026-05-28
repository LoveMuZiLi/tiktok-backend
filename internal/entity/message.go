package entity

import "time"

type Message struct {
	ID             int64     `gorm:"primaryKey;autoIncrement" json:"id"`
	ConversationID int64     `gorm:"not null;index" json:"conversationId"`
	SenderID       int64     `gorm:"not null;index" json:"senderId"`
	Text           string    `gorm:"type:text;not null" json:"text"`
	CreatedAt      time.Time `json:"createdAt"`
}

func (Message) TableName() string { return "messages" }

// MessageView 聊天详情展示
type MessageView struct {
	ID      int64  `json:"id"`
	Text    string `json:"text"`
	IsMine  bool   `json:"isMine"`
	Time    string `json:"time"`
}
