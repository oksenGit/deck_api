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

	if dbDeck.Remaining != 52 {
		t.Errorf("Expected deck remaining to be %v, got %v", len(deckObj.Cards), dbDeck.Remaining)
	}

	dbCards, err := repo.CreateDeckCards(context.Background(), dbDeck.ID, deckObj.Cards, tx)

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if len(dbCards) != len(deckObj.Cards) {
		t.Errorf("Expected %v deck cards, got %v", len(deckObj.Cards), len(dbCards))
	}
}
