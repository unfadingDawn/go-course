package main

import (
	"context"
	"homework3/internal/app"
	"homework3/internal/config"
	moviegrpc "homework3/internal/handler/grpc"
	"log"
	"net"
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
	listener, err := net.Listen("tcp", cfg.GrpcAddr)
	if err != nil {
		return err
	}
	defer listener.Close()
	grpcHandler := moviegrpc.NewHandler(application.Movies, cfg.OperationTimeout)
	grpcServer := moviegrpc.NewServer(grpcHandler)
	defer grpcServer.Stop()
	log.Printf("gRPC server listening on %s", cfg.GrpcAddr)
	errorsCh := make(chan error, 1)
	go func() { errorsCh <- grpcServer.Serve(listener) }()
	select {
	case err := <-errorsCh:
		return err
	case <-ctx.Done():
		shutdownCtx, cancel := context.WithTimeout(context.Background(), cfg.OperationTimeout)
		defer cancel()
		stopped := make(chan struct{})
		go func() {
			grpcServer.GracefulStop()
			close(stopped)
		}()
		select {
		case <-stopped:
		case <-shutdownCtx.Done():
			grpcServer.Stop()
		}
		return nil
	}
}
