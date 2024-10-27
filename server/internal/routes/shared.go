package routes

import (
	"github.com/clerk/clerk-sdk-go/v2"
	"github.com/labstack/echo/v4"
	"paws/internal/auth"
)

// clerkUser returns the current authenticated auth.ClerkAuthedUser from the context if present.
// If not present a default auth.ClerkAuthedUser is returned where Authenticated will be false.
func clerkUser(c echo.Context) auth.ClerkAuthedUser {
	if u, ok := c.Get(string(auth.UserContextKey)).(clerk.User); ok {
		return auth.NewClerkAuthedUser(u)
	}
	return auth.ClerkAuthedUser{}
}
