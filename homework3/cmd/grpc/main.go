package main

import (
	"homework3/internal/client"
	"homework3/internal/config"
	moviegrpc "homework3/internal/transport/grpc"
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
	c := client.NewMovieClient(cfg.MovieBaseURL, cfg.UpstreamTimeout, cfg.OperationTimeout)
	grpcHandler := moviegrpc.NewHandler(c)
	grpcServer := moviegrpc.NewServer(grpcHandler)

	log.Printf("gRPC server listening on %s", cfg.GrpcAddr)

	if err := grpcServer.Serve(listener); err != nil {
		log.Fatal(err)
	}
}
