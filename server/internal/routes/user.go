package routes

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"log/slog"
	"net/http"
	"paws/internal/repository"
	"paws/internal/types"
)

func NewUsersHandler(alertRepo repository.AlertRepository, petRepo repository.PetRepository, logger *slog.Logger) *UsersHandler {
	return &UsersHandler{
		AlertRepo: alertRepo,
		PetRepo:   petRepo,
		Logger:    logger,
	}
}

type UsersHandler struct {
	AlertRepo repository.AlertRepository
	PetRepo   repository.PetRepository
	Logger    *slog.Logger
}

func (h *UsersHandler) MakeRoutes(g *echo.Group) {
	g.GET("/user/alerts", h.ListNotifications())
	g.POST("/user/notifications/read-all", h.MarkAllNotificationsAsRead())
}

func (h *UsersHandler) ListNotifications() echo.HandlerFunc {
	return func(c echo.Context) error {
		user := clerkUser(c)
		if !user.Authenticated {
			return echo.NewHTTPError(http.StatusUnauthorized)
		}

		alerts, err := h.AlertRepo.List(user.ID)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError)
		}

		pets, err := h.PetRepo.List(user.ID)
		petsLookup := make(map[uuid.UUID]string)
		for _, pet := range pets {
			petsLookup[pet.ID] = pet.Name
		}

		notifications := make([]types.Notification, len(alerts))
		for i, alert := range alerts {
			petName, exists := petsLookup[alert.PetId]
			if !exists {
				petName = "Your pet"
			}
			notifications[i] = types.Notification{
				ID:        fmt.Sprintf("alert_%d", alert.ID),
				Type:      types.SpottedPetNotification,
				Message:   fmt.Sprintf("%s was spotted by an anonymous user.", petName),
				Seen:      alert.SeenAt != nil,
				CreatedAt: alert.CreatedAt,
				Links: map[string]string{
					"pet": fmt.Sprintf("/pet/%v", alert.PetId),
				},
			}
		}
		return c.JSON(http.StatusOK, notifications)
	}
}

func (h *UsersHandler) MarkAllNotificationsAsRead() echo.HandlerFunc {
	return func(c echo.Context) error {
		user := clerkUser(c)
		if !user.Authenticated {
			return echo.NewHTTPError(http.StatusUnauthorized)
		}

		if err := h.AlertRepo.MarkAllAsRead(user.ID); err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError)
		}
		return c.NoContent(http.StatusNoContent)
	}
}
