package handler

import (
	"net/http"

	"github.com/LoveMuZiLi/tiktok-backend/internal/entity"
	"github.com/LoveMuZiLi/tiktok-backend/internal/service"
)

type InboxHandler struct {
	inbox *service.InboxService
}

func NewInboxHandler(inbox *service.InboxService) *InboxHandler {
	return &InboxHandler{inbox: inbox}
}

func (h *InboxHandler) ListChats(w http.ResponseWriter, r *http.Request) {
	userID := parseUserID(r)
	chats, err := h.inbox.ListChats(userID)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "failed to list chats"})
		return
	}
	if chats == nil {
		chats = []entity.ChatSummary{}
	}
	writeJSON(w, http.StatusOK, map[string]any{"items": chats})
}

func (h *InboxHandler) ListMessages(w http.ResponseWriter, r *http.Request) {
	convID, ok := parseIDParam(r, "id")
	if !ok {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid id"})
		return
	}
	userID := parseUserID(r)
	msgs, err := h.inbox.ListMessages(convID, userID)
	if err != nil {
		if err.Error() == "forbidden" {
			writeJSON(w, http.StatusForbidden, map[string]string{"error": "forbidden"})
			return
		}
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "failed to list messages"})
		return
	}
	if msgs == nil {
		msgs = []entity.MessageView{}
	}
	writeJSON(w, http.StatusOK, map[string]any{"items": msgs})
}

type sendMessageBody struct {
	Text string `json:"text"`
}

func (h *InboxHandler) SendMessage(w http.ResponseWriter, r *http.Request) {
	convID, ok := parseIDParam(r, "id")
	if !ok {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid id"})
		return
	}
	var body sendMessageBody
	if err := decodeJSON(r, &body); err != nil || body.Text == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "text required"})
		return
	}
	userID := parseUserID(r)
	msg, err := h.inbox.SendMessage(convID, userID, body.Text)
	if err != nil {
		if err.Error() == "forbidden" {
			writeJSON(w, http.StatusForbidden, map[string]string{"error": "forbidden"})
			return
		}
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "failed to send"})
		return
	}
	writeJSON(w, http.StatusCreated, msg)
}

type openChatBody struct {
	PeerID int64 `json:"peerId"`
}

func (h *InboxHandler) OpenChat(w http.ResponseWriter, r *http.Request) {
	var body openChatBody
	if err := decodeJSON(r, &body); err != nil || body.PeerID <= 0 {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "peerId required"})
		return
	}
	userID := parseUserID(r)
	convID, err := h.inbox.OpenOrCreateChat(userID, body.PeerID)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "failed to open chat"})
		return
	}
	writeJSON(w, http.StatusOK, map[string]int64{"conversationId": convID})
}

func (h *InboxHandler) DeleteMessage(w http.ResponseWriter, r *http.Request) {
	msgID, ok := parseIDParam(r, "messageId")
	if !ok {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid message id"})
		return
	}
	if err := h.inbox.DeleteMessage(msgID); err != nil {
		writeJSON(w, http.StatusNotFound, map[string]string{"error": "not found"})
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
