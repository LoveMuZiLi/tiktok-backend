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
	err := r.db.Order("id ASC").Offset(offset).Limit(limit).Find(&items).Error
	if err != nil {
		return entity.VideoFeedResponse{}, err
	}

	return entity.VideoFeedResponse{
		Items:  items,
		Total:  total,
		Offset: offset,
		Limit:  limit,
	}, nil
}

func (r *VideoRepository) GetByID(id int64) (entity.Video, error) {
	var video entity.Video
	err := r.db.First(&video, id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return entity.Video{}, ErrVideoNotFound
	}
	return video, err
}

func (r *VideoRepository) IncrementLikes(id int64) (entity.Video, error) {
	res := r.db.Model(&entity.Video{}).Where("id = ?", id).UpdateColumn("likes", gorm.Expr("likes + ?", 1))
	if res.Error != nil {
		return entity.Video{}, res.Error
	}
	if res.RowsAffected == 0 {
		return entity.Video{}, ErrVideoNotFound
	}
	return r.GetByID(id)
}
