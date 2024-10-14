package routes

import (
	"github.com/labstack/echo/v4"
	"log/slog"
	"paws/internal/store"
)

func NewUserHandler(s *store.PostgresStore, logger *slog.Logger) *UserHandler {
	return &UserHandler{
		UserStore: s.UserStore,
		Logger:    logger,
	}
}

type UserHandler struct {
	UserStore store.UserStore
	Logger    *slog.Logger
}

func (h *UserHandler) MakeRoutes(g *echo.Group) {}
