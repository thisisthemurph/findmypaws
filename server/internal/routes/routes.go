package routes

import (
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"log/slog"
	"net/http"
	"paws/internal/store"
)

type RouteMaker interface {
	MakeRoutes(e *echo.Group)
}

type Router struct {
	*echo.Echo
}

func NewRouter(s *store.PostgresStore, clientBaseURL string, logger *slog.Logger) *Router {
	r := &Router{
		Echo: echo.New(),
	}

	r.Validator = &CustomValidator{validator: validator.New()}

	userMiddleware := NewUserMiddleware(s, logger)
	baseGroup := r.Group("/api/v1")
	baseGroup.Use(userMiddleware.WithUserInContext)

	for _, h := range r.getRouteHandlers(s, logger) {
		h.MakeRoutes(baseGroup)
	}

	configureCORS(r, clientBaseURL)

	return r
}

func (r *Router) getRouteHandlers(s *store.PostgresStore, logger *slog.Logger) []RouteMaker {
	return []RouteMaker{
		NewAuthHandler(s, logger),
		NewUserHandler(s, logger),
		NewPetsHandler(s, logger),
	}
}

func configureCORS(r *Router, clientBaseURL string) {
	allowOrigins := []string{clientBaseURL}
	r.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins:     allowOrigins,
		AllowHeaders:     []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept, echo.HeaderAuthorization},
		AllowMethods:     []string{http.MethodGet, http.MethodPost, http.MethodPut, http.MethodPatch, http.MethodDelete},
		AllowCredentials: true,
	}))
}

type CustomValidator struct {
	validator *validator.Validate
}

func (cv *CustomValidator) Validate(i interface{}) error {
	return cv.validator.Struct(i)
}
