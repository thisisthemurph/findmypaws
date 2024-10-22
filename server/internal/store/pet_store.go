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
		select id, name, tags, dob, avatar_uri, created_at, updated_at, coalesce(type, $2) as type
    	from pets p where id = $1;`

	var p types.Pet
	if err := s.QueryRow(stmt, id, types.PetTypeUnknown).Scan(&p.ID, &p.Name, &p.Tags, &p.DOB, &p.AvatarURI, &p.CreatedAt, &p.UpdatedAt, &p.Type); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return p, ErrPetNotFound
		}
		return p, err
	}

	return p, nil
}

func (s PostgresPetStore) Pets(userID uuid.UUID) ([]types.Pet, error) {
	stmt := `
		select id, name, tags, dob, avatar_uri, created_at, updated_at, coalesce(type, $2) as type
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

func (s PostgresPetStore) Update(p *types.Pet, userID uuid.UUID) error {
	stmt := `
		update pets set name = $1, tags = $2, dob = $3, type = $4
		where id = $5 and user_id = $6
		returning *, coalesce(type, $7) as type;`

	if err := s.Get(p, stmt, p.Name, p.Tags, p.DOB, p.Type, p.ID, userID, types.PetTypeUnknown); err != nil {
		return err
	}
	return nil
}

func (s PostgresPetStore) UpdateAvatar(avatarURI string, petID, userID uuid.UUID) error {
	stmt := `
		update pets set avatar_uri = $1
		where id = $2 and user_id = $3;`

	_, err := s.Exec(stmt, avatarURI, petID, userID)
	if err != nil {
		return err
	}
	return nil
}

func (s PostgresPetStore) Delete(id, userID uuid.UUID) error {
	stmt := `delete from pets where id = $1 and user_id = $2;`
	if _, err := s.Exec(stmt, id, userID); err != nil {
		return err
	}
	return nil
}
