package service

import (
	"errors"

	"gorm.io/gorm"

	"github.com/LoveMuZiLi/tiktok-backend/internal/entity"
	"github.com/LoveMuZiLi/tiktok-backend/internal/repository"
)

type InboxService struct {
	conv *repository.ConversationRepository
	msg  *repository.MessageRepository
}

func NewInboxService(db *gorm.DB) *InboxService {
	return &InboxService{
		conv: repository.NewConversationRepository(db),
		msg:  repository.NewMessageRepository(db),
	}
}

func (s *InboxService) ListChats(userID int64) ([]entity.ChatSummary, error) {
	return s.conv.ListChats(userID)
}

func (s *InboxService) ListMessages(convID, userID int64) ([]entity.MessageView, error) {
	if !s.conv.MemberOf(convID, userID) {
		return nil, errors.New("forbidden")
	}
	_ = s.msg.ClearUnread(convID, userID)
	return s.msg.List(convID, userID)
}

func (s *InboxService) SendMessage(convID, senderID int64, text string) (entity.MessageView, error) {
	if !s.conv.MemberOf(convID, senderID) {
		return entity.MessageView{}, errors.New("forbidden")
	}
	m, err := s.msg.Create(convID, senderID, text)
	if err != nil {
		return entity.MessageView{}, err
	}
	_ = s.conv.Touch(convID)
	_ = s.msg.IncrementUnread(convID, senderID)
	return entity.MessageView{
		ID:     m.ID,
		Text:   m.Text,
		IsMine: true,
		Time:   m.CreatedAt.Format("15:04"),
	}, nil
}

func (s *InboxService) OpenOrCreateChat(userID, peerID int64) (int64, error) {
	convID, err := s.conv.FindDirect(userID, peerID)
	if err == nil {
		return convID, nil
	}
	return s.conv.CreateDirect(userID, peerID)
}

func (s *InboxService) DeleteMessage(id int64) error {
	return s.msg.Delete(id)
}
