-- name: CreateDeck :one
INSERT INTO decks (id, shuffled, remaining) VALUES ($1, $2, $3) RETURNING *;

