package repository

import (
	"errors"

	"gorm.io/gorm"

	"github.com/LoveMuZiLi/tiktok-backend/internal/entity"
)

var ErrUserNotFound = errors.New("user not found")

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) List(offset, limit int) ([]entity.User, int64, error) {
	if limit <= 0 {
		limit = 20
	}
	if offset < 0 {
		offset = 0
	}
	var total int64
	if err := r.db.Model(&entity.User{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}
	var users []entity.User
	err := r.db.Order("id ASC").Offset(offset).Limit(limit).Find(&users).Error
	return users, total, err
}

func (r *UserRepository) GetByID(id int64) (entity.User, error) {
	var u entity.User
	err := r.db.First(&u, id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return entity.User{}, ErrUserNotFound
	}
	return u, err
}

func (r *UserRepository) Create(u *entity.User) error {
	return r.db.Create(u).Error
}

func (r *UserRepository) Update(u *entity.User) error {
	res := r.db.Save(u)
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return ErrUserNotFound
	}
	return nil
}

func (r *UserRepository) Delete(id int64) error {
	res := r.db.Delete(&entity.User{}, id)
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return ErrUserNotFound
	}
	return nil
}

func (r *UserRepository) Profile(id int64) (entity.UserProfile, error) {
	u, err := r.GetByID(id)
	if err != nil {
		return entity.UserProfile{}, err
	}
	var following, followers, videoCount int64
	_ = r.db.Model(&entity.Follow{}).Where("follower_id = ?", id).Count(&following).Error
	_ = r.db.Model(&entity.Follow{}).Where("following_id = ?", id).Count(&followers).Error
	_ = r.db.Model(&entity.Video{}).Where("user_id = ?", id).Count(&videoCount).Error
	var totalLikes int64
	_ = r.db.Model(&entity.Video{}).Where("user_id = ?", id).Select("COALESCE(SUM(likes),0)").Scan(&totalLikes).Error
	return entity.UserProfile{
		User:           u,
		FollowingCount: following,
		FollowerCount:  followers,
		TotalLikes:     totalLikes,
		VideoCount:     videoCount,
	}, nil
}
