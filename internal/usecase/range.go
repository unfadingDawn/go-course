package usecase

import (
	"context"
	"errors"
	"homework3/internal/domain"
	"log"
)

func (s *MovieService) LoadMovies(ctx context.Context, from, to int) (domain.LoadMoviesResult, error) {
	result := domain.LoadMoviesResult{Movies: make([]domain.Movie, 0), Errors: make([]domain.MovieLoadError, 0)}
	if from < 1 || to < from {
		return result, domain.ErrInvalidMoviesRange
	}
	for id := from; id <= to; id++ {
		if err := ctx.Err(); err != nil {
			return result, err
		}
		movie, err := s.provider.GetMovie(ctx, id)
		if ctx.Err() != nil {
			return result, ctx.Err()
		}
		if err != nil {
			log.Printf("LoadMovies: movie_id=%d from=%d to=%d: %v", id, from, to, err)
			failure := domain.MovieLoadError{ID: int64(id), Code: "unavailable", Message: "Movie service request failed"}
			if errors.Is(err, context.DeadlineExceeded) {
				failure.Code = "deadline_exceeded"
				failure.Message = "Movie service timed out"
			}
			result.Errors = append(result.Errors, failure)
		} else {
			result.Movies = append(result.Movies, movie)
		}
	}
	return result, nil
}
