-- name: CreateDeckCard :one
INSERT INTO deck_cards (deck_id, card_code, "order") VALUES ($1, $2, $3) RETURNING *;

-- name: GetDeckRemainingCards :many
SELECT card_code FROM deck_cards 
WHERE deck_id = $1
AND drawn = false
ORDER BY "order" ASC;