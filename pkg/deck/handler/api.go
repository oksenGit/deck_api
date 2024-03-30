package handler

import (
	"net/http"
	"strings"

	"github.com/oksenGit/deck_api/internal/deck"
	"github.com/oksenGit/deck_api/internal/helpers"
	"github.com/oksenGit/deck_api/pkg/deck/db"
	"github.com/oksenGit/deck_api/pkg/deck/repository"
	"github.com/oksenGit/deck_api/pkg/deck/resources"
)

type Handler struct {
	Repo *repository.Repository
}

func NewHandler(repo *repository.Repository) *Handler {
	return &Handler{
		Repo: repo,
	}
}

func (h *Handler) CreateDeck(w http.ResponseWriter, r *http.Request) {
	shuffledStr := r.URL.Query().Get("shuffled")
	shuffled := shuffledStr == "" || shuffledStr == "true"

	cards := r.URL.Query().Get("cards")

	cardsList := []string{}
	if cards != "" {
		cardsList = strings.Split(cards, ",")
	}

	tx, err := db.DB.Begin()

	if err != nil {
		http.Error(w, "Server Error 001", http.StatusInternalServerError)
		return
	}

	defer func() {
		if err != nil {
			tx.Rollback()
			http.Error(w, "Server Error 002", http.StatusInternalServerError)
		}
	}()

	deckObj := deck.NewDeck(shuffled, cardsList)

	dbDeck, err := h.Repo.CreateDeck(r.Context(), deckObj, tx)
	if err != nil {
		return
	}

	_, err = h.Repo.CreateDeckCards(r.Context(), dbDeck.ID, deckObj.Cards, tx)

	if err != nil {
		return
	}

	err = tx.Commit()
	if err != nil {
		http.Error(w, "Server Error 003", http.StatusInternalServerError)
	}

	helpers.RespondWithJSON(w, http.StatusCreated, resources.CreateDeckResource(*dbDeck))
}
