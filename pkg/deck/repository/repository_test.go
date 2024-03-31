package repository

import (
	"context"
	"os"
	"testing"

	"github.com/google/uuid"
	"github.com/joho/godotenv"
	"github.com/oksenGit/deck_api/internal/database"
	"github.com/oksenGit/deck_api/internal/deck"
	"github.com/oksenGit/deck_api/pkg/deck/db"
)

func TestMain(m *testing.M) {

	err := godotenv.Load("../../../.env.test")
	if err != nil {
		panic("Error loading .env.test file")
	}

	db.Init()
	
	defer db.Close()

	code := m.Run()

	os.Exit(code)
}

// TestCreateDeck tests the CreateDeck method.
func TestCreateDeck(t *testing.T) {
	query := database.New(db.DB)
	repo := NewRepository(query)

	deckObj := deck.NewDeck(true, []string{})

	tx, _ := db.DB.Begin()
	defer tx.Rollback()

	dbDeck, err := repo.CreateDeck(context.Background(), deckObj, tx)

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if dbDeck.ID == uuid.Nil {
		t.Errorf("Expected deck ID to be set")
	}

	if dbDeck.Shuffled != deckObj.Shuffled {
		t.Errorf("Expected deck shuffled to be %v, got %v", deckObj.Shuffled, dbDeck.Shuffled)
	}

	dbCards, err := repo.CreateDeckCards(context.Background(), dbDeck.ID, deckObj.Cards, tx)

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if len(dbCards) != len(deckObj.Cards) {
		t.Errorf("Expected %v deck cards, got %v", len(deckObj.Cards), len(dbCards))
	}
}

func TestGetDeckWithRemainingCards(t *testing.T) {
	query := database.New(db.DB)
	repo := NewRepository(query)

	deckObj := deck.NewDeck(true, []string{})

	tx, _ := db.DB.Begin()
	defer tx.Rollback()

	dbDeck, _ := repo.CreateDeck(context.Background(), deckObj, tx)

	repo.CreateDeckCards(context.Background(), dbDeck.ID, deckObj.Cards, tx)
	tx.Commit()

	deck, err := repo.GetDeck(context.Background(), dbDeck.ID)

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if deck.ID != dbDeck.ID {
		t.Errorf("Expected deck ID to be %v, got %v", dbDeck.ID, deck.ID)
	}

	if deck.Shuffled != dbDeck.Shuffled {
		t.Errorf("Expected deck shuffled to be %v, got %v", dbDeck.Shuffled, deck.Shuffled)
	}

	remainingCards, err := repo.GetDeckRemainingCards(context.Background(), dbDeck.ID, nil)

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if len(remainingCards) != len(deckObj.Cards) {
		t.Errorf("Expected %v remaining cards, got %v", len(deckObj.Cards), len(remainingCards))
	}

	for i, card := range remainingCards {
		if card != deckObj.Cards[i].Code {
			t.Errorf("Expected card %v to be %v, got %v", i, deckObj.Cards[i], card)
		}
	}

	// make int pointer with default value of 5
	limit := new(int)
	*limit = 5
	remainingCards, err = repo.GetDeckRemainingCards(context.Background(), dbDeck.ID, limit)

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if len(remainingCards) != *limit {
		t.Errorf("Expected %v remaining cards, got %v", *limit, len(remainingCards))
	}
}