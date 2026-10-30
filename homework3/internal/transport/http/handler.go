package http

import (
	"context"
	"encoding/json"
	"errors"
	"homework3/internal/client"
	"log"
	"net/http"
	"strconv"
)

type Handler struct {
	client *client.MovieClient
}

func NewHandler(c *client.MovieClient) *Handler {
	return &Handler{client: c}
}

func (h *Handler) LoadMovies(w http.ResponseWriter, r *http.Request) {
	from, err := strconv.Atoi(r.URL.Query().Get("from"))
	if err != nil {
		http.Error(w, "invalid from parameter", http.StatusBadRequest)
		return
	}
	to, err := strconv.Atoi(r.URL.Query().Get("to"))
	if err != nil {
		http.Error(w, "invalid to parameter", http.StatusBadRequest)
		return
	}
	if from < 1 || to < from {
		http.Error(w, "expected 1 <= from <= to", http.StatusBadRequest)
		return
	}
	ctx, cancel := context.WithTimeout(r.Context(), h.client.OperationTimeout)
	defer cancel()
	result, err := h.client.LoadMovies(ctx, int64(from), int64(to))
	if err != nil {
		writeUpstreamError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	var body any = result.Movies
	statusCode := http.StatusOK
	if len(result.Errors) > 0 {
		body = result
		statusCode = http.StatusBadGateway
		for _, failure := range result.Errors {
			if failure.Code == "deadline_exceeded" {
				statusCode = http.StatusGatewayTimeout
				break
			}
		}
	}
	w.WriteHeader(statusCode)
	if err := json.NewEncoder(w).Encode(body); err != nil {
		log.Printf("failed to encode response: %v", err)
	}
}

func (h *Handler) LastMovie(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), h.client.OperationTimeout)
	defer cancel()
	movie, err := h.client.GetLastMovie(ctx)
	if err != nil {
		writeUpstreamError(w, err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(movie); err != nil {
		log.Printf("failed to encode response: %v", err)
	}
}

func writeUpstreamError(w http.ResponseWriter, err error) {
	log.Printf("HTTP movie request failed: %v", err)
	switch {
	case errors.Is(err, context.Canceled):
		return
	case errors.Is(err, context.DeadlineExceeded):
		http.Error(w, "Movie service timed out", http.StatusGatewayTimeout)
	case errors.Is(err, client.ErrNoMovies):
		http.Error(w, "Films not found", http.StatusNotFound)
	default:
		http.Error(w, "Movie service request failed", http.StatusBadGateway)
	}
}
