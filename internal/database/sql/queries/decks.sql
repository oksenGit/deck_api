-- name: CreateDeck :one
INSERT INTO decks (id, shuffled) VALUES ($1, $2) RETURNING *;

-- name: GetDeck :one
SELECT * FROM decks WHERE id = $1;