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
		select id, name, tags, dob, created_at, updated_at, coalesce(type, $2) as type
    	from pets p where id = $1;`

	var p types.Pet
	if err := s.QueryRow(stmt, id, types.PetTypeUnknown).Scan(&p.ID, &p.Name, &p.Tags, &p.DOB, &p.CreatedAt, &p.UpdatedAt, &p.Type); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return p, ErrPetNotFound
		}
		return p, err
	}

	return p, nil
}

func (s PostgresPetStore) Pets(userID uuid.UUID) ([]types.Pet, error) {
	return make([]types.Pet, 0), nil
}

func (s PostgresPetStore) Create(p *types.Pet) error {
	stmt := `
		insert into pets (name, tags, dob) 
		values ($1, $2, $3) 
		returning id, created_at, updated_at;`

	if err := s.Get(p, stmt, p.Name, p.Tags, p.DOB); err != nil {
		return err
	}
	return nil
}

func (s PostgresPetStore) Update(p *types.Pet) error {
	stmt := `
		update pets set
			name = $1,
		    tags = $2,
		    dob = $3,
		    type = $4
		where id = $5
		returning *, coalesce(type, $6) as type;`

	if err := s.Get(p, stmt, p.Name, p.Tags, p.DOB, p.Type, p.ID, types.PetTypeUnknown); err != nil {
		return err
	}
	return nil
}
