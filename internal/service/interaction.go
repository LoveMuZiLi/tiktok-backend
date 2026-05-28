package service

import (
	"errors"

	"gorm.io/gorm"

	"github.com/LoveMuZiLi/tiktok-backend/internal/entity"
	"github.com/LoveMuZiLi/tiktok-backend/internal/repository"
)

var ErrCommentNotFound = errors.New("comment not found")

type InteractionService struct {
	interact *repository.InteractionRepository
	videos   *repository.VideoRepository
}

func NewInteractionService(db *gorm.DB) *InteractionService {
	return &InteractionService{
		interact: repository.NewInteractionRepository(db),
		videos:   repository.NewVideoRepository(db),
	}
}

func (s *InteractionService) Status(userID, videoID int64) (entity.InteractionStatus, error) {
	liked, err := s.interact.IsLiked(userID, videoID)
	if err != nil {
		return entity.InteractionStatus{}, err
	}
	favorited, err := s.interact.IsFavorited(userID, videoID)
	if err != nil {
		return entity.InteractionStatus{}, err
	}
	return entity.InteractionStatus{Liked: liked, Favorited: favorited}, nil
}

func (s *InteractionService) Like(userID, videoID int64) (entity.VideoView, error) {
	if _, err := s.videos.GetByID(videoID); err != nil {
		if errors.Is(err, repository.ErrVideoNotFound) {
			return entity.VideoView{}, ErrVideoNotFound
		}
		return entity.VideoView{}, err
	}
	if err := s.interact.AddLike(userID, videoID); err != nil && !errors.Is(err, repository.ErrAlreadyLiked) {
		return entity.VideoView{}, err
	}
	return s.videoView(videoID)
}

func (s *InteractionService) Unlike(userID, videoID int64) (entity.VideoView, error) {
	if err := s.interact.RemoveLike(userID, videoID); err != nil {
		if errors.Is(err, repository.ErrNotLiked) {
			return s.videoView(videoID)
		}
		return entity.VideoView{}, err
	}
	return s.videoView(videoID)
}

func (s *InteractionService) Favorite(userID, videoID int64) (entity.VideoView, error) {
	if _, err := s.videos.GetByID(videoID); err != nil {
		return entity.VideoView{}, ErrVideoNotFound
	}
	if err := s.interact.AddFavorite(userID, videoID); err != nil && !errors.Is(err, repository.ErrAlreadyFavorited) {
		return entity.VideoView{}, err
	}
	return s.videoView(videoID)
}

func (s *InteractionService) Unfavorite(userID, videoID int64) (entity.VideoView, error) {
	if err := s.interact.RemoveFavorite(userID, videoID); err != nil {
		if errors.Is(err, repository.ErrNotFavorited) {
			return s.videoView(videoID)
		}
		return entity.VideoView{}, err
	}
	return s.videoView(videoID)
}

func (s *InteractionService) Share(userID, videoID int64) (entity.VideoView, error) {
	if _, err := s.videos.GetByID(videoID); err != nil {
		return entity.VideoView{}, ErrVideoNotFound
	}
	if err := s.interact.AddShare(userID, videoID); err != nil {
		return entity.VideoView{}, err
	}
	return s.videoView(videoID)
}

func (s *InteractionService) ListComments(videoID int64, offset, limit int) ([]entity.CommentView, int64, error) {
	return s.interact.ListComments(videoID, offset, limit)
}

func (s *InteractionService) CreateComment(userID, videoID int64, content string) (entity.CommentView, error) {
	if _, err := s.videos.GetByID(videoID); err != nil {
		return entity.CommentView{}, ErrVideoNotFound
	}
	return s.interact.CreateComment(userID, videoID, content)
}

func (s *InteractionService) UpdateComment(commentID, userID int64, content string) (entity.CommentView, error) {
	v, err := s.interact.UpdateComment(commentID, userID, content)
	if errors.Is(err, repository.ErrCommentNotFound) {
		return entity.CommentView{}, ErrCommentNotFound
	}
	return v, err
}

func (s *InteractionService) DeleteComment(commentID, userID int64) error {
	if err := s.interact.DeleteComment(commentID, userID); err != nil {
		if errors.Is(err, repository.ErrCommentNotFound) {
			return ErrCommentNotFound
		}
		return err
	}
	return nil
}

func (s *InteractionService) videoView(videoID int64) (entity.VideoView, error) {
	v, err := s.videos.GetByID(videoID)
	if err != nil {
		return entity.VideoView{}, err
	}
	return entity.ToVideoView(v), nil
}
