package app

import (
	"context"
	"database/sql"
	"fmt"
	movieadapter "homework3/internal/adapter/movie"
	"homework3/internal/adapter/postgres"
	"homework3/internal/config"
	"homework3/internal/database"
	"homework3/internal/usecase"
	"log"
)

type App struct {
	Movies *usecase.MovieService
	db     *sql.DB
}

func New(ctx context.Context, cfg *config.Config) (*App, error) {
	dbCtx, cancel := context.WithTimeout(ctx, cfg.OperationTimeout)
	db, err := database.Open(dbCtx, cfg.DatabaseURL)
	cancel()
	if err != nil {
		return nil, err
	}
	repository := postgres.NewMovieRepository(db)
	client := movieadapter.NewMovieClient(cfg.MovieBaseURL, cfg.UpstreamTimeout)
	log.Print("Synchronizing movies with upstream")
	if err := usecase.NewMovieSynchronizer(client, repository).Sync(ctx); err != nil {
		db.Close()
		return nil, fmt.Errorf("synchronize movies: %w", err)
	}
	log.Print("Movie synchronization complete")
	return &App{Movies: usecase.NewMovieService(repository), db: db}, nil
}

func (a *App) Close() error {
	return a.db.Close()
}
