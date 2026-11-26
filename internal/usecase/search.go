package usecase

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/unfadingdawn/go-course/internal/domain"
)

type SearchRepository interface {
	GetMoviesByIDs(ctx context.Context, ids []int) ([]domain.Movie, error)
}

type SearchIndex interface {
	FindIDs(ctx context.Context, words []string) ([]int, error)
}

type SearchCache interface {
	Get(ctx context.Context, query string) ([]byte, bool, error)
	Set(ctx context.Context, query string, body []byte) error
}

type SearchUseCase interface {
	Search(ctx context.Context, query string) ([]byte, error)
}

type SearchService struct {
	store SearchRepository
	index SearchIndex
	cache SearchCache
}

func NewSearchService(store SearchRepository, index SearchIndex, cache SearchCache) *SearchService {
	return &SearchService{store: store, index: index, cache: cache}
}

func (s *SearchService) Search(ctx context.Context, query string) ([]byte, error) {
	query = strings.TrimSpace(query)
	if query == "" {
		return nil, domain.ErrEmptySearchQuery
	}

	body, found, err := s.cache.Get(ctx, query)
	if err != nil {
		return nil, err
	}
	if found {
		return body, nil
	}

	ids, err := s.index.FindIDs(ctx, tokenize(query))
	if err != nil {
		return nil, fmt.Errorf("search movie index: %w", err)
	}
	movies := make([]domain.Movie, 0)
	if len(ids) > 0 {
		movies, err = s.store.GetMoviesByIDs(ctx, ids)
		if err != nil {
			return nil, fmt.Errorf("load search results: %w", err)
		}
		if movies == nil {
			movies = []domain.Movie{}
		}
	}
	body, err = json.Marshal(movies)
	if err != nil {
		return nil, fmt.Errorf("encode search results: %w", err)
	}
	if err := s.cache.Set(ctx, query, body); err != nil {
		return nil, err
	}
	return body, nil
}
