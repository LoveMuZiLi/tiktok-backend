package repository

import (
	"errors"

	"gorm.io/gorm"

	"github.com/LoveMuZiLi/tiktok-backend/internal/entity"
)

var ErrVideoNotFound = errors.New("video not found")

type VideoRepository struct {
	db *gorm.DB
}

func NewVideoRepository(db *gorm.DB) *VideoRepository {
	return &VideoRepository{db: db}
}

func (r *VideoRepository) preload() *gorm.DB {
	return r.db.Preload("User")
}

func (r *VideoRepository) List(offset, limit int) (entity.VideoFeedResponse, error) {
	if offset < 0 {
		offset = 0
	}
	if limit <= 0 {
		limit = 20
	}
	var total int64
	if err := r.db.Model(&entity.Video{}).Count(&total).Error; err != nil {
		return entity.VideoFeedResponse{}, err
	}
	var items []entity.Video
	err := r.preload().Order("id DESC").Offset(offset).Limit(limit).Find(&items).Error
	if err != nil {
		return entity.VideoFeedResponse{}, err
	}
	return entity.VideoFeedResponse{
		Items:  entity.ToVideoViews(items),
		Total:  total,
		Offset: offset,
		Limit:  limit,
	}, nil
}

func (r *VideoRepository) ListByUserIDs(userIDs []int64, offset, limit int) (entity.VideoFeedResponse, error) {
	if len(userIDs) == 0 {
		return entity.VideoFeedResponse{Items: []entity.VideoView{}, Offset: offset, Limit: limit}, nil
	}
	if limit <= 0 {
		limit = 20
	}
	var total int64
	q := r.db.Model(&entity.Video{}).Where("user_id IN ?", userIDs)
	if err := q.Count(&total).Error; err != nil {
		return entity.VideoFeedResponse{}, err
	}
	var items []entity.Video
	err := r.preload().Where("user_id IN ?", userIDs).Order("id DESC").Offset(offset).Limit(limit).Find(&items).Error
	if err != nil {
		return entity.VideoFeedResponse{}, err
	}
	return entity.VideoFeedResponse{
		Items:  entity.ToVideoViews(items),
		Total:  total,
		Offset: offset,
		Limit:  limit,
	}, nil
}

func (r *VideoRepository) ListByUserID(userID int64, offset, limit int) (entity.VideoFeedResponse, error) {
	return r.ListByUserIDs([]int64{userID}, offset, limit)
}

func (r *VideoRepository) GetByID(id int64) (entity.Video, error) {
	var video entity.Video
	err := r.preload().First(&video, id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return entity.Video{}, ErrVideoNotFound
	}
	return video, err
}

func (r *VideoRepository) Create(v *entity.Video) error {
	return r.db.Create(v).Error
}

func (r *VideoRepository) Update(v *entity.Video) error {
	res := r.db.Save(v)
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return ErrVideoNotFound
	}
	return nil
}

func (r *VideoRepository) Delete(id int64) error {
	res := r.db.Delete(&entity.Video{}, id)
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return ErrVideoNotFound
	}
	return nil
}

