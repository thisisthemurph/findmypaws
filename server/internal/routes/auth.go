package routes

import (
	"github.com/labstack/echo/v4"
	"log/slog"
	"net/http"
	"paws/internal/store"
)

func NewAuthHandler(s *store.PostgresStore, logger *slog.Logger) *AuthHandler {
	return &AuthHandler{
		AuthStore: s.AuthStore,
		Logger:    logger,
	}
}

type AuthHandler struct {
	AuthStore store.AuthStore
	Logger    *slog.Logger
}

func (h *AuthHandler) MakeRoutes(g *echo.Group) {
	g.POST("/auth/login", h.HandleLogIn())
	g.POST("/auth/logout", h.HandleLogOut())
	g.POST("/auth/signup", h.HandleSignUp())
	g.POST("/auth/refresh", h.HandleRefreshToken())
}

type SignUpRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (h *AuthHandler) HandleLogIn() echo.HandlerFunc {
	return func(c echo.Context) error {
		var req LoginRequest
		if err := c.Bind(&req); err != nil {
			return echo.NewHTTPError(http.StatusBadRequest)
		}

		session, err := h.AuthStore.LogIn(req.Email, req.Password)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError)
		}
		return c.JSON(http.StatusOK, session)
	}
}

func (h *AuthHandler) HandleLogOut() echo.HandlerFunc {
	return func(c echo.Context) error {
		user := CurrentUser(c)
		if err := h.AuthStore.LogOut(user.Token); err != nil {
			h.Logger.Error("error logging out", "user", user.ID, "error", err)
			return echo.NewHTTPError(http.StatusInternalServerError)
		}
		return c.NoContent(http.StatusNoContent)
	}
}

func (h *AuthHandler) HandleSignUp() echo.HandlerFunc {
	return func(c echo.Context) error {
		var req SignUpRequest
		if err := c.Bind(&req); err != nil {
			return echo.NewHTTPError(http.StatusBadRequest)
		}

		user, err := h.AuthStore.SignUp(req.Email, req.Password, req.Name)
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err)
		}

		return c.JSON(http.StatusOK, user)
	}
}

func (h *AuthHandler) HandleRefreshToken() echo.HandlerFunc {
	return func(c echo.Context) error {
		token := c.Request().Header.Get("X-Refresh-Token")
		if token == "" {
			return echo.NewHTTPError(http.StatusBadRequest)
		}

		session, err := h.AuthStore.RefreshToken(token)
		if err != nil {
			return echo.NewHTTPError(http.StatusUnauthorized)
		}
		return c.JSON(http.StatusOK, session)
	}
}
