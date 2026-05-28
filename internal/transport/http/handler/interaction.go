package handler

import (
	"errors"
	"net/http"

	"github.com/LoveMuZiLi/tiktok-backend/internal/entity"
	"github.com/LoveMuZiLi/tiktok-backend/internal/service"
)

type InteractionHandler struct {
	interactions *service.InteractionService
}

func NewInteractionHandler(interactions *service.InteractionService) *InteractionHandler {
	return &InteractionHandler{interactions: interactions}
}

func (h *InteractionHandler) Status(w http.ResponseWriter, r *http.Request) {
	videoID, ok := parseIDParam(r, "id")
	if !ok {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid id"})
		return
	}
	status, err := h.interactions.Status(parseUserID(r), videoID)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "internal error"})
		return
	}
	writeJSON(w, http.StatusOK, status)
}

func (h *InteractionHandler) Like(w http.ResponseWriter, r *http.Request) {
	videoID, ok := parseIDParam(r, "id")
	if !ok {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid id"})
		return
	}
	video, err := h.interactions.Like(parseUserID(r), videoID)
	if err != nil {
		if errors.Is(err, service.ErrVideoNotFound) {
			writeJSON(w, http.StatusNotFound, map[string]string{"error": "not found"})
			return
		}
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "internal error"})
		return
	}
	writeJSON(w, http.StatusOK, video)
}

func (h *InteractionHandler) Unlike(w http.ResponseWriter, r *http.Request) {
	videoID, ok := parseIDParam(r, "id")
	if !ok {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid id"})
		return
	}
	video, err := h.interactions.Unlike(parseUserID(r), videoID)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "internal error"})
		return
	}
	writeJSON(w, http.StatusOK, video)
}

func (h *InteractionHandler) Favorite(w http.ResponseWriter, r *http.Request) {
	videoID, ok := parseIDParam(r, "id")
	if !ok {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid id"})
		return
	}
	video, err := h.interactions.Favorite(parseUserID(r), videoID)
	if err != nil {
		if errors.Is(err, service.ErrVideoNotFound) {
			writeJSON(w, http.StatusNotFound, map[string]string{"error": "not found"})
			return
		}
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "internal error"})
		return
	}
	writeJSON(w, http.StatusOK, video)
}

func (h *InteractionHandler) Unfavorite(w http.ResponseWriter, r *http.Request) {
	videoID, ok := parseIDParam(r, "id")
	if !ok {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid id"})
		return
	}
	video, err := h.interactions.Unfavorite(parseUserID(r), videoID)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "internal error"})
		return
	}
	writeJSON(w, http.StatusOK, video)
}

func (h *InteractionHandler) Share(w http.ResponseWriter, r *http.Request) {
	videoID, ok := parseIDParam(r, "id")
	if !ok {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid id"})
		return
	}
	video, err := h.interactions.Share(parseUserID(r), videoID)
	if err != nil {
		if errors.Is(err, service.ErrVideoNotFound) {
			writeJSON(w, http.StatusNotFound, map[string]string{"error": "not found"})
			return
		}
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "internal error"})
		return
	}
	writeJSON(w, http.StatusOK, video)
}

func (h *InteractionHandler) ListComments(w http.ResponseWriter, r *http.Request) {
	videoID, ok := parseIDParam(r, "id")
	if !ok {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid id"})
		return
	}
	offset := parseQueryInt(r, "offset", 0)
	limit := parseQueryInt(r, "limit", 20)
	items, total, err := h.interactions.ListComments(videoID, offset, limit)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "failed to list comments"})
		return
	}
	if items == nil {
		items = []entity.CommentView{}
	}
	writeJSON(w, http.StatusOK, map[string]any{"items": items, "total": total, "offset": offset, "limit": limit})
}

func (h *InteractionHandler) CreateComment(w http.ResponseWriter, r *http.Request) {
	videoID, ok := parseIDParam(r, "id")
	if !ok {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid id"})
		return
	}
	var body entity.CreateCommentInput
	if err := decodeJSON(r, &body); err != nil || body.Content == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "content required"})
		return
	}
	userID := body.UserID
	if userID <= 0 {
		userID = parseUserID(r)
	}
	comment, err := h.interactions.CreateComment(userID, videoID, body.Content)
	if err != nil {
		if errors.Is(err, service.ErrVideoNotFound) {
			writeJSON(w, http.StatusNotFound, map[string]string{"error": "not found"})
			return
		}
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "failed to create comment"})
		return
	}
	writeJSON(w, http.StatusCreated, comment)
}

func (h *InteractionHandler) UpdateComment(w http.ResponseWriter, r *http.Request) {
	commentID, ok := parseIDParam(r, "commentId")
	if !ok {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid comment id"})
		return
	}
	var body entity.CreateCommentInput
	if err := decodeJSON(r, &body); err != nil || body.Content == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "content required"})
		return
	}
	comment, err := h.interactions.UpdateComment(commentID, parseUserID(r), body.Content)
	if err != nil {
		if errors.Is(err, service.ErrCommentNotFound) {
			writeJSON(w, http.StatusNotFound, map[string]string{"error": "not found"})
			return
		}
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "internal error"})
		return
	}
	writeJSON(w, http.StatusOK, comment)
}

func (h *InteractionHandler) DeleteComment(w http.ResponseWriter, r *http.Request) {
	commentID, ok := parseIDParam(r, "commentId")
	if !ok {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid comment id"})
		return
	}
	if err := h.interactions.DeleteComment(commentID, parseUserID(r)); err != nil {
		if errors.Is(err, service.ErrCommentNotFound) {
			writeJSON(w, http.StatusNotFound, map[string]string{"error": "not found"})
			return
		}
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "internal error"})
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
