package usecase

import (
	"context"
	"fmt"

	"github.com/unfadingdawn/go-course/internal/domain"
)

type MovieIndexer interface {
	Reset(ctx context.Context) error
	AddMovie(ctx context.Context, id int, words []string) error
}

type IndexRepository interface {
	GetAllMovies(ctx context.Context) ([]domain.Movie, error)
}

type MovieIndexBuilder struct {
	store IndexRepository
	index MovieIndexer
}

func NewMovieIndexBuilder(store IndexRepository, index MovieIndexer) *MovieIndexBuilder {
	return &MovieIndexBuilder{store: store, index: index}
}

func (b *MovieIndexBuilder) Build(ctx context.Context) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	movies, err := b.store.GetAllMovies(ctx)
	if err != nil {
		return fmt.Errorf("load movies for indexing: %w", err)
	}
	if err := b.index.Reset(ctx); err != nil {
		return fmt.Errorf("reset movie index and search cache: %w", err)
	}
	for _, movie := range movies {
		if err := ctx.Err(); err != nil {
			return err
		}
		words := tokenize(movie.Title + " " + movie.Director)
		if err := b.index.AddMovie(ctx, movie.ID, words); err != nil {
			return fmt.Errorf("build index for movie %d: %w", movie.ID, err)
		}
	}
	return nil
}
