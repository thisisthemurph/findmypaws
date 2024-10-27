package auth

import (
	"context"
	"github.com/clerk/clerk-sdk-go/v2/user"
	"net/http"

	"github.com/clerk/clerk-sdk-go/v2"
	clerkhttp "github.com/clerk/clerk-sdk-go/v2/http"
	"github.com/labstack/echo/v4"
)

type ContextKey string

const UserContextKey ContextKey = "user"

// ClerkAuthedUser is a clerk.User with an additional bool Authenticated field.
type ClerkAuthedUser struct {
	clerk.User
	Authenticated bool
}

// NewClerkAuthedUser creates a ClerkAuthedUser with Authenticated set to true.
func NewClerkAuthedUser(u clerk.User) ClerkAuthedUser {
	return ClerkAuthedUser{
		User:          u,
		Authenticated: true,
	}
}

// WithEchoClerkMiddleware maps the Clerk authorization middleware for use with Echo.
func WithEchoClerkMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		handler := clerkhttp.WithHeaderAuthorization()(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			c.SetRequest(r)
			if err := next(c); err != nil {
				c.Error(err)
				return
			}
		}))

		res, req := c.Response(), c.Request()
		handler.ServeHTTP(res, req)

		if !c.Response().Committed {
			return c.JSON(http.StatusUnauthorized, map[string]string{
				"error": "Unauthorized or missing session",
			})
		}

		return nil
	}
}

// WithClerkUserInContextMiddleware middleware adds the clerk.User to the context if present.
func WithClerkUserInContextMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.Request().Context()
		claims, ok := clerk.SessionClaimsFromContext(ctx)
		if !ok {
			return next(c)
		}

		usr, err := user.Get(ctx, claims.Subject)
		if err != nil {
			return next(c)
		}
		if usr == nil {
			return next(c)
		}

		c.Set(string(UserContextKey), *usr)
		c.SetRequest(c.Request().WithContext(context.WithValue(c.Request().Context(), UserContextKey, *usr)))
		return next(c)
	}
}
