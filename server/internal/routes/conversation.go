package routes

import (
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"log/slog"
	"net/http"
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

func (h *ConversationHandler) MakeRoutes(g *echo.Group) {
	g.POST("/conversations", h.CreateIfNotExists())
	g.GET("/conversations", h.ListConversations())
	g.GET("/conversations/:identifier", h.GetConversationByIdentifier())
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

func (h *ConversationHandler) CreateIfNotExists() echo.HandlerFunc {
	return func(c echo.Context) error {
		var req CreateConversationRequest
		if err := c.Bind(&req); err != nil {
			return echo.NewHTTPError(http.StatusBadRequest)
		}
		if err := c.Validate(req); err != nil {
			return echo.NewHTTPError(http.StatusBadRequest)
		}

		if _, err := h.ConversationRepo.GetOrCreate(req.Identifier, req.ParticipantId); err != nil {
			h.Logger.Error("error getting/creating the conversation", "identifier", req.Identifier, "participant", req.ParticipantId, "error", err)
			return echo.NewHTTPError(http.StatusInternalServerError)
		}

		return c.NoContent(http.StatusNoContent)
	}
}

// ListConversations lists all conversations for the current conversation participant.
func (h *ConversationHandler) ListConversations() echo.HandlerFunc {
	return func(c echo.Context) error {
		participantID := getParticipantID(c)
		if participantID == "" {
			return echo.NewHTTPError(http.StatusBadRequest, "participant information missing")
		}

		conversations, err := h.ConversationRepo.List(participantID)
		if err != nil {
			h.Logger.Error("failed to list conversations", "error", err)
			return echo.NewHTTPError(http.StatusInternalServerError)
		}

		if len(conversations) == 0 {
			return c.JSON(http.StatusOK, []ConversationResponse{})
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

		return c.JSON(http.StatusOK, response)
	}
}

func (h *ConversationHandler) GetConversationByIdentifier() echo.HandlerFunc {
	return func(c echo.Context) error {
		currentParticipantID := getParticipantID(c)
		identifier, err := uuid.Parse(c.Param("identifier"))
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, "invalid identifier")
		}

		conversation, err := h.ConversationRepo.Get(identifier, currentParticipantID)
		if err != nil {
			if errors.Is(err, repository.ErrNotFound) {
				return echo.NewHTTPError(http.StatusNotFound, "conversation not found")
			}
			return echo.NewHTTPError(http.StatusInternalServerError)
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

		return c.JSON(http.StatusOK, response)
	}
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

func getAnonymousUserID(c echo.Context) string {
	return c.Request().Header.Get("AnonymousUserId")
}

// getParticipantID returns the ID of the conversation participant.
// If the participant is authenticated, their userID is returned, otherwise their anonymous user ID is returned.
func getParticipantID(c echo.Context) string {
	user := clerkUser(c)
	if user.Authenticated {
		return user.ID
	}
	return getAnonymousUserID(c)
}
