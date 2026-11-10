package main

import (
	movieadapter "homework3/internal/adapter/movie"
	"homework3/internal/config"
	moviegrpc "homework3/internal/handler/grpc"
	"homework3/internal/usecase"
	"log"
	"net"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("%s", err.Error())
	}
	listener, err := net.Listen("tcp", cfg.GrpcAddr)
	if err != nil {
		log.Fatal(err)
	}
	c := movieadapter.NewMovieClient(cfg.MovieBaseURL, cfg.UpstreamTimeout)
	uc := usecase.NewMovieService(c)
	grpcHandler := moviegrpc.NewHandler(uc, cfg.OperationTimeout)
	grpcServer := moviegrpc.NewServer(grpcHandler)

	log.Printf("gRPC server listening on %s", cfg.GrpcAddr)

	if err := grpcServer.Serve(listener); err != nil {
		log.Fatal(err)
	}
}
