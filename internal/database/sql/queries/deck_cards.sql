-- name: CreateDeckCard :one
INSERT INTO deck_cards (deck_id, card_code) VALUES ($1, $2) RETURNING *;