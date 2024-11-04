package routes

import (
	"log/slog"
	"net/http"
	"sync"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"paws/internal/repository"
	"paws/internal/types"
)

func NewUsersHandler(notificationRepo repository.NotificationRepository, petRepo repository.PetRepository, logger *slog.Logger) *UsersHandler {
	return &UsersHandler{
		NotificationRepo: notificationRepo,
		PetRepo:          petRepo,
		Logger:           logger,
	}
}

type UsersHandler struct {
	NotificationRepo repository.NotificationRepository
	PetRepo          repository.PetRepository
	Logger           *slog.Logger
}

func (h *UsersHandler) MakeRoutes(g *echo.Group) {
	g.GET("/user/notifications", h.ListNotifications())
	g.POST("/user/notifications/read-all", h.MarkAllNotificationsAsSeen())
}

func (h *UsersHandler) ListNotifications() echo.HandlerFunc {
	return func(c echo.Context) error {
		user := clerkUser(c)
		if !user.Authenticated {
			return echo.NewHTTPError(http.StatusUnauthorized)
		}

		notificationModels, err := h.NotificationRepo.List(user.ID)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError)
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
		return c.JSON(http.StatusOK, notifications)
	}
}

func (h *UsersHandler) MarkAllNotificationsAsSeen() echo.HandlerFunc {
	return func(c echo.Context) error {
		user := clerkUser(c)
		if !user.Authenticated {
			return echo.NewHTTPError(http.StatusUnauthorized)
		}

		if err := h.NotificationRepo.MarkAllSeen(user.ID); err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError)
		}
		return c.NoContent(http.StatusNoContent)
	}
}
