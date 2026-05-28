package repository

import (
	"errors"

	"gorm.io/gorm"

	"github.com/LoveMuZiLi/tiktok-backend/internal/entity"
)

var ErrFollowNotFound = errors.New("follow not found")

type FollowRepository struct {
	db *gorm.DB
}

func NewFollowRepository(db *gorm.DB) *FollowRepository {
	return &FollowRepository{db: db}
}

func (r *FollowRepository) Create(followerID, followingID int64) error {
	f := entity.Follow{FollowerID: followerID, FollowingID: followingID}
	return r.db.Create(&f).Error
}

func (r *FollowRepository) Delete(followerID, followingID int64) error {
	res := r.db.Where("follower_id = ? AND following_id = ?", followerID, followingID).Delete(&entity.Follow{})
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return ErrFollowNotFound
	}
	return nil
}

func (r *FollowRepository) FollowingIDs(userID int64) ([]int64, error) {
	var ids []int64
	err := r.db.Model(&entity.Follow{}).Where("follower_id = ?", userID).Pluck("following_id", &ids).Error
	return ids, err
}

func (r *FollowRepository) FriendIDs(userID int64) ([]int64, error) {
	var ids []int64
	err := r.db.Raw(`
		SELECT f1.following_id FROM follows f1
		INNER JOIN follows f2 ON f1.following_id = f2.follower_id AND f1.follower_id = f2.following_id
		WHERE f1.follower_id = ?
	`, userID).Scan(&ids).Error
	return ids, err
}

func (r *FollowRepository) IsFollowing(followerID, followingID int64) (bool, error) {
	var count int64
	err := r.db.Model(&entity.Follow{}).
		Where("follower_id = ? AND following_id = ?", followerID, followingID).
		Count(&count).Error
	return count > 0, err
}
