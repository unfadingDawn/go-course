package usecase

import (
	"context"
	"homework3/internal/domain"
)

type MovieProvider interface {
	GetMovie(ctx context.Context, id int) (domain.Movie, error)
}

type MovieRepository interface {
	MovieProvider
	GetLastMovie(ctx context.Context) (domain.Movie, error)
}

type MovieUseCase interface {
	LoadMovies(ctx context.Context, from, to int) (domain.LoadMoviesResult, error)
	GetLastMovie(ctx context.Context) (domain.Movie, error)
}

type MovieService struct {
	provider MovieRepository
}

func NewMovieService(provider MovieRepository) *MovieService {
	return &MovieService{provider: provider}
}

func (s *MovieService) GetLastMovie(ctx context.Context) (domain.Movie, error) {
	return s.provider.GetLastMovie(ctx)
}
