package router

import (
	"github.com/go-chi/chi"
	"github.com/go-chi/cors"
	"github.com/oksenGit/deck_api/internal/database"
	"github.com/oksenGit/deck_api/pkg/deck/db"
	"github.com/oksenGit/deck_api/pkg/deck/handler"
	"github.com/oksenGit/deck_api/pkg/deck/repository"
)

func Init() *chi.Mux {
	router := chi.NewRouter()

	router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"*"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300,
	}))

	v1Router := chi.NewRouter()
	db := db.Init()
	queries := database.New(db)
	repo := repository.NewRepository(queries)
	deckHandler := handler.NewHandler(repo)

	v1Router.Post("/decks", deckHandler.CreateDeck)
	v1Router.Get("/decks/{deckID}", deckHandler.GetDeckWithRemainingCards)
	v1Router.Post("/decks-draw", deckHandler.DrawCards)

	router.Mount("/v1", v1Router)
	return router
}
