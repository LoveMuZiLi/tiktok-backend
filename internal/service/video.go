package service

import (
	"errors"

	"gorm.io/gorm"

	"github.com/LoveMuZiLi/tiktok-backend/internal/entity"
	"github.com/LoveMuZiLi/tiktok-backend/internal/repository"
)

var ErrVideoNotFound = errors.New("video not found")

type VideoService struct {
	repo   *repository.VideoRepository
	follow *repository.FollowRepository
}

func NewVideoService(db *gorm.DB) *VideoService {
	return &VideoService{
		repo:   repository.NewVideoRepository(db),
		follow: repository.NewFollowRepository(db),
	}
}

func (s *VideoService) List(offset, limit int) (entity.VideoFeedResponse, error) {
	return s.repo.List(offset, limit)
}

func (s *VideoService) ListFollowing(userID int64, offset, limit int) (entity.VideoFeedResponse, error) {
	ids, err := s.follow.FollowingIDs(userID)
	if err != nil {
		return entity.VideoFeedResponse{}, err
	}
	return s.repo.ListByUserIDs(ids, offset, limit)
}

func (s *VideoService) ListFriends(userID int64, offset, limit int) (entity.VideoFeedResponse, error) {
	ids, err := s.follow.FriendIDs(userID)
	if err != nil {
		return entity.VideoFeedResponse{}, err
	}
	return s.repo.ListByUserIDs(ids, offset, limit)
}

func (s *VideoService) ListByUser(userID int64, offset, limit int) (entity.VideoFeedResponse, error) {
	return s.repo.ListByUserID(userID, offset, limit)
}

func (s *VideoService) Get(id int64) (entity.VideoView, error) {
	v, err := s.repo.GetByID(id)
	if errors.Is(err, repository.ErrVideoNotFound) {
		return entity.VideoView{}, ErrVideoNotFound
	}
	return entity.ToVideoView(v), err
}

func (s *VideoService) Create(in entity.CreateVideoInput) (entity.VideoView, error) {
	v := entity.Video{
		UserID:      in.UserID,
		Description: in.Description,
		Music:       in.Music,
		Image:       in.Image,
	}
	if err := s.repo.Create(&v); err != nil {
		return entity.VideoView{}, err
	}
	full, err := s.repo.GetByID(v.ID)
	if err != nil {
		return entity.VideoView{}, err
	}
	return entity.ToVideoView(full), nil
}

func (s *VideoService) Update(id int64, in entity.UpdateVideoInput) (entity.VideoView, error) {
	v, err := s.repo.GetByID(id)
	if errors.Is(err, repository.ErrVideoNotFound) {
		return entity.VideoView{}, ErrVideoNotFound
	}
	if in.Description != nil {
		v.Description = *in.Description
	}
	if in.Music != nil {
		v.Music = *in.Music
	}
	if in.Image != nil {
		v.Image = *in.Image
	}
	if err := s.repo.Update(&v); err != nil {
		return entity.VideoView{}, err
	}
	full, err := s.repo.GetByID(id)
	if err != nil {
		return entity.VideoView{}, err
	}
	return entity.ToVideoView(full), nil
}

func (s *VideoService) Delete(id int64) error {
	if err := s.repo.Delete(id); errors.Is(err, repository.ErrVideoNotFound) {
		return ErrVideoNotFound
	}
	return nil
}
