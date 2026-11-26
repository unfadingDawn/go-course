package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/unfadingdawn/go-course/internal/domain"
)

type MovieRepository struct {
	db *sql.DB
}

func NewMovieRepository(db *sql.DB) *MovieRepository {
	return &MovieRepository{db: db}
}

func (r *MovieRepository) GetMovie(ctx context.Context, id int) (domain.Movie, error) {
	var movie domain.Movie
	err := r.db.QueryRowContext(ctx, "SELECT id, title, year, director FROM movies WHERE id = $1", id).
		Scan(&movie.ID, &movie.Title, &movie.Year, &movie.Director)
	if errors.Is(err, sql.ErrNoRows) {
		return domain.Movie{}, domain.ErrMovieNotFound
	}
	if err != nil {
		return domain.Movie{}, fmt.Errorf("get movie %d: %w", id, err)
	}
	return movie, nil
}

func (r *MovieRepository) GetLastMovie(ctx context.Context) (domain.Movie, error) {
	var movie domain.Movie
	err := r.db.QueryRowContext(ctx, "SELECT id, title, year, director FROM movies ORDER BY id DESC LIMIT 1").
		Scan(&movie.ID, &movie.Title, &movie.Year, &movie.Director)
	if errors.Is(err, sql.ErrNoRows) {
		return domain.Movie{}, domain.ErrNoMovies
	}
	if err != nil {
		return domain.Movie{}, fmt.Errorf("get last movie: %w", err)
	}
	return movie, nil
}

func (r *MovieRepository) GetAllMovies(ctx context.Context) ([]domain.Movie, error) {
	rows, err := r.db.QueryContext(ctx, "SELECT id, title, year, director FROM movies")
	if err != nil {
		return nil, fmt.Errorf("get all movies: %w", err)
	}
	defer rows.Close()

	movies, err := scanMovies(rows)
	if err != nil {
		return nil, fmt.Errorf("get all movies: %w", err)
	}
	return movies, nil
}

func (r *MovieRepository) GetMoviesByIDs(ctx context.Context, ids []int) ([]domain.Movie, error) {
	if len(ids) == 0 {
		return []domain.Movie{}, nil
	}

	rows, err := r.db.QueryContext(ctx,
		"SELECT id, title, year, director FROM movies WHERE id = ANY($1::bigint[])", ids)
	if err != nil {
		return nil, fmt.Errorf("get movies by IDs: %w", err)
	}
	defer rows.Close()

	movies, err := scanMovies(rows)
	if err != nil {
		return nil, fmt.Errorf("get movies by IDs: %w", err)
	}
	return movies, nil
}

func scanMovies(rows *sql.Rows) ([]domain.Movie, error) {
	movies := make([]domain.Movie, 0)
	for rows.Next() {
		var movie domain.Movie
		if err := rows.Scan(&movie.ID, &movie.Title, &movie.Year, &movie.Director); err != nil {
			return nil, fmt.Errorf("read movie: %w", err)
		}
		movies = append(movies, movie)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("read movies: %w", err)
	}
	return movies, nil
}

func (r *MovieRepository) ExistingIDs(ctx context.Context, lastID int) (map[int]bool, error) {
	rows, err := r.db.QueryContext(ctx, "SELECT id FROM movies WHERE id BETWEEN 1 AND $1", lastID)
	if err != nil {
		return nil, fmt.Errorf("list movie IDs: %w", err)
	}
	defer rows.Close()
	ids := make(map[int]bool)
	for rows.Next() {
		var id int
		if err := rows.Scan(&id); err != nil {
			return nil, fmt.Errorf("read movie ID: %w", err)
		}
		ids[id] = true
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("list movie IDs: %w", err)
	}
	return ids, nil
}

func (r *MovieRepository) SaveMovie(ctx context.Context, movie domain.Movie) error {
	_, err := r.db.ExecContext(ctx,
		"INSERT INTO movies (id, title, year, director) VALUES ($1, $2, $3, $4) ON CONFLICT (id) DO NOTHING",
		movie.ID, movie.Title, movie.Year, movie.Director)
	if err != nil {
		return fmt.Errorf("save movie %d: %w", movie.ID, err)
	}
	return nil
}
