package grpc

import (
	moviepb "github.com/unfadingdawn/go-course/api/proto"

	"google.golang.org/grpc"

	"google.golang.org/grpc/reflection"
)

func NewServer(h *Handler) *grpc.Server {
	server := grpc.NewServer()

	moviepb.RegisterMovieServiceServer(server, h)

	reflection.Register(server)

	return server
}
