package routes

import (
	"fmt"
	"github.com/google/uuid"
	"log/slog"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"paws/internal/store"
)

type UserMiddleware struct {
	userStore store.UserStore
	logger    *slog.Logger
	jwtSecret []byte
}

func NewUserMiddleware(
	s *store.PostgresStore,
	logger *slog.Logger,
) *UserMiddleware {
	return &UserMiddleware{
		userStore: s.UserStore,
		logger:    logger,
		jwtSecret: []byte(os.Getenv("SUPABASE_JWT_SECRET")),
	}
}

type AppMetadata struct {
	Provider  string   `json:"provider"`
	Providers []string `json:"providers"`
}

type UserMetadata struct {
	Email         string `json:"email"`
	EmailVerified bool   `json:"email_verified"`
	Name          string `json:"name"`
	PhoneVerified bool   `json:"phone_verified"`
	Subscriber    string `json:"sub"`
}

type AuthenticationMethodReference struct {
	Method    string `json:"method"`
	Timestamp int64  `json:"timestamp"`
}

type UserClaims struct {
	jwt.RegisteredClaims
	Email                        string                          `json:"email"`
	Phone                        string                          `json:"phone"`
	AppMetadata                  AppMetadata                     `json:"app_metadata"`
	UserMetadata                 UserMetadata                    `json:"user_metadata"`
	Role                         string                          `json:"role"`
	AuthenticationLevelAssurance string                          `json:"aal"`
	AuthenticationMethods        []AuthenticationMethodReference `json:"amr"`
	SessionID                    string                          `json:"session_id"`
	IsAnonymous                  bool                            `json:"is_anonymous"`
}

type UserSession struct {
	ID       uuid.UUID
	Email    string
	Name     string
	LoggedIn bool
}

// WithUserInContext middleware parses any present JWT and persists the content in the context if
// the JWT is valid and has not expired. A valid and unexpired JWT indicates an authenticated user.
func (m *UserMiddleware) WithUserInContext(next echo.HandlerFunc) echo.HandlerFunc {
	logger := m.logger.With("middleware", "WithUserInContext")
	return func(c echo.Context) error {
		token := c.Request().Header.Get("Authorization")
		if token == "" || !strings.HasPrefix(token, "Bearer ") {
			return next(c)
		}

		jwtToken, err := jwt.ParseWithClaims(token[7:], &UserClaims{}, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				logger.Debug("unexpected signing method when parsing JWT", "alg", token.Header["alg"])
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return m.jwtSecret, nil
		})

		if err != nil {
			logger.Error("error parsing claims from JWT", "error", err)
			return next(c)
		}

		if claims, ok := jwtToken.Claims.(*UserClaims); ok && jwtToken.Valid {
			if claims.ExpiresAt.After(time.Now().UTC()) {
				userID, err := uuid.Parse(claims.Subject)
				if err != nil {
					logger.Error("error parsing user ID from JWT", "error", err, "ID", claims.Subject)
					return next(c)
				}

				user := UserSession{
					ID:       userID,
					Email:    claims.Email,
					Name:     claims.UserMetadata.Name,
					LoggedIn: true,
				}

				c.Set("user", user)
				c.Set("user-claims", claims)
				return next(c)
			}
		}

		return next(c)
	}
}

// WithAuthedUser middleware prevents access to protected routes.
func (m *UserMiddleware) WithAuthedUser(next echo.HandlerFunc) echo.HandlerFunc {
	logger := m.logger.With("middleware", "WithAuthedUser")
	return func(c echo.Context) error {
		user := CurrentUser(c)
		if !user.LoggedIn {
			logger.Debug("user not logged in", "user", user)
			return echo.NewHTTPError(http.StatusUnauthorized)
		}
		return next(c)
	}
}
