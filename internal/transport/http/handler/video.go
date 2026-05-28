package handler

import (
	"errors"
	"net/http"

	"github.com/LoveMuZiLi/tiktok-backend/internal/entity"
	"github.com/LoveMuZiLi/tiktok-backend/internal/service"
)

type VideoHandler struct {
	videos *service.VideoService
}

func NewVideoHandler(videos *service.VideoService) *VideoHandler {
	return &VideoHandler{videos: videos}
}

func (h *VideoHandler) List(w http.ResponseWriter, r *http.Request) {
	offset := parseQueryInt(r, "offset", 0)
	limit := parseQueryInt(r, "limit", 20)
	feedType := r.URL.Query().Get("feed")
	userID := parseUserID(r)

	var feed entity.VideoFeedResponse
	var err error
	switch feedType {
	case "following":
		feed, err = h.videos.ListFollowing(userID, offset, limit)
	case "friends":
		feed, err = h.videos.ListFriends(userID, offset, limit)
	case "user":
		targetID := parseQueryInt(r, "target_id", int(userID))
		feed, err = h.videos.ListByUser(int64(targetID), offset, limit)
	default:
		feed, err = h.videos.List(offset, limit)
	}
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "failed to list videos"})
		return
	}
	writeJSON(w, http.StatusOK, feed)
}

func (h *VideoHandler) Get(w http.ResponseWriter, r *http.Request) {
	id, ok := parseIDParam(r, "id")
	if !ok {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid id"})
		return
	}
	video, err := h.videos.Get(id)
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

func (h *VideoHandler) Create(w http.ResponseWriter, r *http.Request) {
	var body entity.CreateVideoInput
	if err := decodeJSON(r, &body); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid body"})
		return
	}
	if body.UserID <= 0 {
		body.UserID = parseUserID(r)
	}
	if body.Description == "" || body.Image == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "description and image required"})
		return
	}
	video, err := h.videos.Create(body)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "failed to create video"})
		return
	}
	writeJSON(w, http.StatusCreated, video)
}

func (h *VideoHandler) Update(w http.ResponseWriter, r *http.Request) {
	id, ok := parseIDParam(r, "id")
	if !ok {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid id"})
		return
	}
	var body entity.UpdateVideoInput
	if err := decodeJSON(r, &body); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid body"})
		return
	}
	video, err := h.videos.Update(id, body)
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

func (h *VideoHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id, ok := parseIDParam(r, "id")
	if !ok {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid id"})
		return
	}
	if err := h.videos.Delete(id); err != nil {
		if errors.Is(err, service.ErrVideoNotFound) {
			writeJSON(w, http.StatusNotFound, map[string]string{"error": "not found"})
			return
		}
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "internal error"})
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
