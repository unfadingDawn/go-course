package client

import (
	"context"
	"errors"
	"homework3/internal/model"
	"log"
)

type MovieLoadError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	ID      int64  `json:"id"`
}

type LoadMoviesResult struct {
	Movies []model.Movie    `json:"movies"`
	Errors []MovieLoadError `json:"errors"`
}

func (c *MovieClient) LoadMovies(ctx context.Context, from, to int64) (LoadMoviesResult, error) {
	result := LoadMoviesResult{Movies: make([]model.Movie, 0), Errors: make([]MovieLoadError, 0)}
	for id := from; id <= to; id++ {
		if err := ctx.Err(); err != nil {
			return result, err
		}
		movie, err := c.GetMovie(ctx, int(id))
		if ctx.Err() != nil {
			return result, ctx.Err()
		}
		if err != nil {
			log.Printf("LoadMovies: movie_id=%d from=%d to=%d: %v", id, from, to, err)
			failure := MovieLoadError{ID: id, Code: "unavailable", Message: "Movie service request failed"}
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
