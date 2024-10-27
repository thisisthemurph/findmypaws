package store

import (
	"database/sql"
	"errors"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"paws/internal/types"
)

var ErrPetNotFound = errors.New("pet not found")

type PostgresPetStore struct {
	*sqlx.DB
}

func NewPostgresPetStore(db *sqlx.DB) *PostgresPetStore {
	return &PostgresPetStore{db}
}

func (s PostgresPetStore) Pet(id uuid.UUID) (types.Pet, error) {
	stmt := `
		select id, user_id, name, tags, dob, avatar_uri, blurb, created_at, updated_at, coalesce(type, $2) as type
    	from pets
    	where id = $1;`

	var p types.Pet
	if err := s.Get(&p, stmt, id, types.PetTypeUnknown); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return p, ErrPetNotFound
		}
		return p, err
	}
	return p, nil
}

func (s PostgresPetStore) Pets(userID string) ([]types.Pet, error) {
	stmt := `
		select id, user_id, name, tags, dob, avatar_uri, blurb, created_at, updated_at, coalesce(type, $2) as type
		from pets
		where user_id = $1;`

	pp := make([]types.Pet, 0)
	if err := s.Select(&pp, stmt, userID, types.PetTypeUnknown); err != nil {
		return pp, err
	}
	return pp, nil
}

func (s PostgresPetStore) Create(p *types.Pet) error {
	stmt := `
		insert into pets (user_id, name, type, tags, dob) 
		values ($1, $2, $3, $4, $5) 
		returning id, created_at, updated_at;`

	if err := s.Get(p, stmt, p.UserID, p.Name, p.Type, p.Tags, p.DOB); err != nil {
		return err
	}
	return nil
}

func (s PostgresPetStore) Update(p *types.Pet, userID string) error {
	stmt := `
		update pets set name = $1, tags = $2, dob = $3, type = $4, blurb = $5
		where id = $6 and user_id = $7
		returning *, coalesce(type, $8) as type;`

	if err := s.Get(p, stmt, p.Name, p.Tags, p.DOB, p.Type, p.Blurb, p.ID, userID, types.PetTypeUnknown); err != nil {
		return err
	}
	return nil
}

func (s PostgresPetStore) UpdateAvatar(avatarURI string, petID uuid.UUID, userID string) error {
	stmt := `
		update pets set avatar_uri = $1
		where id = $2 and user_id = $3;`

	_, err := s.Exec(stmt, avatarURI, petID, userID)
	if err != nil {
		return err
	}
	return nil
}

func (s PostgresPetStore) Delete(id uuid.UUID, userID string) error {
	stmt := `delete from pets where id = $1 and user_id = $2;`
	if _, err := s.Exec(stmt, id, userID); err != nil {
		return err
	}
	return nil
}
