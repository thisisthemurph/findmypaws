package routes

import (
	"errors"
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
	logger *slog.Logger) *ConversationHandler {
	return &ConversationHandler{
		ConversationRepo: conversationRepo,
		PetRepository:    petRepo,
		Logger:           logger,
	}
}

type ConversationHandler struct {
	PetRepository    repository.PetRepository
	ConversationRepo repository.ConversationRepository
	Logger           *slog.Logger
}

func (h *ConversationHandler) MakeRoutes(g *echo.Group) {
	g.GET("/conversations", h.ListConversations())
	g.GET("/conversations/:identifier", h.GetConversationByIdentifier())
}

type ConversationPetDetail struct {
	Name string `json:"name"`
	Type string `json:"type"`
}

type ConversationResponse struct {
	types.Conversation
	Pet   ConversationPetDetail `json:"pet"`
	Title string                `json:"title"`
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
			petDetail, _ := petLookup[conversation.Identifier]
			response[i] = ConversationResponse{
				Conversation: conversation,
				Pet:          petDetail,
				Title:        petDetail.Name,
			}
		}

		return c.JSON(http.StatusOK, response)
	}
}

func (h *ConversationHandler) GetConversationByIdentifier() echo.HandlerFunc {
	return func(c echo.Context) error {
		participantID := getParticipantID(c)
		identifier, err := uuid.Parse(c.Param("identifier"))
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, "invalid identifier")
		}

		conversation, err := h.ConversationRepo.Get(identifier, participantID)
		if err != nil {
			if errors.Is(err, repository.ErrNotFound) {
				return echo.NewHTTPError(http.StatusNotFound, "conversation not found")
			}
			return echo.NewHTTPError(http.StatusInternalServerError)
		}

		primaryUserID := conversation.PrimaryParticipantID
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

		petDetail := petLookup[conversation.Identifier]
		response := ConversationResponse{
			Conversation: *conversation,
			Pet:          petDetail,
			Title:        petDetail.Name,
		}

		return c.JSON(http.StatusOK, response)
	}
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
