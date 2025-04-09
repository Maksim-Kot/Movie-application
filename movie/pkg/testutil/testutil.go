package testutil

import (
	"github.com/Maksim-Kot/Movie-application/gen"
	"github.com/Maksim-Kot/Movie-application/movie/internal/controller/movie"
	metadatagateway "github.com/Maksim-Kot/Movie-application/movie/internal/gateway/metadata/grpc"
	ratinggateway "github.com/Maksim-Kot/Movie-application/movie/internal/gateway/rating/grpc"
	grpchandler "github.com/Maksim-Kot/Movie-application/movie/internal/handler/grpc"
	"github.com/Maksim-Kot/Movie-application/pkg/discovery"
)

// NewTestMovieGRPCServer creates a new movie gRPC server to be used in tests.
func NewTestMovieGRPCServer(registry discovery.Registry) gen.MovieServiceServer {
	metadataGateway := metadatagateway.New(registry)
	ratingGateway := ratinggateway.New(registry)
	ctrl := movie.New(ratingGateway, metadataGateway)
	return grpchandler.New(ctrl)
}
