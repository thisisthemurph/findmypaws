package routes

import (
	"errors"
	"log/slog"
	"net/http"

	"github.com/google/uuid"
	"paws/pkg/chat"
)

type ChatHandler struct {
	manager *chat.Manager
	logger  *slog.Logger
}

func NewChatHandler(manager *chat.Manager, logger *slog.Logger) *ChatHandler {
	return &ChatHandler{
		manager: manager,
		logger:  logger,
	}
}

func (h *ChatHandler) RegisterRoutes(mux *http.ServeMux, mf MiddlewareFunc) {
	mux.HandleFunc("GET /room", h.HandleRoom)
}

func (h *ChatHandler) HandleRoom(w http.ResponseWriter, r *http.Request) {
	participantID := r.URL.Query().Get("pid")
	if participantID == "" {
		http.Error(w, "Missing required parameter pid", http.StatusBadRequest)
		return
	}
	roomID := r.URL.Query().Get("r")
	if roomID == "" {
		http.Error(w, "Missing Room ID", http.StatusBadRequest)
		return
	}

	roomIdentifier, err := uuid.Parse(roomID)
	if err != nil {
		http.Error(w, "Invalid Room ID", http.StatusBadRequest)
		return
	}

	// Retrieve or create a new Room for the given Room ID.
	room, err := h.manager.GetOrCreateRoom(roomIdentifier, participantID)
	if err != nil {
		if errors.Is(err, chat.ErrUnauthorized) {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	room.ServeWS(w, r)
}
