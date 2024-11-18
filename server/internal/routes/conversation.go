package routes

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"log/slog"
	"net/http"
	"paws/internal/auth"
	"paws/internal/repository"
	"paws/internal/types"
)

func NewConversationHandler(
	conversationRepo repository.ConversationRepository,
	petRepo repository.PetRepository,
	userRepo repository.UserRepository,
	logger *slog.Logger) *ConversationHandler {
	return &ConversationHandler{
		ConversationRepo: conversationRepo,
		PetRepository:    petRepo,
		UserRepo:         userRepo,
		Logger:           logger,
	}
}

type ConversationHandler struct {
	PetRepository    repository.PetRepository
	ConversationRepo repository.ConversationRepository
	UserRepo         repository.UserRepository
	Logger           *slog.Logger
}

func (h *ConversationHandler) RegisterRoutes(mux *http.ServeMux, mf MiddlewareFunc) {
	mux.HandleFunc("GET /api/v1/conversations", mf(h.ListConversations))
	mux.HandleFunc("GET /api/v1/conversations/{identifier}", mf(h.GetConversationByIdentifier))
	mux.HandleFunc("POST /api/v1/conversations", mf(h.CreateIfNotExists))
}

type ConversationPetDetail struct {
	Name string `json:"name"`
	Type string `json:"type"`
}

type ConversationParticipant struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type ConversationResponse struct {
	types.Conversation
	Pet              ConversationPetDetail   `json:"pet"`
	Title            string                  `json:"title"`
	Participant      ConversationParticipant `json:"participant"`
	OtherParticipant ConversationParticipant `json:"otherParticipant"`
}

type CreateConversationRequest struct {
	Identifier    uuid.UUID `json:"identifier" validate:"required"`
	ParticipantId string    `json:"participantId" validate:"required"`
}

