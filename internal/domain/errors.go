package domain

import "errors"

var (
	ErrMovieNotFound      = errors.New("movie not found")
	ErrNoMovies           = errors.New("films not found")
	ErrInvalidMoviesRange = errors.New("expected 1 <= from <= to")
)
