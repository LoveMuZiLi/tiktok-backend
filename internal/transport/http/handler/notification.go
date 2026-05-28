package handler

import (
	"errors"
	"net/http"

	"github.com/LoveMuZiLi/tiktok-backend/internal/entity"
	"github.com/LoveMuZiLi/tiktok-backend/internal/service"
)

type NotificationHandler struct {
	notifications *service.NotificationService
}

func NewNotificationHandler(notifications *service.NotificationService) *NotificationHandler {
	return &NotificationHandler{notifications: notifications}
}

func (h *NotificationHandler) List(w http.ResponseWriter, r *http.Request) {
	userID := parseUserID(r)
	items, err := h.notifications.List(userID)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "failed to list"})
		return
	}
	if items == nil {
		items = []entity.Notification{}
	}
	writeJSON(w, http.StatusOK, map[string]any{"items": items})
}

func (h *NotificationHandler) Create(w http.ResponseWriter, r *http.Request) {
	var body entity.Notification
	if err := decodeJSON(r, &body); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid body"})
		return
	}
	if body.UserID <= 0 {
		body.UserID = parseUserID(r)
	}
	if err := h.notifications.Create(&body); err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "failed to create"})
		return
	}
	writeJSON(w, http.StatusCreated, body)
}

func (h *NotificationHandler) MarkRead(w http.ResponseWriter, r *http.Request) {
	id, ok := parseIDParam(r, "id")
	if !ok {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid id"})
		return
	}
	if err := h.notifications.MarkRead(id); err != nil {
		if errors.Is(err, service.ErrNotificationNotFound) {
			writeJSON(w, http.StatusNotFound, map[string]string{"error": "not found"})
			return
		}
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "internal error"})
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

func (h *NotificationHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id, ok := parseIDParam(r, "id")
	if !ok {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid id"})
		return
	}
	if err := h.notifications.Delete(id); err != nil {
		if errors.Is(err, service.ErrNotificationNotFound) {
			writeJSON(w, http.StatusNotFound, map[string]string{"error": "not found"})
			return
		}
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "internal error"})
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
