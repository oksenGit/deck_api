-- +goose Up

CREATE TABLE IF NOT EXISTS decks (
    id UUID PRIMARY KEY,
    shuffled BOOLEAN NOT NULL,
    remaining INTEGER NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP
);

-- +goose Down

DROP TABLE IF EXISTS decks;