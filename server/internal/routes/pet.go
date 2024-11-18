package routes

import (
	"encoding/json"
	"errors"
	"io"
	"log/slog"
	"net/http"
	"time"

	"github.com/google/uuid"
	"paws/internal/auth"
	"paws/internal/repository"
	"paws/internal/types"
	"paws/pkg/blight"
)

func NewPetsHandler(
	notificationRepo repository.NotificationRepository,
	petRepo repository.PetRepository,
	logger *slog.Logger,
) *PetsHandler {
	b, err := blight.New("./avatars.db")
	if err != nil {
		panic(err)
	}

	return &PetsHandler{
		NotificationRepo: notificationRepo,
		PetRepo:          petRepo,
		Blight:           b,
		Logger:           logger,
	}
}

type PetsHandler struct {
	NotificationRepo repository.NotificationRepository
	PetRepo          repository.PetRepository
	Blight           *blight.Client
	Logger           *slog.Logger
}

func (h *PetsHandler) RegisterRoutes(mux *http.ServeMux, mf MiddlewareFunc) {
	mux.HandleFunc("GET /api/v1/pets/{id}", mf(h.GetPetByID))
	mux.HandleFunc("GET /api/v1/pets", mf(h.ListPets))
	mux.HandleFunc("POST /api/v1/pets", mf(h.CreateNewPet))
	mux.HandleFunc("PUT /api/v1/pets/{id}", mf(h.UpdatePet))
	mux.HandleFunc("POST /api/v1/pets/{id}/tag", mf(h.AddTag))
	mux.HandleFunc("DELETE /api/v1/pets/{id}/tag/{key}", mf(h.DeleteTag))
	mux.HandleFunc("DELETE /api/v1/pets/{id}", mf(h.DeletePet))
	mux.HandleFunc("PUT /api/v1/pets/{id}/avatar", mf(h.UpdateAvatar))
	mux.HandleFunc("GET /api/v1/pets/{id}/avatar", mf(h.GetAvatar))
	mux.HandleFunc("POST /api/v1/pets/{id}/alert", mf(h.CreateNotificationOnPetPageVisit))
}

