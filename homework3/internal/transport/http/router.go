package http

import (
	"homework3/internal/client"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func NewRouter(c *client.MovieClient) http.Handler {
	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	h := NewHandler(c)

	r.Get("/load-movies", h.LoadMovies)
	r.Get("/last-movie", h.LastMovie)
	r.Get("/swagger/index.html", swaggerIndex)
	r.Get("/swagger/openapi.json", swaggerSpec)

	return r
}
