package service

import (
	"errors"

	"gorm.io/gorm"

	"github.com/LoveMuZiLi/tiktok-backend/internal/entity"
	"github.com/LoveMuZiLi/tiktok-backend/internal/repository"
)

var ErrNotificationNotFound = errors.New("notification not found")

type NotificationService struct {
	repo *repository.NotificationRepository
}

func NewNotificationService(db *gorm.DB) *NotificationService {
	return &NotificationService{repo: repository.NewNotificationRepository(db)}
}

func (s *NotificationService) List(userID int64) ([]entity.Notification, error) {
	return s.repo.List(userID)
}

func (s *NotificationService) Create(n *entity.Notification) error {
	return s.repo.Create(n)
}

func (s *NotificationService) MarkRead(id int64) error {
	if err := s.repo.MarkRead(id); errors.Is(err, repository.ErrNotificationNotFound) {
		return ErrNotificationNotFound
	}
	return nil
}

func (s *NotificationService) Delete(id int64) error {
	if err := s.repo.Delete(id); errors.Is(err, repository.ErrNotificationNotFound) {
		return ErrNotificationNotFound
	}
	return nil
}
