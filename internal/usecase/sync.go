package usecase

import (
	"context"
	"errors"
	"fmt"
	"homework3/internal/domain"
)

type SyncRepository interface {
	ExistingIDs(ctx context.Context, lastID int) (map[int]bool, error)
	SaveMovie(ctx context.Context, movie domain.Movie) error
}

type MovieSynchronizer struct {
	source MovieProvider
	store  SyncRepository
}

func NewMovieSynchronizer(source MovieProvider, store SyncRepository) *MovieSynchronizer {
	return &MovieSynchronizer{source: source, store: store}
}

func (s *MovieSynchronizer) Sync(ctx context.Context) error {
	cache := make(map[int]domain.Movie)
	get := func(id int) (domain.Movie, error) {
		if movie, ok := cache[id]; ok {
			return movie, nil
		}
		movie, err := s.source.GetMovie(ctx, id)
		if err != nil {
			return domain.Movie{}, err
		}
		if movie.ID != id {
			return domain.Movie{}, fmt.Errorf("upstream movie ID mismatch: requested %d, received %d", id, movie.ID)
		}
		cache[id] = movie
		return movie, nil
	}
	exists := func(id int) (bool, error) {
		_, err := get(id)
		if errors.Is(err, domain.ErrMovieNotFound) {
			return false, nil
		}
		return err == nil, err
	}
	low, high := 0, 1
	for {
		found, err := exists(high)
		if err != nil {
			return fmt.Errorf("discover last movie: %w", err)
		}
		if !found {
			break
		}
		low = high
		if high > int(^uint(0)>>1)/2 {
			return fmt.Errorf("upstream movie ID exceeds supported range")
		}
		high *= 2
	}
	for high-low > 1 {
		mid := low + (high-low)/2
		found, err := exists(mid)
		if err != nil {
			return fmt.Errorf("discover last movie: %w", err)
		}
		if found {
			low = mid
		} else {
			high = mid
		}
	}
	ids, err := s.store.ExistingIDs(ctx, low)
	if err != nil {
		return err
	}
	for id := 1; id <= low; id++ {
		if err := ctx.Err(); err != nil {
			return err
		}
		if ids[id] {
			delete(cache, id)
			continue
		}
		movie, err := get(id)
		if err != nil {
			return fmt.Errorf("download movie %d: %w", id, err)
		}
		if err := s.store.SaveMovie(ctx, movie); err != nil {
			return err
		}
		delete(cache, id)
	}
	return nil
}
