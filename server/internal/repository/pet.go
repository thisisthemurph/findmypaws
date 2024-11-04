package repository

import (
	"database/sql"
	"errors"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"paws/internal/types"
)

type PetRepository interface {
	Get(id uuid.UUID) (types.Pet, error)
	List(userID string) ([]types.Pet, error)
	Create(pet *types.Pet) error
	Update(pet *types.Pet) error
	Delete(id uuid.UUID) error
}

type postgresPetRepository struct {
	db *sqlx.DB
}

func NewPetRepository(db *sqlx.DB) PetRepository {
	return &postgresPetRepository{
		db: db,
	}
}

func (r *postgresPetRepository) Get(id uuid.UUID) (types.Pet, error) {
	stmt := `
		select id, user_id, name, tags, dob, avatar_uri, blurb, created_at, updated_at, coalesce(type, $2) as type
    	from pets
    	where id = $1;`

	var p types.Pet
	if err := r.db.Get(&p, stmt, id, types.PetTypeUnknown); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return p, ErrNotFound
		}
		return p, err
	}
	return p, nil
}

func (r *postgresPetRepository) List(userID string) ([]types.Pet, error) {
	stmt := `
		select id, user_id, name, tags, dob, avatar_uri, blurb, created_at, updated_at, coalesce(type, $2) as type
		from pets
		where user_id = $1;`

	pp := make([]types.Pet, 0)
	if err := r.db.Select(&pp, stmt, userID, types.PetTypeUnknown); err != nil {
		return pp, err
	}
	return pp, nil
}

func (r *postgresPetRepository) Create(p *types.Pet) error {
	stmt := `
		insert into pets (user_id, name, type, tags, dob) 
		values ($1, $2, $3, $4, $5) 
		returning id, created_at, updated_at;`

	if err := r.db.Get(p, stmt, p.UserID, p.Name, p.Type, p.Tags, p.DOB); err != nil {
		return err
	}
	return nil
}

func (r *postgresPetRepository) Update(p *types.Pet) error {
	stmt := `
		update pets set 
		    name = $1,
		    tags = $2,
		    dob = $3,
		    type = $4,
		    blurb = $5,
		    avatar_uri = $6
		where id = $7
		returning *, coalesce(type, $8) as type;`

	if p.UserID == "" {
		return errors.New("userId is required")
	}

	if err := r.db.Get(p, stmt, p.Name, p.Tags, p.DOB, p.Type, p.Blurb, p.AvatarURI, p.ID, types.PetTypeUnknown); err != nil {
		return err
	}
	return nil
}

func (r *postgresPetRepository) Delete(id uuid.UUID) error {
	stmt := `delete from pets where id = $1;`
	if _, err := r.db.Exec(stmt, id); err != nil {
		return err
	}
	return nil
}
