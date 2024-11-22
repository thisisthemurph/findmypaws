package routes

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"

	"github.com/clerk/clerk-sdk-go/v2"
	svix "github.com/svix/svix-webhooks/go"
	"paws/internal/repository"
	"paws/internal/response/clerktype"
)

type WebhookHandler struct {
	ClerkSigningSecret string
	UserRepo           repository.UserRepository
	Logger             *slog.Logger
	Handlers           map[string]func(clerktype.WebhookEvent) (int, error)
}

func NewWebhookHandler(clerkSigningSecret string, userRepo repository.UserRepository, logger *slog.Logger) *WebhookHandler {
	whh := &WebhookHandler{
		ClerkSigningSecret: clerkSigningSecret,
		UserRepo:           userRepo,
		Logger:             logger,
	}

	whh.Handlers = map[string]func(clerktype.WebhookEvent) (int, error){
		"user.created": whh.handleUserCreatedEvent,
		"user.deleted": whh.handleUserDeletedEvent,
	}
	return whh
}

func (h *WebhookHandler) RegisterRoutes(mux *http.ServeMux, mf MiddlewareFunc) {
	mux.HandleFunc("POST /api/webhooks", mf(h.HandleClerkWebhook))
}

func (h *WebhookHandler) HandleClerkWebhook(w http.ResponseWriter, r *http.Request) {
	wh, err := svix.NewWebhook(h.ClerkSigningSecret)
	if err != nil {
		http.Error(w, "error creating webhook", http.StatusInternalServerError)
		return
	}

	svixID := r.Header.Get("svix-id")
	svixTimestamp := r.Header.Get("svix-timestamp")
	svixSignature := r.Header.Get("svix-signature")

	if svixID == "" || svixTimestamp == "" || svixSignature == "" {
		http.Error(w, "missing svix headers", http.StatusBadRequest)
		return
	}

	r.Body = http.MaxBytesReader(w, r.Body, 1_048_576)
	var event clerktype.WebhookEvent
	byteData, err := io.ReadAll(r.Body)
	if err != nil {
		h.Logger.Error("error parsing webhook event JSON", "error", err)
		http.Error(w, "error parsing webhook event JSON", http.StatusBadRequest)
		return
	}

	if err := wh.Verify(byteData, r.Header); err != nil {
		h.Logger.Error("error verifying webhook", "error", err)
		http.Error(w, "error verifying webhook", http.StatusBadRequest)
		return
	}

	if err := json.Unmarshal(byteData, &event); err != nil {
		http.Error(w, "error reading body", http.StatusBadRequest)
		return
	}

	h.handleEvent(w, event)
	return
}

func (h *WebhookHandler) handleEvent(w http.ResponseWriter, event clerktype.WebhookEvent) {
	handler, ok := h.Handlers[event.Type]
	if !ok {
		h.Logger.Error("unknown event type", "type", event.Type)
		http.Error(w, fmt.Sprintf("unknown event type: %s", event.Type), http.StatusBadRequest)
		return
	}

	if status, err := handler(event); err != nil {
		h.Logger.Error("error handling event", "type", event.Type, "error", err)
		http.Error(w, fmt.Sprintf("error handling event: %s", event.Type), status)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (h *WebhookHandler) handleUserCreatedEvent(event clerktype.WebhookEvent) (int, error) {
	if event.Type != "user.created" {
		return http.StatusBadRequest, fmt.Errorf("invalid event type: %s", event.Type)
	}

	var usr clerk.User
	if err := json.Unmarshal(event.Data, &usr); err != nil {
		return http.StatusBadRequest, fmt.Errorf("error parsing event data: %w", err)
	}
	if usr.Object != "user" {
		return http.StatusBadRequest, fmt.Errorf("invalid user object: %s", usr.Object)
	}

	if err := h.UserRepo.UpsertUser(usr); err != nil {
		return http.StatusInternalServerError, err
	}

	fmt.Printf("%+v\n", usr)
	return http.StatusOK, nil
}

func (h *WebhookHandler) handleUserDeletedEvent(event clerktype.WebhookEvent) (int, error) {
	if event.Type != "user.deleted" {
		return http.StatusBadRequest, fmt.Errorf("invalid event type: %s", event.Type)
	}

	var deletion clerktype.UserDeletedWebhookEventData
	if err := json.Unmarshal(event.Data, &deletion); err != nil {
		return http.StatusBadRequest, fmt.Errorf("error parsing event data: %w", err)
	}
	if deletion.Object != "user" {
		return http.StatusBadRequest, fmt.Errorf("invalid user object: %s", deletion.Object)
	}
	if !deletion.Deleted {
		return http.StatusBadRequest, errors.New("error deleting user, event indicated user is not deleted")
	}

	if err := h.UserRepo.DeleteUser(deletion.UserId); err != nil {
		return http.StatusInternalServerError, err
	}
	return http.StatusOK, nil
}
