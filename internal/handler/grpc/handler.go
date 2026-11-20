package grpc

import (
	"context"
	"errors"
	"log"
	"time"

	moviepb "homework3/api/proto"
	"homework3/internal/domain"
	"homework3/internal/usecase"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Handler struct {
	moviepb.UnimplementedMovieServiceServer
	usecase          usecase.MovieUseCase
	operationTimeout time.Duration
}

func NewHandler(uc usecase.MovieUseCase, operationTimeout time.Duration) *Handler {
	return &Handler{usecase: uc, operationTimeout: operationTimeout}
}

func toProtoMovie(m domain.Movie) *moviepb.Movie {
	return &moviepb.Movie{
		Id:       int64(m.ID),
		Title:    m.Title,
		Year:     int32(m.Year),
		Director: m.Director,
	}
}

func (h *Handler) LoadMovies(
	ctx context.Context,
	req *moviepb.LoadMoviesRequest,
) (*moviepb.LoadMoviesResponse, error) {
	from := req.GetFrom()
	to := req.GetTo()
	ctx, cancel := context.WithTimeout(ctx, h.operationTimeout)
	defer cancel()
	result, err := h.usecase.LoadMovies(ctx, int(from), int(to))
	if err != nil {
		log.Printf("gRPC LoadMovies: from=%d to=%d: %v", from, to, err)
		return nil, toStatusError(err)
	}
	response := &moviepb.LoadMoviesResponse{}
	for _, movie := range result.Movies {
		response.Movies = append(response.Movies, toProtoMovie(movie))
	}
	for _, failure := range result.Errors {
		response.Errors = append(response.Errors, &moviepb.MovieLoadError{
			Id: failure.ID, Code: failure.Code, Message: failure.Message,
		})
	}
	return response, nil
}

func (h *Handler) GetLastMovie(
	ctx context.Context,
	_ *moviepb.GetLastMovieRequest,
) (*moviepb.Movie, error) {
	ctx, cancel := context.WithTimeout(ctx, h.operationTimeout)
	defer cancel()
	movie, err := h.usecase.GetLastMovie(ctx)
	if err != nil {
		log.Printf("gRPC GetLastMovie: %v", err)
		return nil, toStatusError(err)
	}
	return toProtoMovie(movie), nil
}

func toStatusError(err error) error {
	switch {
	case errors.Is(err, context.Canceled):
		return status.Error(codes.Canceled, "Operation canceled")
	case errors.Is(err, context.DeadlineExceeded):
		return status.Error(codes.DeadlineExceeded, "Deadline exceeded")
	case errors.Is(err, domain.ErrNoMovies):
		return status.Error(codes.NotFound, "Films not found")
	case errors.Is(err, domain.ErrInvalidMoviesRange):
		return status.Error(codes.InvalidArgument, "expected 1 <= from <= to")
	default:
		return status.Error(codes.Unavailable, "Service is unavailable")
	}
}
