package service

import (
	"errors"

	"gorm.io/gorm"

	"github.com/LoveMuZiLi/tiktok-backend/internal/entity"
	"github.com/LoveMuZiLi/tiktok-backend/internal/repository"
)

var ErrVideoNotFound = errors.New("video not found")

type VideoService struct {
	repo *repository.VideoRepository
}

func NewVideoService(db *gorm.DB) *VideoService {
	return &VideoService{repo: repository.NewVideoRepository(db)}
}

func (s *VideoService) List(offset, limit int) (entity.VideoFeedResponse, error) {
	return s.repo.List(offset, limit)
}

func (s *VideoService) Get(id int64) (entity.Video, error) {
	v, err := s.repo.GetByID(id)
	if errors.Is(err, repository.ErrVideoNotFound) {
		return entity.Video{}, ErrVideoNotFound
	}
	return v, err
}

func (s *VideoService) Like(id int64) (entity.Video, error) {
	v, err := s.repo.IncrementLikes(id)
	if errors.Is(err, repository.ErrVideoNotFound) {
		return entity.Video{}, ErrVideoNotFound
	}
	return v, err
}
