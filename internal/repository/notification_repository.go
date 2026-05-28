package repository

import (
	"errors"

	"gorm.io/gorm"

	"github.com/LoveMuZiLi/tiktok-backend/internal/entity"
)

var ErrNotificationNotFound = errors.New("notification not found")

type NotificationRepository struct {
	db *gorm.DB
}

func NewNotificationRepository(db *gorm.DB) *NotificationRepository {
	return &NotificationRepository{db: db}
}

func (r *NotificationRepository) List(userID int64) ([]entity.Notification, error) {
	var items []entity.Notification
	err := r.db.Where("user_id = ?", userID).Order("id DESC").Find(&items).Error
	return items, err
}

func (r *NotificationRepository) GetByID(id int64) (entity.Notification, error) {
	var n entity.Notification
	err := r.db.First(&n, id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return entity.Notification{}, ErrNotificationNotFound
	}
	return n, err
}

func (r *NotificationRepository) Create(n *entity.Notification) error {
	return r.db.Create(n).Error
}

func (r *NotificationRepository) MarkRead(id int64) error {
	res := r.db.Model(&entity.Notification{}).Where("id = ?", id).Update("is_read", true)
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return ErrNotificationNotFound
	}
	return nil
}

func (r *NotificationRepository) Delete(id int64) error {
	res := r.db.Delete(&entity.Notification{}, id)
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return ErrNotificationNotFound
	}
	return nil
}
