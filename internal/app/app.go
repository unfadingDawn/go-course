package app

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"

	"github.com/redis/go-redis/v9"
	movieadapter "github.com/unfadingdawn/go-course/internal/adapter/movie"
	"github.com/unfadingdawn/go-course/internal/adapter/postgres"
	redisadapter "github.com/unfadingdawn/go-course/internal/adapter/redis"
	"github.com/unfadingdawn/go-course/internal/config"
	"github.com/unfadingdawn/go-course/internal/database"
	"github.com/unfadingdawn/go-course/internal/usecase"
)

type App struct {
	Movies *usecase.MovieService
	Search *usecase.SearchService
	db     *sql.DB
	redis  *redis.Client
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
	redisClient := redis.NewClient(&redis.Options{Addr: cfg.RedisAddr})
	redisCtx, cancel := context.WithTimeout(ctx, cfg.OperationTimeout)
	err = redisClient.Ping(redisCtx).Err()
	cancel()
	if err != nil {
		return nil, errors.Join(fmt.Errorf("connect to Redis: %w", err), redisClient.Close(), db.Close())
	}
	log.Print("Synchronizing movies with upstream")
	if err := usecase.NewMovieSynchronizer(client, repository).Sync(ctx); err != nil {
		return nil, errors.Join(fmt.Errorf("synchronize movies: %w", err), redisClient.Close(), db.Close())
	}
	log.Print("Movie synchronization complete")
	log.Print("Building movie index")
	index := redisadapter.NewMovieIndex(redisClient)
	if err := usecase.NewMovieIndexBuilder(repository, index).Build(ctx); err != nil {
		return nil, errors.Join(fmt.Errorf("build movie index: %w", err), redisClient.Close(), db.Close())
	}
	log.Print("Movie index complete")
	cache := redisadapter.NewSearchCache(redisClient)
	return &App{
		Movies: usecase.NewMovieService(repository),
		Search: usecase.NewSearchService(repository, index, cache),
		db:     db,
		redis:  redisClient,
	}, nil
}

func (a *App) Close() error {
	return errors.Join(a.db.Close(), a.redis.Close())
}
