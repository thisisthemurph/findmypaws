package store

import (
	"database/sql"
	"errors"
	"fmt"
	"log/slog"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	sbauth "github.com/supabase-community/gotrue-go/types"
	"github.com/supabase-community/supabase-go"
	"paws/internal/types"
)

var (
	ErrInvalidCredentials error = errors.New("invalid credentials")
)

type PostgresAuthStore struct {
	*sqlx.DB
	Supabase  *supabase.Client
	Logger    *slog.Logger
	JWTSecret string
}

func NewPostgresAuthStore(
	db *sqlx.DB,
	sb *supabase.Client,
	jwtSecret string,
) *PostgresAuthStore {
	return &PostgresAuthStore{
		DB:        db,
		Supabase:  sb,
		JWTSecret: jwtSecret,
	}
}

func (s *PostgresAuthStore) LogIn(email, password string) (sbauth.Session, error) {
	session, err := s.Supabase.SignInWithEmailPassword(email, password)
	if err != nil {
		if sbErr, ok := NewSupabaseAuthError(err); ok {
			if sbErr.ErrorName == "invalid_grant" {
				return sbauth.Session{}, ErrInvalidCredentials
			}
		}
		return sbauth.Session{}, err
	}
	return session, nil
}

func (s *PostgresAuthStore) LogOut(token string) error {
	return s.Supabase.Auth.WithToken(token).Logout()
}

func (s *PostgresAuthStore) SignUp(email, password, name string) (uuid.UUID, error) {
	// Step 1: Try to sign up or fetch the existing user.
	user, existing, err := s.getOrSignupSupabaseUser(email, password, name)
	if err != nil {
		return uuid.Nil, err
	}

	// Step 2: Check if the user already has a profile.
	// We only need to do this if there was an existing record in the auth.users table.
	if existing {
		stmt := `select exists(select 1 from profiles where id = $1);`
		var profileExists bool
		if err := s.Get(&profileExists, stmt, user.ID); err != nil {
			return uuid.Nil, fmt.Errorf("failed to check if profile exists: %w", err)
		}
		if profileExists {
			return uuid.Nil, ErrUserWithEmailExists
		}
	}

	// Step 3: Create a profile for the new user.
	stmt := `insert into profiles (id, name) values ($1, $2);`
	if _, err := s.Exec(stmt, user.ID, name); err != nil {
		return uuid.Nil, err
	}

	return user.ID, nil
}

func (s *PostgresAuthStore) RefreshToken(refreshToken string) (sbauth.Session, error) {
	return s.Supabase.RefreshToken(refreshToken)
}

// getOrSignupSupabaseUser handles user signup or fetching an existing user via Supabase.
func (s *PostgresAuthStore) getOrSignupSupabaseUser(email, password, name string) (*sbauth.User, bool, error) {
	existingUser, err := s.getSupabaseUserByEmail(email)
	if errors.Is(err, ErrUserNotFound) {
		// User not found, attempt to sign up via Supabase.
		resp, signupErr := s.Supabase.Auth.Signup(sbauth.SignupRequest{
			Email:    email,
			Password: password,
			Data: map[string]interface{}{
				"name": name,
			},
		})
		if signupErr != nil {
			return nil, false, fmt.Errorf("failed to sign up user in Supabase: %w", signupErr)
		}
		return &resp.User, false, nil
	} else if err != nil {
		return nil, false, err
	}

	return existingUser, true, nil
}

func (s *PostgresAuthStore) getSupabaseUserByEmail(email string) (*sbauth.User, error) {
	stmt := `select id as auth_id, email, created_at from auth.users where email = $1;`
	var user types.UserAuth
	if err := s.Get(&user, stmt, email); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}
	return user.SupabaseUser(), nil
}
