package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net"
	"net/http"
	"os/signal"
	"syscall"

	"github.com/unfadingdawn/go-course/internal/app"
	"github.com/unfadingdawn/go-course/internal/config"
	moviegrpc "github.com/unfadingdawn/go-course/internal/handler/grpc"
	moviehttp "github.com/unfadingdawn/go-course/internal/handler/http"

	"google.golang.org/grpc"
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
		Handler:           moviehttp.NewRouter(application.Movies, application.Search, cfg.OperationTimeout),
		ReadHeaderTimeout: cfg.OperationTimeout,
	}
	defer server.Close()
	httpListener, err := net.Listen("tcp", cfg.HTTPAddr)
	if err != nil {
		return fmt.Errorf("listen HTTP: %w", err)
	}
	defer httpListener.Close()
	grpcListener, err := net.Listen("tcp", cfg.GrpcAddr)
	if err != nil {
		return fmt.Errorf("listen gRPC: %w", err)
	}
	defer grpcListener.Close()
	grpcHandler := moviegrpc.NewHandler(application.Movies, cfg.OperationTimeout)
	grpcServer := moviegrpc.NewServer(grpcHandler)
	defer grpcServer.Stop()

	errorsCh := make(chan error, 2)
	go func() {
		err := server.Serve(httpListener)
		if errors.Is(err, http.ErrServerClosed) {
			err = nil
		}
		if err != nil {
			err = fmt.Errorf("HTTP server failed: %w", err)
		}
		errorsCh <- err
	}()
	go func() {
		err := grpcServer.Serve(grpcListener)
		if errors.Is(err, grpc.ErrServerStopped) {
			err = nil
		}
		if err != nil {
			err = fmt.Errorf("gRPC server failed: %w", err)
		}
		errorsCh <- err
	}()
	log.Printf("HTTP server listening on %s", httpListener.Addr())
	log.Printf("gRPC server listening on %s", grpcListener.Addr())
	var serveErr error
	select {
	case serveErr = <-errorsCh:
	case <-ctx.Done():
	}
	shutdownCtx, cancel := context.WithTimeout(context.Background(), cfg.OperationTimeout)
	defer cancel()
	stopped := make(chan struct{})
	go func() {
		grpcServer.GracefulStop()
		close(stopped)
	}()
	shutdownErr := server.Shutdown(shutdownCtx)
	select {
	case <-stopped:
	case <-shutdownCtx.Done():
		grpcServer.Stop()
		<-stopped
	}
	return errors.Join(serveErr, shutdownErr)
}
