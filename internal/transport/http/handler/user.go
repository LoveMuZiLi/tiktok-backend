package handler

import (
	"errors"
	"net/http"

	"github.com/LoveMuZiLi/tiktok-backend/internal/entity"
	"github.com/LoveMuZiLi/tiktok-backend/internal/service"
)

type UserHandler struct {
	users *service.UserService
}

func NewUserHandler(users *service.UserService) *UserHandler {
	return &UserHandler{users: users}
}

func (h *UserHandler) List(w http.ResponseWriter, r *http.Request) {
	offset := parseQueryInt(r, "offset", 0)
	limit := parseQueryInt(r, "limit", 20)
	items, total, err := h.users.List(offset, limit)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "failed to list users"})
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"items": items, "total": total, "offset": offset, "limit": limit})
}

func (h *UserHandler) Get(w http.ResponseWriter, r *http.Request) {
	id, ok := parseIDParam(r, "id")
	if !ok {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid id"})
		return
	}
	u, err := h.users.Get(id)
	if err != nil {
		if errors.Is(err, service.ErrUserNotFound) {
			writeJSON(w, http.StatusNotFound, map[string]string{"error": "not found"})
			return
		}
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "internal error"})
		return
	}
	writeJSON(w, http.StatusOK, u)
}

func (h *UserHandler) Profile(w http.ResponseWriter, r *http.Request) {
	id, ok := parseIDParam(r, "id")
	if !ok {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid id"})
		return
	}
	p, err := h.users.Profile(id)
	if err != nil {
		if errors.Is(err, service.ErrUserNotFound) {
			writeJSON(w, http.StatusNotFound, map[string]string{"error": "not found"})
			return
		}
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "internal error"})
		return
	}
	writeJSON(w, http.StatusOK, p)
}

func (h *UserHandler) Create(w http.ResponseWriter, r *http.Request) {
	var body entity.User
	if err := decodeJSON(r, &body); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid body"})
		return
	}
	if body.Username == "" || body.Nickname == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "username and nickname required"})
		return
	}
	if err := h.users.Create(&body); err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "failed to create user"})
		return
	}
	writeJSON(w, http.StatusCreated, body)
}

func (h *UserHandler) Update(w http.ResponseWriter, r *http.Request) {
	id, ok := parseIDParam(r, "id")
	if !ok {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid id"})
		return
	}
	existing, err := h.users.Get(id)
	if err != nil {
		if errors.Is(err, service.ErrUserNotFound) {
			writeJSON(w, http.StatusNotFound, map[string]string{"error": "not found"})
			return
		}
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "internal error"})
		return
	}
	var body entity.User
	if err := decodeJSON(r, &body); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid body"})
		return
	}
	body.ID = existing.ID
	if err := h.users.Update(&body); err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "failed to update"})
		return
	}
	updated, _ := h.users.Get(id)
	writeJSON(w, http.StatusOK, updated)
}

func (h *UserHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id, ok := parseIDParam(r, "id")
	if !ok {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid id"})
		return
	}
	if err := h.users.Delete(id); err != nil {
		if errors.Is(err, service.ErrUserNotFound) {
			writeJSON(w, http.StatusNotFound, map[string]string{"error": "not found"})
			return
		}
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "internal error"})
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
