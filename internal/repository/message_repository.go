package repository

import (
	"errors"
	"time"

	"gorm.io/gorm"

	"github.com/LoveMuZiLi/tiktok-backend/internal/entity"
)

type MessageRepository struct {
	db *gorm.DB
}

func NewMessageRepository(db *gorm.DB) *MessageRepository {
	return &MessageRepository{db: db}
}

func (r *MessageRepository) List(convID, viewerID int64) ([]entity.MessageView, error) {
	var msgs []entity.Message
	err := r.db.Where("conversation_id = ?", convID).Order("id ASC").Find(&msgs).Error
	if err != nil {
		return nil, err
	}
	out := make([]entity.MessageView, len(msgs))
	for i, m := range msgs {
		out[i] = entity.MessageView{
			ID:     m.ID,
			Text:   m.Text,
			IsMine: m.SenderID == viewerID,
			Time:   m.CreatedAt.Format("15:04"),
		}
	}
	return out, nil
}

func (r *MessageRepository) Create(convID, senderID int64, text string) (entity.Message, error) {
	m := entity.Message{
		ConversationID: convID,
		SenderID:       senderID,
		Text:           text,
		CreatedAt:      time.Now(),
	}
	err := r.db.Create(&m).Error
	return m, err
}

func (r *MessageRepository) Delete(id int64) error {
	res := r.db.Delete(&entity.Message{}, id)
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return errors.New("message not found")
	}
	return nil
}

func (r *MessageRepository) IncrementUnread(convID, exceptUserID int64) error {
	return r.db.Model(&entity.ConversationMember{}).
		Where("conversation_id = ? AND user_id != ?", convID, exceptUserID).
		UpdateColumn("unread_count", gorm.Expr("unread_count + 1")).Error
}

func (r *MessageRepository) ClearUnread(convID, userID int64) error {
	return r.db.Model(&entity.ConversationMember{}).
		Where("conversation_id = ? AND user_id = ?", convID, userID).
		Update("unread_count", 0).Error
}
