package routes

import (
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"log/slog"
	"net/http"
	"paws/internal/auth"
	"paws/internal/repository"
	"paws/internal/types"
	"sync"
)

type UsersHandler struct {
	UserRepo         repository.UserRepository
	NotificationRepo repository.NotificationRepository
	PetRepo          repository.PetRepository
	Logger           *slog.Logger
}

func NewUsersHandler(
	userRepo repository.UserRepository,
	notificationRepo repository.NotificationRepository,
	petRepo repository.PetRepository,
	logger *slog.Logger,
) *UsersHandler {
	return &UsersHandler{
		UserRepo:         userRepo,
		NotificationRepo: notificationRepo,
		PetRepo:          petRepo,
		Logger:           logger,
	}
}

func (h *UsersHandler) RegisterRoutes(mux *http.ServeMux, mf MiddlewareFunc) {
	mux.HandleFunc("GET /api/v1/user/notifications", mf(h.ListNotifications))
	mux.HandleFunc("POST /api/v1/user/notifications/read-all", mf(h.MarkAllNotificationsAsSeen))
	mux.HandleFunc("PUT /api/v1/user/anonymous/{id}", mf(h.UpdateAnonymousUser))

	mux.HandleFunc("GET /api/v1/test", mf(h.Test))
}

type UpdateAnonymousUserRequest struct {
	Name string `json:"name"`
}

func (h *UsersHandler) Test(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(struct {
		Name string `json:"name"`
	}{
		Name: "this is only a test",
	})
}

func (h *UsersHandler) UpdateAnonymousUser(w http.ResponseWriter, r *http.Request) {
	// Extract anonymous user ID from URL
	anonymousUserId := r.PathValue("id")

	// Parse the request body
	var req UpdateAnonymousUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}
	if err := validateRequest(req); err != nil {
		http.Error(w, "Validation Error", http.StatusBadRequest)
		return
	}

	user := types.AnonymousUser{
		ID:   anonymousUserId,
		Name: req.Name,
	}
	if err := h.UserRepo.UpsertAnonymousUser(&user); err != nil {
		h.Logger.Error("error upserting anonymous user", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(user)
}

func (h *UsersHandler) ListNotifications(w http.ResponseWriter, r *http.Request) {
	user := auth.GetUserFromContext(r.Context())
	if !user.Authenticated {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	notificationModels, err := h.NotificationRepo.List(user.ID)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	pets, err := h.PetRepo.List(user.ID)
	petsLookup := make(map[uuid.UUID]string)
	for _, pet := range pets {
		petsLookup[pet.ID] = pet.Name
	}

	notifications := make([]types.Notification, len(notificationModels))
	var wg sync.WaitGroup
	wg.Add(len(notificationModels))
	for i, model := range notificationModels {
		go func(i int, model types.NotificationModel) {
			defer wg.Done()
			n, ok := model.Notification()
			if !ok {
				h.Logger.Error("error parsing notification model", "model", model)
				return
			}
			notifications[i] = n
		}(i, model)
	}

	wg.Wait()

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(notifications)
}

func (h *UsersHandler) MarkAllNotificationsAsSeen(w http.ResponseWriter, r *http.Request) {
	user := auth.GetUserFromContext(r.Context())
	if !user.Authenticated {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	if err := h.NotificationRepo.MarkAllSeen(user.ID); err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func validateRequest(req UpdateAnonymousUserRequest) error {
	// Perform necessary validation, e.g., ensuring the name is not empty
	if req.Name == "" {
		return fmt.Errorf("name is required")
	}
	return nil
}
