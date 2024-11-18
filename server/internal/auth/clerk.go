package auth

import (
	"context"
	"github.com/clerk/clerk-sdk-go/v2"
	clerkhttp "github.com/clerk/clerk-sdk-go/v2/http"
	"github.com/clerk/clerk-sdk-go/v2/user"
	"net/http"
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

// WithClerkUserInContextMiddleware adds the authenticated Clerk user to the context.
func WithClerkUserInContextMiddleware(next http.HandlerFunc) http.HandlerFunc {
	clerkHandler := clerkhttp.WithHeaderAuthorization()(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		claims, ok := clerk.SessionClaimsFromContext(r.Context())
		if !ok {
			next.ServeHTTP(w, r)
			return
		}

		usr, err := user.Get(r.Context(), claims.Subject)
		if err != nil || usr == nil {
			next.ServeHTTP(w, r)
			return
		}

		authedUser := NewClerkAuthedUser(*usr)
		ctx := context.WithValue(r.Context(), UserContextKey, authedUser)
		next.ServeHTTP(w, r.WithContext(ctx))
	}))

	return func(w http.ResponseWriter, r *http.Request) {
		clerkHandler.ServeHTTP(w, r)
	}
}

// GetUserFromContext retrieves the authenticated user from the context.
func GetUserFromContext(ctx context.Context) ClerkAuthedUser {
	u, ok := ctx.Value(UserContextKey).(ClerkAuthedUser)
	if !ok {
		return ClerkAuthedUser{}
	}
	return u
}
