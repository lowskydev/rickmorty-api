package main

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"github.com/lowskydev/rickmorty-api/handlers"
)

func main() {
	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Get("/search", handlers.Search)

	fmt.Println("Server running on http://localhost:8080")
	http.ListenAndServe(":8080", r)
}