func (h *PetsHandler) GetPetByID(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		http.Error(w, "invalid pet id", http.StatusBadRequest)
		return
	}

	pet, err := h.PetRepo.Get(id)
	if err != nil {
		if errors.As(err, &repository.ErrNotFound) {
			http.Error(w, "pet not found", http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(pet)
}

func (h *PetsHandler) ListPets(w http.ResponseWriter, r *http.Request) {
	user := auth.GetUserFromContext(r.Context())
	if !user.Authenticated {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	pets, err := h.PetRepo.List(user.ID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(pets)
}

type NewPetRequest struct {
	Name string         `json:"name" validate:"required"`
	Type *types.PetType `json:"type"`
	Tags types.PetTags  `json:"tags"`
	DOB  *time.Time     `json:"dob"`
}

func (h *PetsHandler) CreateNewPet(w http.ResponseWriter, r *http.Request) {
	user := auth.GetUserFromContext(r.Context())
	if !user.Authenticated {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	var req NewPetRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.Logger.Error("bad request", "error", err)
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}

	currentPets, err := h.PetRepo.List(user.ID)
	if err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	for _, p := range currentPets {
		if p.Name == req.Name {
			http.Error(w, "A pet with that name already exists", http.StatusBadRequest)
			return
		}
	}

	pet := &types.Pet{
		UserID: user.ID,
		Name:   req.Name,
		Type:   req.Type,
		Tags:   req.Tags,
		DOB:    req.DOB,
	}

	if err := h.PetRepo.Create(pet); err != nil {
		h.Logger.Error("error creating pet", "error", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(pet)
}

type UpdatePetRequest struct {
	Type  *types.PetType `json:"type" validate:"required,max=16"`
	Name  string         `json:"name" validate:"required"`
	DOB   *time.Time     `json:"dob"`
	Blurb *string        `json:"blurb"`
}

func (h *PetsHandler) UpdatePet(w http.ResponseWriter, r *http.Request) {
	user := auth.GetUserFromContext(r.Context())
	if !user.Authenticated {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	var req UpdatePetRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.Logger.Error("bad request", "error", err)
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}
	if req.Name == "" {
		http.Error(w, "a name must be provided", http.StatusBadRequest)
		return
	}
	if req.Type == nil || len(string(*req.Type)) == 0 || len(string(*req.Type)) > 16 {
		http.Error(w, "invalid type", http.StatusBadRequest)
		return
	}

	petId, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		http.Error(w, "invalid pet id", http.StatusBadRequest)
		return
	}

	pet, err := h.PetRepo.Get(petId)
	if err != nil {
		h.Logger.Error("error getting pet", "error", err)
		if notFound := errors.As(err, &repository.ErrNotFound); notFound {
			http.Error(w, "pet not found", http.StatusNotFound)
			return
		}
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	if pet.UserID != user.ID {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	pet.Name = req.Name
	pet.Type = req.Type
	pet.Blurb = req.Blurb
	pet.DOB = req.DOB

	if err := h.PetRepo.Update(&pet); err != nil {
		h.Logger.Error("error updating pet", "error", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(pet)
}

type NewTagRequest struct {
	Key   string `json:"key" validate:"required"`
	Value string `json:"value" validate:"required"`
}

func (h *PetsHandler) AddTag(w http.ResponseWriter, r *http.Request) {
	user := auth.GetUserFromContext(r.Context())
	if !user.Authenticated {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	petID, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		http.Error(w, "invalid pet id", http.StatusBadRequest)
		return
	}

	var req NewTagRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}
	if req.Key == "" || req.Value == "" {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}

	pet, err := h.PetRepo.Get(petID)
	if err != nil {
		if errors.As(err, &repository.ErrNotFound) {
			http.Error(w, "pet not found", http.StatusNotFound)
			return
		}
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	if pet.UserID != user.ID {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
	}

	if pet.Tags == nil {
		pet.Tags = make(types.PetTags)
	}
	pet.Tags[req.Key] = req.Value
	if err := h.PetRepo.Update(&pet); err != nil {
		h.Logger.Error("error updating pet", "error", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(pet)
}

func (h *PetsHandler) DeleteTag(w http.ResponseWriter, r *http.Request) {
	user := auth.GetUserFromContext(r.Context())
	if !user.Authenticated {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	petID, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		http.Error(w, "invalid pet id", http.StatusBadRequest)
		return
	}
	tagKey := r.PathValue("key")

	pet, err := h.PetRepo.Get(petID)
	if err != nil {
		if errors.As(err, &repository.ErrNotFound) {
			http.Error(w, "pet not found", http.StatusNotFound)
			return
		}
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	if pet.UserID != user.ID {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	if pet.Tags == nil {
		pet.Tags = make(types.PetTags)
	}

	delete(pet.Tags, tagKey)
	if err := h.PetRepo.Update(&pet); err != nil {
		h.Logger.Error("error updating pet", "error", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(pet)
}

func (h *PetsHandler) DeletePet(w http.ResponseWriter, r *http.Request) {
	user := auth.GetUserFromContext(r.Context())
	if !user.Authenticated {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	petID, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		http.Error(w, "invalid pet id", http.StatusBadRequest)
		return
	}

	existingPet, err := h.PetRepo.Get(petID)
	if err != nil {
		if errors.As(err, &repository.ErrNotFound) {
			http.Error(w, "pet not found", http.StatusNotFound)
			return
		}
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	if existingPet.UserID != user.ID {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	if err := h.PetRepo.Delete(petID); err != nil {
		h.Logger.Error("error deleting pet", "error", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(existingPet)
}

func (h *PetsHandler) UpdateAvatar(w http.ResponseWriter, r *http.Request) {
	user := auth.GetUserFromContext(r.Context())
	if !user.Authenticated {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	petID, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		http.Error(w, "invalid pet id", http.StatusBadRequest)
		return
	}

	if err := r.ParseMultipartForm(10 << 20); err != nil {
		h.Logger.Error("error parsing multipart form", "error", err)
		http.Error(w, "unable to parse form", http.StatusBadRequest)
		return
	}

	file, _, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "file upload error", http.StatusBadRequest)
		return
	}
	defer file.Close()

	path := petID.String()
	if err := h.Blight.Add(path, file); err != nil {
		http.Error(w, "failed to save file", http.StatusBadRequest)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{
		"avatar_uri": path,
	})
}

func (h *PetsHandler) GetAvatar(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		http.Error(w, "invalid pet id", http.StatusBadRequest)
		return
	}

	result, err := h.Blight.Get(id.String())
	if err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	r.Header.Set("Content-Type", "image/jpeg")
	if _, err := io.Copy(w, result.BLOB); err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
	}
}

type NewAlertRequest struct {
	AlertingUserId          string `json:"user_id"`
	AlertingAnonymousUserId string `json:"anonymous_user_id"`
}

func (a NewAlertRequest) IsAnonymous() bool {
	return a.AlertingAnonymousUserId != "" && a.AlertingUserId == ""
}
func (h *PetsHandler) CreateNotificationOnPetPageVisit(w http.ResponseWriter, r *http.Request) {
	alertCreatedResponse := func(w http.ResponseWriter, created bool) {
		status := http.StatusOK
		if created {
			status = http.StatusCreated
		}

		w.WriteHeader(status)
		json.NewEncoder(w).Encode(map[string]bool{
			"alert_created": created,
		})
	}

	processAndValidateRequest := func(r *http.Request) (NewAlertRequest, uuid.UUID, error) {
		var req NewAlertRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			return req, uuid.Nil, err
		}
		if req.AlertingAnonymousUserId == "" && req.AlertingUserId == "" {
			h.Logger.Error("cannot create notification; user_id or anonymous_user_id is required", "req", req)
			return req, uuid.Nil, errors.New("user_id or anonymous_user_id is required")
		}
		petID, err := uuid.Parse(r.PathValue("id"))
		if err != nil {
			h.Logger.Error("error parsing petID", "petID", petID, "error", err)
			return req, uuid.Nil, errors.New("error parsing petID")
		}
		return req, petID, nil
	}

	makeSpottedPetNotificationModel := func(req NewAlertRequest, pet types.Pet) (types.NotificationModel, error) {
		spotterName := ""
		if !req.IsAnonymous() {
			spotterName = "a registered user"
		}

		detail := types.SpottedPetNotificationDetail{
			SpotterName: spotterName,
			PetID:       pet.ID,
			PetName:     pet.Name,
			IsAnonymous: req.IsAnonymous(),
		}

		detailJSON, err := json.Marshal(detail)
		if err != nil {
			h.Logger.Error("failed to marshal notification detail JSON", "detail", detail, "error", err)
			return types.NotificationModel{}, errors.New("error marshalling notification detail")
		}

		return types.NotificationModel{
			UserID: pet.UserID,
			Type:   types.SpottedPetNotification,
			Detail: detailJSON,
		}, nil
	}

	user := auth.GetUserFromContext(r.Context())

	req, petID, err := processAndValidateRequest(r)
	if err != nil {
		h.Logger.Error("invalid request", "err", err)
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	pet, err := h.PetRepo.Get(petID)
	if err != nil {
		if errors.As(err, &repository.ErrNotFound) {
			http.Error(w, "pet not found", http.StatusNotFound)
			return
		}
		h.Logger.Error("failed to get pet", "petID", petID, "error", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	if pet.UserID == user.ID {
		// Ensure the user is not creating alerts for themselves.
		alertCreatedResponse(w, false)
		return
	}

	model, err := makeSpottedPetNotificationModel(req, pet)
	if err != nil {
		h.Logger.Error("failed to create notification model", "err", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	recentlyNotified, err := h.NotificationRepo.RecentlyNotified(model)
	if err != nil {
		h.Logger.Error("error determining if recently notified", "error", err)
	}
	if recentlyNotified {
		alertCreatedResponse(w, false)
		return
	}

	if err := h.NotificationRepo.Create(&model); err != nil {
		h.Logger.Error("error creating notification", "error", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	alertCreatedResponse(w, true)
}
