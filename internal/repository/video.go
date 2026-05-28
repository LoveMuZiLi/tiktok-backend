package repository

import (
	"errors"

	"gorm.io/gorm"

	"tiktok-app/backend/internal/model"
)

var ErrVideoNotFound = errors.New("video not found")

type VideoRepository struct {
	db *gorm.DB
}

func NewVideoRepository(db *gorm.DB) *VideoRepository {
	return &VideoRepository{db: db}
}

func (r *VideoRepository) List(offset, limit int) (model.VideoFeedResponse, error) {
	if offset < 0 {
		offset = 0
	}
	if limit <= 0 {
		limit = 20
	}

	var total int64
	if err := r.db.Model(&model.Video{}).Count(&total).Error; err != nil {
		return model.VideoFeedResponse{}, err
	}

	var items []model.Video
	err := r.db.Order("id ASC").Offset(offset).Limit(limit).Find(&items).Error
	if err != nil {
		return model.VideoFeedResponse{}, err
	}

	return model.VideoFeedResponse{
		Items:  items,
		Total:  total,
		Offset: offset,
		Limit:  limit,
	}, nil
}

func (r *VideoRepository) GetByID(id int64) (model.Video, error) {
	var video model.Video
	err := r.db.First(&video, id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return model.Video{}, ErrVideoNotFound
	}
	return video, err
}

func (r *VideoRepository) IncrementLikes(id int64) (model.Video, error) {
	res := r.db.Model(&model.Video{}).Where("id = ?", id).UpdateColumn("likes", gorm.Expr("likes + ?", 1))
	if res.Error != nil {
		return model.Video{}, res.Error
	}
	if res.RowsAffected == 0 {
		return model.Video{}, ErrVideoNotFound
	}
	return r.GetByID(id)
}
