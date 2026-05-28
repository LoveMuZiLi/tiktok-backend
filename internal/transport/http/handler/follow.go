package handler

import (
	"errors"
	"net/http"

	"github.com/LoveMuZiLi/tiktok-backend/internal/service"
)

type FollowHandler struct {
	follows *service.FollowService
}

func NewFollowHandler(follows *service.FollowService) *FollowHandler {
	return &FollowHandler{follows: follows}
}

type followBody struct {
	FollowerID  int64 `json:"followerId"`
	FollowingID int64 `json:"followingId"`
}

func (h *FollowHandler) Create(w http.ResponseWriter, r *http.Request) {
	var body followBody
	if err := decodeJSON(r, &body); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid body"})
		return
	}
	if body.FollowerID <= 0 {
		body.FollowerID = parseUserID(r)
	}
	if err := h.follows.Follow(body.FollowerID, body.FollowingID); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}
	writeJSON(w, http.StatusCreated, map[string]string{"status": "ok"})
}

func (h *FollowHandler) Delete(w http.ResponseWriter, r *http.Request) {
	followingID, ok := parseIDParam(r, "followingId")
	if !ok {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid following id"})
		return
	}
	followerID := parseUserID(r)
	if err := h.follows.Unfollow(followerID, followingID); err != nil {
		if errors.Is(err, service.ErrFollowNotFound) {
			writeJSON(w, http.StatusNotFound, map[string]string{"error": "not found"})
			return
		}
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "internal error"})
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *FollowHandler) Status(w http.ResponseWriter, r *http.Request) {
	followingID, ok := parseIDParam(r, "followingId")
	if !ok {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid following id"})
		return
	}
	followerID := parseUserID(r)
	isFollowing, err := h.follows.IsFollowing(followerID, followingID)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "internal error"})
		return
	}
	writeJSON(w, http.StatusOK, map[string]bool{"isFollowing": isFollowing})
}
