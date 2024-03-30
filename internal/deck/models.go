package deck

type Card struct {
	Value string
	Suit  string
	Code  string
}

// Deck represents a deck of cards.
type Deck struct {
	ID        string
	Shuffled  bool  
	Remaining int   
	Cards     []Card
}
