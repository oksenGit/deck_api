package handler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/go-chi/chi"
	"github.com/oksenGit/deck_api/internal/database"
	"github.com/oksenGit/deck_api/pkg/deck/db"
	"github.com/oksenGit/deck_api/pkg/deck/repository"
	"github.com/stretchr/testify/assert"
)

func TestMain(m *testing.M) {
	db.Init()
	defer db.Close()
	code := m.Run()
	os.Exit(code)
}

func TestCreateDeck(t *testing.T) {
	query := database.New(db.DB)
	repo := repository.NewRepository(query)
	h := NewHandler(repo)
	t.Run("CreateDefaultDeck", func(t *testing.T) {
		req, err := http.NewRequest("POST", "/v1/decks", nil)

		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()

		handler := http.HandlerFunc(h.CreateDeck)

		handler.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusCreated, rr.Code, "handler returned wrong status code")
		assert.NotNil(t, rr.Body, "handler returned no body")

		var response map[string]interface{}
		err = json.Unmarshal(rr.Body.Bytes(), &response)

		if err != nil {
			t.Fatal(err)
		}

		assert.NotNil(t, response["deck_id"], "deck_id is nil")
		assert.NotNil(t, response["shuffled"], "shuffled is nil")
		assert.Equal(t, response["shuffled"], false, "shuffled is not true")
		assert.NotNil(t, response["remaining"], "remaining is nil")
		assert.Equal(t, response["remaining"], float64(52), "remaining is not 52")

		// check that cards key does not exist
		_, ok := response["cards"]
		assert.False(t, ok, "cards key exists")
	})

	t.Run("CreateUnshuffledDeck", func(t *testing.T) {
		req, err := http.NewRequest("POST", "/v1/decks?shuffled=false", nil)

		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()

		handler := http.HandlerFunc(h.CreateDeck)

		handler.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusCreated, rr.Code, "handler returned wrong status code")
		assert.NotNil(t, rr.Body, "handler returned no body")

		var response map[string]interface{}
		err = json.Unmarshal(rr.Body.Bytes(), &response)

		if err != nil {
			t.Fatal(err)
		}

		assert.NotNil(t, response["deck_id"], "deck_id is nil")
		assert.NotNil(t, response["shuffled"], "shuffled is nil")
		assert.Equal(t, response["shuffled"], false, "shuffled is not false")
		assert.NotNil(t, response["remaining"], "remaining is nil")
		assert.Equal(t, response["remaining"], float64(52), "remaining is not 52")

		// check that cards key does not exist
		_, ok := response["cards"]
		assert.False(t, ok, "cards key exists")
	})

	t.Run("CreateCustomCardsDeckWithFaultyCard", func(t *testing.T) {
		req, err := http.NewRequest("POST", "/v1/decks?cards=AS,KH,BD", nil)

		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()

		handler := http.HandlerFunc(h.CreateDeck)

		handler.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusCreated, rr.Code, "handler returned wrong status code")
		assert.NotNil(t, rr.Body, "handler returned no body")

		var response map[string]interface{}
		err = json.Unmarshal(rr.Body.Bytes(), &response)

		if err != nil {
			t.Fatal(err)
		}

		assert.NotNil(t, response["deck_id"], "deck_id is nil")
		assert.NotNil(t, response["shuffled"], "shuffled is nil")
		assert.Equal(t, response["shuffled"], false, "shuffled is not false")
		assert.NotNil(t, response["remaining"], "remaining is nil")
		assert.Equal(t, response["remaining"], float64(2), "remaining is not 2")

		// check that cards key does not exist
		_, ok := response["cards"]
		assert.False(t, ok, "cards key exists")
	})
}

func TestGetDeckWithRemainingCards(t *testing.T) {
	query := database.New(db.DB)
	repo := repository.NewRepository(query)
	h := NewHandler(repo)

	t.Run("GetDeckWithRemainingCards", func(t *testing.T) {
		req, _ := http.NewRequest("POST", "/v1/decks", nil)
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(h.CreateDeck)

		handler.ServeHTTP(rr, req)
		server := httptest.NewServer(handler)
		defer server.Close()

		resp, err := http.Post(server.URL+"/v1/decks", "application/json", nil)
		if err != nil {
			t.Fatal(err)
		}

		var response map[string]interface{}
		err = json.NewDecoder(resp.Body).Decode(&response)
		if err != nil {
			t.Fatal(err)
		}

		deckID := response["deck_id"].(string)
		router := chi.NewRouter()
		router.Get("/v1/decks/{deckID}", h.GetDeckWithRemainingCards)

		req, _ = http.NewRequest("GET", "/v1/decks/"+deckID, nil)
		rr = httptest.NewRecorder()
		router.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code, "handler returned wrong status code")
		assert.NotNil(t, rr.Body, "handler returned no body")

		err = json.Unmarshal(rr.Body.Bytes(), &response)

		if err != nil {
			t.Fatal(err)
		}

		assert.NotNil(t, response["deck_id"], "deck_id is nil")
		assert.NotNil(t, response["shuffled"], "shuffled is nil")
		assert.NotNil(t, response["remaining"], "remaining is nil")
		assert.NotNil(t, response["cards"], "cards is nil")
		assert.Equal(t, len(response["cards"].([]interface{})), 52, "cards length is not 52")
	})
}

func TestDrawCards(t *testing.T) {
	query := database.New(db.DB)
	repo := repository.NewRepository(query)
	h := NewHandler(repo)

	t.Run("DrawCards", func(t *testing.T) {
		req, _ := http.NewRequest("POST", "/v1/decks", nil)
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(h.CreateDeck)

		handler.ServeHTTP(rr, req)
		server := httptest.NewServer(handler)
		defer server.Close()

		resp, err := http.Post(server.URL+"/v1/decks", "application/json", nil)
		if err != nil {
			t.Fatal(err)
		}

		var response map[string]interface{}
		err = json.NewDecoder(resp.Body).Decode(&response)
		if err != nil {
			t.Fatal(err)
		}

		deckID := response["deck_id"].(string)
		router := chi.NewRouter()
		router.Post("/v1/decks-draw", h.DrawCards)

		requestBody := map[string]interface{}{
			"deck_id": deckID,
			"count":   2,
		}

		requestBodyBytes, err := json.Marshal(requestBody)
		if err != nil {
			t.Fatal(err)
		}

		req, _ = http.NewRequest("POST", "/v1/decks-draw", bytes.NewBuffer(requestBodyBytes))
		rr = httptest.NewRecorder()
		router.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code, "handler returned wrong status code")
		assert.NotNil(t, rr.Body, "handler returned no body")

		err = json.Unmarshal(rr.Body.Bytes(), &response)

		if err != nil {
			t.Fatal(err)
		}
		assert.NotNil(t, response["cards"], "cards is nil")
		assert.Equal(t, len(response["cards"].([]interface{})), 2, "cards length is not 1")
	})
}
