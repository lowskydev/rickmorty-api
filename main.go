package main

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"github.com/lowskydev/rickmorty-api/client"
	"github.com/lowskydev/rickmorty-api/handlers"
)

func main() {
	fmt.Println("Warming cache...")
	if _, err := client.GetAllCharacters(); err != nil {
		fmt.Printf("warning: could not pre-fetch characters: %v\n", err)
	}
	if _, err := client.GetAllEpisodes(); err != nil {
		fmt.Printf("warning: could not pre-fetch episodes: %v\n", err)
	}
	fmt.Println("Cache ready.")

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Get("/search", handlers.Search)
	r.Get("/top-pairs", handlers.TopPairs)

	fmt.Println("Server running on http://localhost:8080")
	http.ListenAndServe(":8080", r)
}
