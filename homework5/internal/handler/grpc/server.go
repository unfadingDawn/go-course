package grpc

import (
	moviepb "homework3/api/proto"

	"google.golang.org/grpc"

	"google.golang.org/grpc/reflection"
)

func NewServer(h *Handler) *grpc.Server {
	server := grpc.NewServer()

	moviepb.RegisterMovieServiceServer(server, h)

	reflection.Register(server)

	return server
}
