package repository

import (
	"errors"

	"gorm.io/gorm"

	"github.com/LoveMuZiLi/tiktok-backend/internal/entity"
)

var (
	ErrAlreadyLiked     = errors.New("already liked")
	ErrNotLiked         = errors.New("not liked")
	ErrAlreadyFavorited = errors.New("already favorited")
	ErrNotFavorited     = errors.New("not favorited")
	ErrCommentNotFound  = errors.New("comment not found")
)

type InteractionRepository struct {
	db *gorm.DB
}

func NewInteractionRepository(db *gorm.DB) *InteractionRepository {
	return &InteractionRepository{db: db}
}

func (r *InteractionRepository) SyncVideoCounts(videoID int64) error {
	var likes, favorites, comments, shares int64
	_ = r.db.Model(&entity.VideoLike{}).Where("video_id = ?", videoID).Count(&likes).Error
	_ = r.db.Model(&entity.VideoFavorite{}).Where("video_id = ?", videoID).Count(&favorites).Error
	_ = r.db.Model(&entity.VideoComment{}).Where("video_id = ?", videoID).Count(&comments).Error
	_ = r.db.Model(&entity.VideoShare{}).Where("video_id = ?", videoID).Count(&shares).Error
	return r.db.Model(&entity.Video{}).Where("id = ?", videoID).Updates(map[string]any{
		"likes": likes, "favorites": favorites, "comments": comments, "shares": shares,
	}).Error
}

func (r *InteractionRepository) AddLike(userID, videoID int64) error {
	var n int64
	if err := r.db.Model(&entity.VideoLike{}).Where("user_id = ? AND video_id = ?", userID, videoID).Count(&n).Error; err != nil {
		return err
	}
	if n > 0 {
		return ErrAlreadyLiked
	}
	if err := r.db.Create(&entity.VideoLike{UserID: userID, VideoID: videoID}).Error; err != nil {
		return err
	}
	return r.SyncVideoCounts(videoID)
}

func (r *InteractionRepository) RemoveLike(userID, videoID int64) error {
	res := r.db.Where("user_id = ? AND video_id = ?", userID, videoID).Delete(&entity.VideoLike{})
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return ErrNotLiked
	}
	return r.SyncVideoCounts(videoID)
}

func (r *InteractionRepository) IsLiked(userID, videoID int64) (bool, error) {
	var n int64
	err := r.db.Model(&entity.VideoLike{}).Where("user_id = ? AND video_id = ?", userID, videoID).Count(&n).Error
	return n > 0, err
}

func (r *InteractionRepository) AddFavorite(userID, videoID int64) error {
	var n int64
	if err := r.db.Model(&entity.VideoFavorite{}).Where("user_id = ? AND video_id = ?", userID, videoID).Count(&n).Error; err != nil {
		return err
	}
	if n > 0 {
		return ErrAlreadyFavorited
	}
	if err := r.db.Create(&entity.VideoFavorite{UserID: userID, VideoID: videoID}).Error; err != nil {
		return err
	}
	return r.SyncVideoCounts(videoID)
}

func (r *InteractionRepository) RemoveFavorite(userID, videoID int64) error {
	res := r.db.Where("user_id = ? AND video_id = ?", userID, videoID).Delete(&entity.VideoFavorite{})
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return ErrNotFavorited
	}
	return r.SyncVideoCounts(videoID)
}

func (r *InteractionRepository) IsFavorited(userID, videoID int64) (bool, error) {
	var n int64
	err := r.db.Model(&entity.VideoFavorite{}).Where("user_id = ? AND video_id = ?", userID, videoID).Count(&n).Error
	return n > 0, err
}

func (r *InteractionRepository) AddShare(userID, videoID int64) error {
	if err := r.db.Create(&entity.VideoShare{UserID: userID, VideoID: videoID}).Error; err != nil {
		return err
	}
	return r.SyncVideoCounts(videoID)
}

func (r *InteractionRepository) ListComments(videoID int64, offset, limit int) ([]entity.CommentView, int64, error) {
	if limit <= 0 {
		limit = 20
	}
	if offset < 0 {
		offset = 0
	}
	var total int64
	if err := r.db.Model(&entity.VideoComment{}).Where("video_id = ?", videoID).Count(&total).Error; err != nil {
		return nil, 0, err
	}
	var rows []entity.VideoComment
	err := r.db.Preload("User").Where("video_id = ?", videoID).
		Order("id DESC").Offset(offset).Limit(limit).Find(&rows).Error
	if err != nil {
		return nil, 0, err
	}
	out := make([]entity.CommentView, len(rows))
	for i, c := range rows {
		out[i] = entity.ToCommentView(c)
	}
	return out, total, nil
}

func (r *InteractionRepository) CreateComment(userID, videoID int64, content string) (entity.CommentView, error) {
	c := entity.VideoComment{UserID: userID, VideoID: videoID, Content: content}
	if err := r.db.Create(&c).Error; err != nil {
		return entity.CommentView{}, err
	}
	if err := r.SyncVideoCounts(videoID); err != nil {
		return entity.CommentView{}, err
	}
	var full entity.VideoComment
	if err := r.db.Preload("User").First(&full, c.ID).Error; err != nil {
		return entity.CommentView{}, err
	}
	return entity.ToCommentView(full), nil
}

func (r *InteractionRepository) UpdateComment(id, userID int64, content string) (entity.CommentView, error) {
	var c entity.VideoComment
	if err := r.db.Preload("User").First(&c, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return entity.CommentView{}, ErrCommentNotFound
		}
		return entity.CommentView{}, err
	}
	if c.UserID != userID {
		return entity.CommentView{}, ErrCommentNotFound
	}
	c.Content = content
	if err := r.db.Save(&c).Error; err != nil {
		return entity.CommentView{}, err
	}
	return entity.ToCommentView(c), nil
}

func (r *InteractionRepository) DeleteComment(id, userID int64) error {
	var c entity.VideoComment
	if err := r.db.First(&c, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrCommentNotFound
		}
		return err
	}
	if c.UserID != userID {
		return ErrCommentNotFound
	}
	videoID := c.VideoID
	if err := r.db.Delete(&c).Error; err != nil {
		return err
	}
	return r.SyncVideoCounts(videoID)
}

func (r *InteractionRepository) GetComment(id int64) (entity.VideoComment, error) {
	var c entity.VideoComment
	err := r.db.Preload("User").First(&c, id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return entity.VideoComment{}, ErrCommentNotFound
	}
	return c, err
}
