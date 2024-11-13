package routes

import (
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"log/slog"
	"net/http"
	"paws/internal/repository"
	"paws/internal/types"
	"sync"
)

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

type UsersHandler struct {
	UserRepo         repository.UserRepository
	NotificationRepo repository.NotificationRepository
	PetRepo          repository.PetRepository
	Logger           *slog.Logger
}

func (h *UsersHandler) MakeRoutes(g *echo.Group) {
	g.GET("/user/notifications", h.ListNotifications())
	g.POST("/user/notifications/read-all", h.MarkAllNotificationsAsSeen())
	g.PUT("/user/anonymous/:id", h.UpdateAnonymousUser())
}

type UpdateAnonymousUserRequest struct {
	Name string `json:"name" validate:"required"`
}

func (h *UsersHandler) UpdateAnonymousUser() echo.HandlerFunc {
	return func(c echo.Context) error {
		anonymousUserId := c.Param("id")

		var req UpdateAnonymousUserRequest
		if err := c.Bind(&req); err != nil {
			return echo.NewHTTPError(http.StatusBadRequest)
		}
		if err := c.Validate(req); err != nil {
			return echo.NewHTTPError(http.StatusBadRequest)
		}

		user := types.AnonymousUser{
			ID:   anonymousUserId,
			Name: req.Name,
		}
		if err := h.UserRepo.UpsertAnonymousUser(&user); err != nil {
			h.Logger.Error("error upserting anonymous user", "error", err)
			return echo.NewHTTPError(http.StatusInternalServerError)
		}

		return c.JSON(http.StatusOK, user)
	}
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
