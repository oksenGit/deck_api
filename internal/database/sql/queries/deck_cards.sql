-- name: CreateDeckCard :one
INSERT INTO deck_cards (deck_id, card_code, "order") VALUES ($1, $2, $3) RETURNING *;

-- name: GetDeckRemainingCards :many
SELECT card_code FROM deck_cards 
WHERE deck_id = $1
AND drawn = false
ORDER BY "order" ASC
LIMIT sqlc.narg('limit')::int;

-- name: SetDeckCardsDrawn :exec
UPDATE deck_cards SET drawn = true WHERE deck_id = $1 AND card_code = ANY($2::text[]);