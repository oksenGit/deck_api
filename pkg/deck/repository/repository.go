package repository

import (
	"context"
	"database/sql"

	"github.com/google/uuid"
	"github.com/oksenGit/deck_api/internal/database"
	deck "github.com/oksenGit/deck_api/internal/deck"
)

// Repository is a struct that holds the database connection.
type Repository struct {
	DB *database.Queries
}

// NewRepository creates a new instance of Repository.
func NewRepository(db *database.Queries) *Repository {
	return &Repository{
		DB: db,
	}
}

// CreateDeck creates a new deck in the database.
func (r *Repository) CreateDeck(ctx context.Context, deck *deck.Deck, tx *sql.Tx) (*database.Deck, error) {
	params := database.CreateDeckParams{
		ID:       uuid.MustParse(deck.ID),
		Shuffled: deck.Shuffled,
	}
	databaseDeck, _err := r.DB.WithTx(tx).CreateDeck(ctx, params)
	if _err != nil {
		return nil, _err
	}

	return &databaseDeck, nil

}

// CreateDeckCards creates new deck cards in the database.
func (r *Repository) CreateDeckCards(ctx context.Context, deckID uuid.UUID, cards []deck.Card, tx *sql.Tx) ([]database.DeckCard, error) {
	dbCards := make([]database.DeckCard, 0, len(cards))
	for index, card := range cards {
		params := database.CreateDeckCardParams{
			DeckID:   deckID,
			CardCode: card.Code,
			Order:    int32(index),
		}
		dbCard, err := r.DB.WithTx(tx).CreateDeckCard(ctx, params)
		dbCards = append(dbCards, dbCard)
		if err != nil {
			return nil, err
		}
	}

	return dbCards, nil
}

// Get Deck With Cards
func (r *Repository) GetDeck(ctx context.Context, deckID uuid.UUID) (*database.Deck, error) {
	deck, err := r.DB.GetDeck(ctx, deckID)
	if err != nil {
		return nil, err
	}

	return &deck, nil
}

func (r *Repository) GetDeckRemainingCards(ctx context.Context, deckID uuid.UUID, limit *int) ([]string, error) {
	dbLimit := sql.NullInt32{Valid: false}
	if limit != nil {
		dbLimit = sql.NullInt32{Valid: true, Int32: int32(*limit)}
	}
	params := database.GetDeckRemainingCardsParams{
		DeckID: deckID,
		Limit:  dbLimit,
	}
	cards, err := r.DB.GetDeckRemainingCards(ctx, params)
	if err != nil {
		return nil, err
	}

	return cards, nil
}

func (r *Repository) SetDeckCardsDrawn(ctx context.Context, deckID uuid.UUID, cardCodes []string) error {
	params := database.SetDeckCardsDrawnParams{
		DeckID:  deckID,
		Column2: cardCodes,
	}
	err := r.DB.SetDeckCardsDrawn(ctx, params)
	return err
}
