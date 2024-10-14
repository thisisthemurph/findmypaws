package store

import (
	"database/sql"
	"errors"
	"fmt"
	"paws/internal/types"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	sbauth "github.com/supabase-community/gotrue-go/types"
	"github.com/supabase-community/supabase-go"
)

var (
	ErrUserNotFound        = errors.New("user not found")
	ErrUserWithEmailExists = errors.New("email already exists")
)

type PostgresUserStore struct {
	*sqlx.DB
	Supabase *supabase.Client
}

func NewPostgresUserStore(db *sqlx.DB, sb *supabase.Client) *PostgresUserStore {
	return &PostgresUserStore{db, sb}
}

func (s *PostgresUserStore) SignUp(email, password, name string) (types.User, error) {
	// Step 1: Try to sign up or fetch the existing user.
	user, existing, err := s.getOrSignupSupabaseUser(email, password, name)
	if err != nil {
		return types.User{}, err
	}

	// Step 2: Check if the user already has a profile.
	// We only need to do this if there was an existing record in the auth.users table.
	if existing {
		stmt := `select exists(select 1 from profiles where id = $1);`
		var profileExists bool
		if err := s.Get(&profileExists, stmt, user.ID); err != nil {
			return types.User{}, fmt.Errorf("failed to check if profile exists: %w", err)
		}
		if profileExists {
			return types.User{}, ErrUserWithEmailExists
		}
	}

	// Step 3: Create a profile for the new user.
	stmt := `insert into profiles (id, name) values ($1, $2);`
	if _, err := s.Exec(stmt, user.ID, name); err != nil {
		return types.User{}, err
	}

	return s.User(user.ID)
}

func (s *PostgresUserStore) User(id uuid.UUID) (types.User, error) {
	stmt := `
		select p.*, u.id as auth_id, u.email, u.created_at, u.last_sign_in_at
		from public.profiles p
		join auth.users u on u.id = p.id
		where p.id = $1;`

	var u types.User
	if err := s.Get(&u, stmt, id); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return types.User{}, ErrUserNotFound
		}
		return types.User{}, err
	}

	return u, nil
}

func (s *PostgresUserStore) UserByEmail(email string) (types.User, error) {
	stmt := `
		select p.*, u.id as auth_id, u.email, u.created_at, u.last_sign_in_at
		from public.profiles p
		join auth.users u on u.id = p.id
		where u.email = $1 limit 1;`

	var u types.User
	if err := s.Get(&u, stmt, email); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return types.User{}, ErrUserNotFound
		}
		return types.User{}, err
	}

	return u, nil
}

// getOrSignupSupabaseUser handles user signup or fetching an existing user via Supabase.
func (s *PostgresUserStore) getOrSignupSupabaseUser(email, password, name string) (*sbauth.User, bool, error) {
	existingUser, err := s.getSupabaseUserByEmail(email)
	if errors.Is(err, ErrUserNotFound) {
		// User not found, attempt to sign up via Supabase.
		resp, err := s.Supabase.Auth.Signup(sbauth.SignupRequest{
			Email:    email,
			Password: password,
			Data: map[string]interface{}{
				"name": name,
			},
		})
		if err != nil {
			return nil, false, fmt.Errorf("failed to sign up user in Supabase: %w", err)
		}
		return &resp.User, false, nil
	} else if err != nil {
		return nil, false, err
	}

	return existingUser, true, nil
}

func (s *PostgresUserStore) getSupabaseUserByEmail(email string) (*sbauth.User, error) {
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
