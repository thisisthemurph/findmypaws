package model

import (
	"encoding/json"
	"time"
)

type User struct {
	ID        string          `db:"id"`
	Data      json.RawMessage `db:"data"`
	CreatedAt time.Time       `db:"created_at"`
	UpdatedAt time.Time       `db:"updated_at"`
}

type AnonymousUser struct {
	ID        string    `db:"id"`
	Name      string    `db:"name"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}
