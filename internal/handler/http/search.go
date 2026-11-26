package http

import (
	"context"
	"errors"
	"log"
	"net/http"
	"strings"

	"github.com/unfadingdawn/go-course/internal/domain"
)

func (h *Handler) Search(w http.ResponseWriter, r *http.Request) {
	query := strings.TrimSpace(r.URL.Query().Get("q"))
	if query == "" {
		http.Error(w, domain.ErrEmptySearchQuery.Error(), http.StatusBadRequest)
		return
	}
	ctx, cancel := context.WithTimeout(r.Context(), h.operationTimeout)
	defer cancel()
	body, err := h.search.Search(ctx, query)
	if err != nil {
		if errors.Is(err, domain.ErrEmptySearchQuery) {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		log.Printf("HTTP search request failed: %v", err)
		http.Error(w, "Search request failed", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if _, err := w.Write(body); err != nil {
		log.Printf("failed to write search response: %v", err)
	}
}
