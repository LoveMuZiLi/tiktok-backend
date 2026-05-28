package service

import (
	"errors"

	"gorm.io/gorm"

	"tiktok-app/backend/internal/model"
	"tiktok-app/backend/internal/repository"
)

var ErrVideoNotFound = errors.New("video not found")

type VideoService struct {
	repo *repository.VideoRepository
}

func NewVideoService(db *gorm.DB) *VideoService {
	return &VideoService{repo: repository.NewVideoRepository(db)}
}

func (s *VideoService) List(offset, limit int) (model.VideoFeedResponse, error) {
	return s.repo.List(offset, limit)
}

func (s *VideoService) Get(id int64) (model.Video, error) {
	v, err := s.repo.GetByID(id)
	if errors.Is(err, repository.ErrVideoNotFound) {
		return model.Video{}, ErrVideoNotFound
	}
	return v, err
}

func (s *VideoService) Like(id int64) (model.Video, error) {
	v, err := s.repo.IncrementLikes(id)
	if errors.Is(err, repository.ErrVideoNotFound) {
		return model.Video{}, ErrVideoNotFound
	}
	return v, err
}
