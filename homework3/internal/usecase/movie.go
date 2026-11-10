package usecase

import (
	"context"
	"errors"
	"homework3/internal/domain"
)

type MovieProvider interface {
	GetMovie(ctx context.Context, id int) (domain.Movie, error)
}

type MovieUseCase interface {
	LoadMovies(ctx context.Context, from, to int) (domain.LoadMoviesResult, error)
	GetLastMovie(ctx context.Context) (domain.Movie, error)
}

type MovieService struct {
	provider MovieProvider
}

func NewMovieService(provider MovieProvider) *MovieService {
	return &MovieService{provider: provider}
}

func (s *MovieService) existsMovie(ctx context.Context, id int) (bool, error) {
	_, err := s.provider.GetMovie(ctx, id)
	if errors.Is(err, domain.ErrMovieNotFound) {
		return false, nil
	}
	return err == nil, err
}

func (s *MovieService) GetLastMovie(ctx context.Context) (domain.Movie, error) {
	low, high := 0, 1
	for {
		exists, err := s.existsMovie(ctx, high)
		if err != nil {
			return domain.Movie{}, err
		}
		if !exists {
			break
		}
		low = high
		high *= 2
	}
	for high-low > 1 {
		mid := low + (high-low)/2
		exists, err := s.existsMovie(ctx, mid)
		if err != nil {
			return domain.Movie{}, err
		}
		if exists {
			low = mid
		} else {
			high = mid
		}
	}
	if high == 1 {
		return domain.Movie{}, domain.ErrNoMovies
	}
	return s.provider.GetMovie(ctx, low)
}
