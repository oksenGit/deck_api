package deck

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGenerateStandardDeck(t *testing.T) {
	deck := GetStandardDeck()
	assert.Equal(t, 52, len(deck))
}

func TestNewDeck(t *testing.T) {
	t.Run("Create new shuffled deck", func(t *testing.T) {
		deck := NewDeck(true, []string{})
		assert.NotNil(t, deck)
		assert.Equal(t, 52, len(deck.Cards))

		assert.Equal(t, true, deck.Shuffled)
		shuffled := checkDeckShuffled(*deck, GetStandardDeck())
		assert.Equal(t, true, shuffled)
	})

	t.Run("Create new unshuffled deck", func(t *testing.T) {
		deck := NewDeck(false, []string{})
		assert.NotNil(t, deck)
		assert.Equal(t, 52, len(deck.Cards))

		assert.Equal(t, false, deck.Shuffled)
		shuffled := checkDeckShuffled(*deck, GetStandardDeck())
		assert.Equal(t, false, shuffled)
	})


	t.Run("Create new shuffled deck with custom cards", func(t *testing.T) {
		cards := []string{"AS", "KD", "QC", "JH", "TS"}
		deck := NewDeck(true, cards)
		assert.NotNil(t, deck)
		assert.Equal(t, 5, len(deck.Cards))
		assert.Equal(t, true, deck.Shuffled)
		// avoid checking shuffled deck as it may be the same as the input for small sets of cards
	})

	t.Run("Create new unshuffled deck with custom cards", func(t *testing.T) {
		cards := []string{"AS", "KD", "QC", "JH", "TS"}
		deck := NewDeck(false, cards)
		assert.NotNil(t, deck)
		assert.Equal(t, 5, len(deck.Cards))
		assert.Equal(t, false, deck.Shuffled)
		shuffled := checkDeckShuffled(*deck, cards)
		assert.Equal(t, false, shuffled)
	})
}

func TestDecodeCardCode(t *testing.T) {
	card := decodeCardCode("AS")
	assert.NotNil(t, card)
	assert.Equal(t, "ACE", card.Value)
	assert.Equal(t, "SPADES", card.Suit)
	assert.Equal(t, "AS", card.Code)
}

func checkDeckShuffled(deck Deck, cards []string) bool {
	shuffled := false
	for i, card := range deck.Cards {
		if card.Code != cards[i] {
			shuffled = true
			break
		}
	}
	return shuffled
}
