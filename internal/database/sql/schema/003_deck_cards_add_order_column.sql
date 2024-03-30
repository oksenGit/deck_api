-- +goose Up

ALTER TABLE deck_cards
ADD COLUMN "order" INTEGER NOT NULL DEFAULT 0;

CREATE INDEX deck_cards_order_idx ON deck_cards ("order");

-- +goose Down

ALTER TABLE deck_cards
DROP COLUMN "order";

DROP INDEX deck_cards_order_idx;