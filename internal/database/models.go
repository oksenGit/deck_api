// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.25.0

package database

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
)

type Deck struct {
	ID        uuid.UUID
	Shuffled  bool
	Remaining int32
	CreatedAt time.Time
	UpdatedAt sql.NullTime
}

type DeckCard struct {
	DeckID    uuid.UUID
	CardCode  string
	Drawn     bool
	CreatedAt time.Time
	UpdatedAt sql.NullTime
	Order     int32
}
