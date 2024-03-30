-- +goose Up

CREATE TABLE IF NOT EXISTS deck_cards (
    deck_id UUID REFERENCES decks(id),
    card_code VARCHAR(2) NOT NULL,
    PRIMARY KEY (deck_id, card_code),
    drawn BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP
);

CREATE INDEX deck_cards_deck_id_idx ON deck_cards(deck_id);

-- +goose Down

DROP TABLE IF EXISTS deck_cards;