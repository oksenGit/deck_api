-- +goose Up

ALTER TABLE decks
DROP COLUMN remaining;

-- +goose Down

ALTER TABLE decks
ADD COLUMN remaining INTEGER NOT NULL;