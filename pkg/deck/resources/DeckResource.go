package resources

import (
	"github.com/oksenGit/deck_api/internal/database"
	"github.com/oksenGit/deck_api/internal/deck"
)

type deckResource struct {
	ID        string `json:"deck_id"`
	Shuffled  bool   `json:"shuffled"`
	Remaining int32  `json:"remaining"`
}

type cardResource struct {
	Value string `json:"value"`
	Suit  string `json:"suit"`
	Code  string `json:"code"`
}

type cardsResource struct {
	Cards []cardResource `json:"cards"`
}

type deckWithCardsResource struct {
	deckResource
	Cards []cardResource `json:"cards"`
}

func CreateDeckResource(dbDeck *database.Deck, remaining int32) deckResource {
	return deckResource{
		ID:        (*dbDeck).ID.String(),
		Shuffled:  (*dbDeck).Shuffled,
		Remaining: remaining,
	}
}

func GetDeckWithRemainingCards(dbDeck *database.Deck, cards []string) deckWithCardsResource {
	deckWithCardsResource := deckWithCardsResource{
		deckResource: deckResource{
			ID:        (*dbDeck).ID.String(),
			Shuffled:  (*dbDeck).Shuffled,
			Remaining: int32(len(cards)),
		},
		Cards: make([]cardResource, 0, len(cards)),
	}

	for _, cardCode := range cards {
		deckWithCardsResource.Cards = append(deckWithCardsResource.Cards, cardResource{
			Value: deck.DecodeCardValue(cardCode[:1]),
			Suit:  deck.DecodeCardSuit(cardCode[1:]),
			Code:  cardCode,
		})
	}

	return deckWithCardsResource
}

func DrawCardsResource(cards []string) cardsResource {
	cardsResource := cardsResource{
		Cards: make([]cardResource, 0, len(cards)),
	}

	for _, cardCode := range cards {
		cardsResource.Cards = append(cardsResource.Cards, cardResource{
			Value: deck.DecodeCardValue(cardCode[:1]),
			Suit:  deck.DecodeCardSuit(cardCode[1:]),
			Code:  cardCode,
		})
	}

	return cardsResource
}