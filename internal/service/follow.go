package service

import (
	"errors"

	"gorm.io/gorm"

	"github.com/LoveMuZiLi/tiktok-backend/internal/repository"
)

var ErrFollowNotFound = errors.New("follow not found")

type FollowService struct {
	repo *repository.FollowRepository
}

func NewFollowService(db *gorm.DB) *FollowService {
	return &FollowService{repo: repository.NewFollowRepository(db)}
}

func (s *FollowService) Follow(followerID, followingID int64) error {
	if followerID == followingID {
		return errors.New("cannot follow yourself")
	}
	return s.repo.Create(followerID, followingID)
}

func (s *FollowService) Unfollow(followerID, followingID int64) error {
	if err := s.repo.Delete(followerID, followingID); errors.Is(err, repository.ErrFollowNotFound) {
		return ErrFollowNotFound
	}
	return nil
}

func (s *FollowService) IsFollowing(followerID, followingID int64) (bool, error) {
	return s.repo.IsFollowing(followerID, followingID)
}
