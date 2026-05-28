package service

import (
	"errors"

	"gorm.io/gorm"

	"github.com/LoveMuZiLi/tiktok-backend/internal/entity"
	"github.com/LoveMuZiLi/tiktok-backend/internal/repository"
)

var ErrUserNotFound = errors.New("user not found")

type UserService struct {
	repo *repository.UserRepository
}

func NewUserService(db *gorm.DB) *UserService {
	return &UserService{repo: repository.NewUserRepository(db)}
}

func (s *UserService) List(offset, limit int) ([]entity.User, int64, error) {
	return s.repo.List(offset, limit)
}

func (s *UserService) Get(id int64) (entity.User, error) {
	u, err := s.repo.GetByID(id)
	if errors.Is(err, repository.ErrUserNotFound) {
		return entity.User{}, ErrUserNotFound
	}
	return u, err
}

func (s *UserService) Profile(id int64) (entity.UserProfile, error) {
	p, err := s.repo.Profile(id)
	if errors.Is(err, repository.ErrUserNotFound) {
		return entity.UserProfile{}, ErrUserNotFound
	}
	return p, err
}

func (s *UserService) Create(u *entity.User) error {
	return s.repo.Create(u)
}

func (s *UserService) Update(u *entity.User) error {
	if err := s.repo.Update(u); errors.Is(err, repository.ErrUserNotFound) {
		return ErrUserNotFound
	}
	return nil
}

func (s *UserService) Delete(id int64) error {
	if err := s.repo.Delete(id); errors.Is(err, repository.ErrUserNotFound) {
		return ErrUserNotFound
	}
	return nil
}
