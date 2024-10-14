package types

import (
	"github.com/google/uuid"
	sbauth "github.com/supabase-community/gotrue-go/types"
	"time"
)

// UserAuth represents a subset of the data in the auth.users table.
type UserAuth struct {
	AuthID       uuid.UUID  `json:"-" db:"auth_id"`
	Email        string     `json:"email" db:"email"`
	CreatedAt    time.Time  `json:"created_at" db:"created_at"`
	LastSignInAt *time.Time `json:"last_sign_in_at" db:"last_sign_in_at"`
}

func (u UserAuth) SupabaseUser() *sbauth.User {
	return &sbauth.User{
		ID:           u.AuthID,
		Email:        u.Email,
		LastSignInAt: u.LastSignInAt,
		CreatedAt:    u.CreatedAt,
	}
}

// UserProfile represents the columns in the public.profiles table.
type UserProfile struct {
	ID        uuid.UUID `json:"id" db:"id"`
	Name      string    `json:"name" db:"name"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

// User represents an amalgamation of the UserAuth and UserProfile columns.
type User struct {
	UserAuth
	UserProfile
}

func (u User) SupabaseUser() *sbauth.User {
	return &sbauth.User{
		ID:           u.AuthID,
		Email:        u.Email,
		LastSignInAt: u.LastSignInAt,
		CreatedAt:    u.CreatedAt,
		UpdatedAt:    u.UpdatedAt,
	}
}