func (h *ConversationHandler) CreateIfNotExists(w http.ResponseWriter, r *http.Request) {
	var req CreateConversationRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.Logger.Error("error getting/creating the conversation", "identifier", req.Identifier, "participant", req.ParticipantId, "error", err)
		return
	}
	if req.Identifier == uuid.Nil || req.ParticipantId == "" {
		http.Error(w, "bad request", http.StatusBadRequest)
	}

	if _, err := h.ConversationRepo.GetOrCreate(req.Identifier, req.ParticipantId); err != nil {
		h.Logger.Error("error getting/creating the conv")
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// ListConversations lists all conversations for the current conversation participant.
func (h *ConversationHandler) ListConversations(w http.ResponseWriter, r *http.Request) {
	participantID, err := getParticipantIDFromRequest(r)
	if err != nil {
		h.Logger.Error("failed to determine participant ID", "error", err)
	}

	conversations, err := h.ConversationRepo.List(participantID)
	if err != nil {
		h.Logger.Error("failed to list conversations", "error", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if len(conversations) == 0 {
		json.NewEncoder(w).Encode([]ConversationResponse{})
		return
	}

	primaryUserID := conversations[0].PrimaryParticipantID
	pets, err := h.PetRepository.List(primaryUserID)
	if err != nil {
		h.Logger.Error("failed to list conversations", "error", err)
	}

	petLookup := make(map[uuid.UUID]ConversationPetDetail)
	for _, pet := range pets {
		petLookup[pet.ID] = ConversationPetDetail{
			Name: pet.Name,
			Type: string(*pet.Type),
		}
	}

	response := make([]ConversationResponse, len(conversations))
	for i, conversation := range conversations {
		petDetail, petFound := petLookup[conversation.Identifier]
		participant, otherParticipant := h.getParticipants(participantID, conversation, petLookup)
		title := otherParticipant.Name
		if petFound {
			title = fmt.Sprintf("%s - %s", otherParticipant.Name, petDetail.Name)
		}

		response[i] = ConversationResponse{
			Conversation:     conversation,
			Pet:              petDetail,
			Participant:      participant,
			OtherParticipant: otherParticipant,
			Title:            title,
		}
	}

	json.NewEncoder(w).Encode(response)
}

func (h *ConversationHandler) GetConversationByIdentifier(w http.ResponseWriter, r *http.Request) {
	currentParticipantID, err := getParticipantIDFromRequest(r)
	if err != nil {
		h.Logger.Error("failed to determine participant ID", "error", err)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	identifier, err := uuid.Parse(r.PathValue("identifier"))
	if err != nil {
		h.Logger.Error("failed to determine participant ID", "error", err)
		http.Error(w, "invalid identifier", http.StatusBadRequest)
		return
	}

	conversation, err := h.ConversationRepo.Get(identifier, currentParticipantID)
	if err != nil {
		h.Logger.Error("failed to determine conversation", "error", err)
		if errors.Is(err, repository.ErrNotFound) {
			http.Error(w, "conversation not found", http.StatusNotFound)
			return
		}
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	petLookup := h.getPetLookup(conversation.PrimaryParticipantID)

	petDetail, petFound := petLookup[conversation.Identifier]
	currentParticipant, otherParticipant := h.getParticipants(currentParticipantID, *conversation, petLookup)

	title := otherParticipant.Name
	if petFound {
		title = fmt.Sprintf("%s - %s", otherParticipant.Name, petDetail.Name)
	}

	response := ConversationResponse{
		Conversation:     *conversation,
		Pet:              petDetail,
		Participant:      currentParticipant,
		OtherParticipant: otherParticipant,
		Title:            title,
	}

	json.NewEncoder(w).Encode(response)
}

func (h *ConversationHandler) getParticipants(currentParticipantID string, conversation types.Conversation, petLookup map[uuid.UUID]ConversationPetDetail) (ConversationParticipant, ConversationParticipant) {
	participant := ConversationParticipant{}
	otherParticipant := ConversationParticipant{}

	petDetail, petFound := petLookup[conversation.Identifier]
	if !petFound {
		h.Logger.Error("pet details not found", "conversation_id", conversation.Identifier)
	}

	// Fetch secondary participant name
	secondaryParticipantName := "Anonymous"
	secondaryParticipant, err := h.UserRepo.GetAnonymousUser(conversation.SecondaryParticipantID)
	if err != nil {
		h.Logger.Error("failed to get anonymous user", "error", err)
	} else if secondaryParticipant.Name != "" {
		secondaryParticipantName = secondaryParticipant.Name
	}

	if currentParticipantID == conversation.PrimaryParticipantID {
		participant.ID = conversation.PrimaryParticipantID
		participant.Name = "You"
		otherParticipant.ID = conversation.SecondaryParticipantID
		otherParticipant.Name = secondaryParticipantName
		return participant, otherParticipant
	}

	participant.ID = conversation.SecondaryParticipantID
	participant.Name = secondaryParticipantName
	otherParticipant.ID = conversation.PrimaryParticipantID
	if petFound {
		otherParticipant.Name = petDetail.Name
	} else {
		otherParticipant.Name = "Owner"
	}
	return participant, otherParticipant
}

func (h *ConversationHandler) getPetLookup(primaryParticipantID string) map[uuid.UUID]ConversationPetDetail {
	petLookup := make(map[uuid.UUID]ConversationPetDetail)
	pets, err := h.PetRepository.List(primaryParticipantID)
	if err != nil {
		h.Logger.Error("failed to list conversations", "error", err)
		return petLookup
	}

	for _, pet := range pets {
		petLookup[pet.ID] = ConversationPetDetail{
			Name: pet.Name,
			Type: string(*pet.Type),
		}
	}
	return petLookup
}

func getParticipantIDFromRequest(r *http.Request) (string, error) {
	user := auth.GetUserFromContext(r.Context())
	if user.Authenticated {
		return user.ID, nil
	}
	anonymousUserID := r.Header.Get("AnonymousUserId")
	if anonymousUserID == "" {
		return "", errors.New("could not determine participant ID")
	}
	return anonymousUserID, nil
}
