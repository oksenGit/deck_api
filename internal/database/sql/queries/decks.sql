-- name: CreateDeck :one
INSERT INTO decks (id, shuffled, remaining) VALUES ($1, $2, $3) RETURNING *;

-- name: GetDeck :one
SELECT * FROM decks WHERE id = $1;