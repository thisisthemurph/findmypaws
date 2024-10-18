package routes

import (
	"errors"
	"log/slog"
	"net/http"

	"github.com/labstack/echo/v4"
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

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
}

func (h *AuthHandler) HandleLogIn() echo.HandlerFunc {
	defaultErr := "There has been an issue signing you in, please ensure the form is complete and try again."
	return func(c echo.Context) error {
		var req LoginRequest
		if err := c.Bind(&req); err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, defaultErr)
		}
		if err := c.Validate(req); err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, defaultErr)
		}

		session, err := h.AuthStore.LogIn(req.Email, req.Password)
		if err != nil {
			if errors.Is(err, store.ErrInvalidCredentials) {
				return echo.NewHTTPError(http.StatusUnauthorized, "The email and password combination does not match our records.")
			}
			return echo.NewHTTPError(http.StatusInternalServerError, "There has been an unexpected error signing you in.")
		}
		return c.JSON(http.StatusOK, session)
	}
}

func (h *AuthHandler) HandleLogOut() echo.HandlerFunc {
	return func(c echo.Context) error {
		user := CurrentUser(c)
		if err := h.AuthStore.LogOut(user.Token); err != nil {
			h.Logger.Error("error logging out", "user", user.ID, "error", err)
			return echo.NewHTTPError(http.StatusInternalServerError, "There has been an unexpected error logging you out")
		}
		return c.NoContent(http.StatusNoContent)
	}
}

type SignUpRequest struct {
	Name     string `json:"username" validate:"required"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

func (h *AuthHandler) HandleSignUp() echo.HandlerFunc {
	defaultErr := "There has been an issue signing you up, please ensure the form is complete and try again."
	return func(c echo.Context) error {
		var req SignUpRequest
		if err := c.Bind(&req); err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, defaultErr)
		}
		if err := c.Validate(&req); err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, defaultErr)
		}

		_, err := h.AuthStore.SignUp(req.Email, req.Password, req.Name)
		if err != nil {
			h.Logger.Error("error signing up", "error", err)
			return echo.NewHTTPError(http.StatusInternalServerError, "There has been an unexpected error signing you in.")
		}

		return c.NoContent(http.StatusCreated)
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
			h.Logger.Error("error refreshing token", "error", err)
			return echo.NewHTTPError(http.StatusUnauthorized)
		}
		return c.JSON(http.StatusOK, session)
	}
}
