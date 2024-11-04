package routes

import (
	"log/slog"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"paws/internal/auth"
	"paws/internal/repository"
)

type RouteMaker interface {
	MakeRoutes(e *echo.Group)
}

type Router struct {
	*echo.Echo
}

func NewRouter(repos *repository.Repositories, clientBaseURL string, logger *slog.Logger) *Router {
	e := echo.New()
	e.Validator = NewCustomValidator()
	e.Static("/static", "./static")

	//e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	configureCORS(e, clientBaseURL)

	baseGroup := e.Group("/api/v1")
	baseGroup.Use(auth.WithEchoClerkMiddleware)
	baseGroup.Use(auth.WithClerkUserInContextMiddleware)

	for _, h := range getRouteHandlers(repos, logger) {
		h.MakeRoutes(baseGroup)
	}

	return &Router{
		Echo: e,
	}
}

func getRouteHandlers(repos *repository.Repositories, logger *slog.Logger) []RouteMaker {
	return []RouteMaker{
		NewPetsHandler(repos.NotificationRepository, repos.PetRepository, logger),
		NewUsersHandler(repos.NotificationRepository, repos.PetRepository, logger),
	}
}

func configureCORS(e *echo.Echo, clientBaseURL string) {
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{clientBaseURL},
		AllowHeaders: []string{
			echo.HeaderAuthorization,
			echo.HeaderAccept,
			"Host",
			echo.HeaderOrigin,
			"Referer",
			"Sec-Fetch-Dest",
			"User-Agent",
			"X-Forwarded-Host",
			"X-Forwarded-Proto",
			echo.HeaderContentType,
		},
		AllowMethods: []string{
			http.MethodGet,
			http.MethodPost,
			http.MethodPut,
			http.MethodPatch,
			http.MethodDelete,
		},
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
