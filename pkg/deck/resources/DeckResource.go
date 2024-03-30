package resources

import "github.com/oksenGit/deck_api/internal/database"

type deckResource struct {
	ID        string `json:"deck_id"`
	Shuffled  bool   `json:"shuffled"`
	Remaining int32  `json:"remaining"`
}

func CreateDeckResource(dbDeck database.Deck, remaining int32) deckResource {
	return deckResource{
		ID:        dbDeck.ID.String(),
		Shuffled:  dbDeck.Shuffled,
		Remaining: remaining,
	}
}
