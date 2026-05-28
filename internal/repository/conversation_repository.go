package repository

import (
	"errors"
	"fmt"
	"time"

	"gorm.io/gorm"

	"github.com/LoveMuZiLi/tiktok-backend/internal/entity"
)

var ErrConversationNotFound = errors.New("conversation not found")

type ConversationRepository struct {
	db *gorm.DB
}

func NewConversationRepository(db *gorm.DB) *ConversationRepository {
	return &ConversationRepository{db: db}
}

func (r *ConversationRepository) ListChats(userID int64) ([]entity.ChatSummary, error) {
	type row struct {
		ConversationID int64
		PeerID         int64
		Unread         int
		UpdatedAt      time.Time
	}
	var rows []row
	err := r.db.Raw(`
		SELECT cm.conversation_id, cm2.user_id AS peer_id, cm.unread_count AS unread, c.updated_at
		FROM conversation_members cm
		INNER JOIN conversation_members cm2 ON cm.conversation_id = cm2.conversation_id AND cm2.user_id != cm.user_id
		INNER JOIN conversations c ON c.id = cm.conversation_id
		WHERE cm.user_id = ?
		ORDER BY c.updated_at DESC
	`, userID).Scan(&rows).Error
	if err != nil {
		return nil, err
	}

	out := make([]entity.ChatSummary, 0, len(rows))
	for _, row := range rows {
		var peer entity.User
		if err := r.db.First(&peer, row.PeerID).Error; err != nil {
			continue
		}
		lastMsg, _ := r.lastMessage(row.ConversationID)
		out = append(out, entity.ChatSummary{
			ID:          row.ConversationID,
			PeerID:      row.PeerID,
			Name:        peer.Nickname,
			Avatar:      peer.Avatar,
			LastMessage: lastMsg,
			Time:        formatRelative(row.UpdatedAt),
			Unread:      row.Unread,
			HasCamera:   true,
		})
	}
	return out, nil
}

func (r *ConversationRepository) lastMessage(convID int64) (string, error) {
	var m entity.Message
	err := r.db.Where("conversation_id = ?", convID).Order("id DESC").First(&m).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return "", nil
	}
	return m.Text, err
}

func formatRelative(t time.Time) string {
	d := time.Since(t)
	switch {
	case d < time.Minute:
		return "刚刚"
	case d < 24*time.Hour:
		return fmt.Sprintf("%d小时", int(d.Hours()))
	case d < 48*time.Hour:
		return "1天"
	default:
		return fmt.Sprintf("%d天", int(d.Hours()/24))
	}
}

func (r *ConversationRepository) FindDirect(userID, peerID int64) (int64, error) {
	var convID int64
	err := r.db.Raw(`
		SELECT cm1.conversation_id FROM conversation_members cm1
		INNER JOIN conversation_members cm2 ON cm1.conversation_id = cm2.conversation_id
		WHERE cm1.user_id = ? AND cm2.user_id = ?
		LIMIT 1
	`, userID, peerID).Scan(&convID).Error
	if convID == 0 {
		return 0, gorm.ErrRecordNotFound
	}
	return convID, err
}

func (r *ConversationRepository) CreateDirect(userID, peerID int64) (int64, error) {
	var convID int64
	err := r.db.Transaction(func(tx *gorm.DB) error {
		conv := entity.Conversation{}
		if err := tx.Create(&conv).Error; err != nil {
			return err
		}
		convID = conv.ID
		members := []entity.ConversationMember{
			{ConversationID: convID, UserID: userID},
			{ConversationID: convID, UserID: peerID},
		}
		return tx.Create(&members).Error
	})
	return convID, err
}

func (r *ConversationRepository) Touch(convID int64) error {
	return r.db.Model(&entity.Conversation{}).Where("id = ?", convID).Update("updated_at", time.Now()).Error
}

func (r *ConversationRepository) MemberOf(convID, userID int64) bool {
	var count int64
	r.db.Model(&entity.ConversationMember{}).
		Where("conversation_id = ? AND user_id = ?", convID, userID).
		Count(&count)
	return count > 0
}
