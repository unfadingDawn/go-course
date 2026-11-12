package main

import (
	"context"
	"errors"
	"fmt"
	"homework3/internal/app"
	"homework3/internal/config"
	moviehttp "homework3/internal/handler/http"
	"log"
	"net/http"
	"os/signal"
	"syscall"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()
	cfg, err := config.Load()
	if err != nil {
		return err
	}
	application, err := app.New(ctx, cfg)
	if err != nil {
		return err
	}
	defer application.Close()
	server := &http.Server{
		Addr:              cfg.HTTPAddr,
		Handler:           moviehttp.NewRouter(application.Movies, cfg.OperationTimeout),
		ReadHeaderTimeout: cfg.OperationTimeout,
	}
	defer server.Close()
	errorsCh := make(chan error, 1)
	go func() { errorsCh <- server.ListenAndServe() }()
	log.Printf("HTTP server listening on %s", cfg.HTTPAddr)
	select {
	case err := <-errorsCh:
		if errors.Is(err, http.ErrServerClosed) {
			return nil
		}
		return fmt.Errorf("HTTP server failed: %w", err)
	case <-ctx.Done():
		shutdownCtx, cancel := context.WithTimeout(context.Background(), cfg.OperationTimeout)
		defer cancel()
		return server.Shutdown(shutdownCtx)
	}
}
