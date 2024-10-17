package routes

import (
	"log/slog"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"paws/internal/store"
)

type RouteMaker interface {
	MakeRoutes(e *echo.Group)
}

type Router struct {
	*echo.Echo
}

func NewRouter(s *store.PostgresStore, clientBaseURL string, logger *slog.Logger) *Router {
	e := echo.New()
	e.Validator = NewCustomValidator()

	userMiddleware := NewUserMiddleware(s, logger)
	baseGroup := e.Group("/api/v1")
	baseGroup.Use(userMiddleware.WithUserInContext)

	for _, h := range getRouteHandlers(s, logger) {
		h.MakeRoutes(baseGroup)
	}

	configureCORS(e, clientBaseURL)

	return &Router{
		Echo: e,
	}
}

func getRouteHandlers(s *store.PostgresStore, logger *slog.Logger) []RouteMaker {
	return []RouteMaker{
		NewAuthHandler(s, logger),
		NewUserHandler(s, logger),
		NewPetsHandler(s, logger),
	}
}

func configureCORS(e *echo.Echo, clientBaseURL string) {
	allowOrigins := []string{clientBaseURL}
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins:     allowOrigins,
		AllowHeaders:     []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept, echo.HeaderAuthorization},
		AllowMethods:     []string{http.MethodGet, http.MethodPost, http.MethodPut, http.MethodPatch, http.MethodDelete},
		AllowCredentials: true,
	}))
}

type RequestValidator struct {
	validator *validator.Validate
}

func NewCustomValidator() *RequestValidator {
	return &RequestValidator{validator: validator.New()}
}

func (cv *RequestValidator) Validate(i interface{}) error {
	return cv.validator.Struct(i)
}
