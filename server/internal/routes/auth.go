package routes

import (
	"github.com/labstack/echo/v4"
	"log/slog"
	"net/http"
	"paws/internal/store"
)

func NewAuthHandler(s *store.PostgresStore, logger *slog.Logger) *AuthHandler {
	return &AuthHandler{
		UserStore: s.UserStore,
		Logger:    logger,
	}
}

type AuthHandler struct {
	UserStore store.UserStore
	Logger    *slog.Logger
}

func (h *AuthHandler) MakeRoutes(g *echo.Group) {
	g.POST("/auth/signup", h.HandleSignUp())
}

type SignUpRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (h *AuthHandler) HandleSignUp() echo.HandlerFunc {
	return func(c echo.Context) error {
		var req SignUpRequest
		if err := c.Bind(&req); err != nil {
			return echo.NewHTTPError(http.StatusBadRequest)
		}

		user, err := h.UserStore.SignUp(req.Email, req.Password, req.Name)
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err)
		}

		return c.JSON(http.StatusOK, user)
	}
}
