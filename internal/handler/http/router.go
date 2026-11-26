package http

import (
	"net/http"
	"time"

	"github.com/unfadingdawn/go-course/internal/usecase"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func NewRouter(uc usecase.MovieUseCase, search usecase.SearchUseCase, operationTimeout time.Duration) http.Handler {
	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	h := NewHandler(uc, search, operationTimeout)

	r.Get("/load-movies", h.LoadMovies)
	r.Get("/last-movie", h.LastMovie)
	r.Get("/search", h.Search)
	r.Get("/swagger/index.html", swaggerIndex)
	r.Get("/swagger/openapi.json", swaggerSpec)

	return r
}
