package handler

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/go-chi/chi"
	"github.com/google/uuid"
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
	shuffled := shuffledStr == "true"

	cards := r.URL.Query().Get("cards")

	cardsList := []string{}
	if cards != "" {
		cardsList = strings.Split(cards, ",")
	}

	tx, err := db.DB.Begin()

	if err != nil {
		helpers.RespondWithError(w, http.StatusInternalServerError, "Server Error 001")
		return
	}

	defer func() {
		if err != nil {
			tx.Rollback()
			helpers.RespondWithError(w, http.StatusInternalServerError, "Server Error 002")
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
		helpers.RespondWithError(w, http.StatusInternalServerError, "Server Error 003")
	}

	helpers.RespondWithJSON(w, http.StatusCreated, resources.CreateDeckResource(dbDeck, int32(deckObj.Remaining)))
}

func (h *Handler) GetDeckWithRemainingCards(w http.ResponseWriter, r *http.Request) {
	deckID := chi.URLParam(r, "deckID")

	if deckID == "" {
		helpers.RespondWithError(w, http.StatusBadRequest, "deck_id is required")
		return
	}

	deckUUID, err := uuid.Parse(deckID)

	if err != nil {
		helpers.RespondWithError(w, http.StatusBadRequest, "deck_id is invalid")
		return
	}

	deck, err := h.Repo.GetDeck(r.Context(), deckUUID)

	if err != nil {
		helpers.RespondWithError(w, http.StatusInternalServerError, "Server Error 004")
		return
	}

	if deck == nil {
		helpers.RespondWithError(w, http.StatusNotFound, "Deck not found")
		return
	}

	deckCards, err := h.Repo.GetDeckRemainingCards(r.Context(), deckUUID, nil)

	if err != nil {
		helpers.RespondWithError(w, http.StatusInternalServerError, "Server Error 005")
		return
	}

	helpers.RespondWithJSON(w, http.StatusOK, resources.GetDeckWithRemainingCards(deck, deckCards))
}

func (h *Handler) DrawCards(w http.ResponseWriter, r *http.Request) {
	// get deck_id and count from request body
	type drawCardsRequest struct {
		DeckID string `json:"deck_id"`
		Count  int    `json:"count"`
	}
	var req drawCardsRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		helpers.RespondWithError(w, http.StatusBadRequest,  err.Error())
		return
	}

	deckID := req.DeckID
	count := req.Count

	if deckID == "" {
		helpers.RespondWithError(w, http.StatusBadRequest, "deck_id is required")
		return
	}

	if count <= 0 {
		helpers.RespondWithError(w, http.StatusBadRequest, "count must be greater than 0")
		return
	}

	deckUUID, err := uuid.Parse(deckID)

	if err != nil {
		helpers.RespondWithError(w, http.StatusBadRequest, "deck_id is invalid")
		return
	}

	deck, err := h.Repo.GetDeck(r.Context(), deckUUID)

	if err != nil {
		helpers.RespondWithError(w, http.StatusInternalServerError, "Server Error 006")
		return
	}

	if deck == nil {
		helpers.RespondWithError(w, http.StatusNotFound, "Deck not found")
		return
	}

	deckCards, err := h.Repo.GetDeckRemainingCards(r.Context(), deckUUID, &count)

	if err != nil {
		helpers.RespondWithError(w, http.StatusInternalServerError, "Server Error 007")
		return
	}

	err = h.Repo.SetDeckCardsDrawn(r.Context(), deckUUID, deckCards)

	if err != nil {
		helpers.RespondWithError(w, http.StatusInternalServerError, "Server Error 008")
		return
	}

	helpers.RespondWithJSON(w, http.StatusOK, resources.DrawCardsResource(deckCards))
}
