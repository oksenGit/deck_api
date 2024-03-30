package deck

import (
	"math/rand"
	"time"
	"github.com/google/uuid"
)

func NewDeck(shuffled bool, cards []string) *Deck {
	generatedCards := generateCards(cards)
	deck := &Deck{
		ID:        generateUUID(),
		Shuffled:  shuffled,
		Remaining: len(generatedCards),
		Cards:     generateCards(cards),
	}

	if shuffled {
		rand.NewSource(time.Now().UnixNano())
		rand.Shuffle(len(deck.Cards), func(i, j int) {
			deck.Cards[i], deck.Cards[j] = deck.Cards[j], deck.Cards[i]
		})
	}

	return deck
}

func generateUUID() string {
	return uuid.New().String()
}

func generateCards(cards []string) []Card {
	if len(cards) == 0 {
		cards = GetStandardDeck()
	}

	var deck []Card
	var cardsSet = make(map[string]bool)
	for _, code := range cards {
		card := decodeCardCode(code)
		if card != nil  && !cardsSet[card.Code] {
			deck = append(deck, *card)
			cardsSet[card.Code] = true
		}
	}

	return deck
}

func GetStandardDeck() []string {
	return []string{
		"AS", "2S", "3S", "4S", "5S", "6S", "7S", "8S", "9S", "TS", "JS", "QS", "KS",
		"AC", "2C", "3C", "4C", "5C", "6C", "7C", "8C", "9C", "TC", "JC", "QC", "KC",
		"AH", "2H", "3H", "4H", "5H", "6H", "7H", "8H", "9H", "TH", "JH", "QH", "KH",
		"AD", "2D", "3D", "4D", "5D", "6D", "7D", "8D", "9D", "TD", "JD", "QD", "KD",
	}
}

func decodeCardCode(code string) *Card {
	if len(code) != 2 {
		return nil
	}

	value := DecodeCardValue(code[:1])
	suit := DecodeCardSuit(code[1:])

	if value == "" || suit == "" {
		return nil
	}

	return &Card{
		Value: value,
		Suit:  suit,
		Code:  code,
	}
}

func DecodeCardValue(value string) string {
    cardValues := map[string]string{
        "A": "ACE",
        "2": "2",
        "3": "3",
        "4": "4",
        "5": "5",
        "6": "6",
        "7": "7",
        "8": "8",
        "9": "9",
        "T": "10",
        "J": "JACK",
        "Q": "QUEEN",
        "K": "KING",
    }

    return cardValues[value]
}

func DecodeCardSuit(suit string) string {
    cardSuits := map[string]string{
        "S": "SPADES",
        "C": "CLUBS",
        "H": "HEARTS",
        "D": "DIAMONDS",
    }

    return cardSuits[suit]
}