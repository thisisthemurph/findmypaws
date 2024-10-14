package routes

import (
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"log/slog"
	"paws/internal/store"
)

type RouteMaker interface {
	MakeRoutes(e *echo.Group)
}

type Router struct {
	*echo.Echo
}

func NewRouter(s *store.PostgresStore, logger *slog.Logger) *Router {
	r := &Router{
		Echo: echo.New(),
	}

	r.Validator = &CustomValidator{validator: validator.New()}
	baseGroup := r.Group("/api/v1")
	for _, h := range r.getRouteHandlers(s, logger) {
		h.MakeRoutes(baseGroup)
	}

	return r
}

func (r *Router) getRouteHandlers(s *store.PostgresStore, logger *slog.Logger) []RouteMaker {
	return []RouteMaker{
		NewPetsHandler(s, logger),
	}
}

type CustomValidator struct {
	validator *validator.Validate
}

func (cv *CustomValidator) Validate(i interface{}) error {
	return cv.validator.Struct(i)
}
