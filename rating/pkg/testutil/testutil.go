package testutil

import (
	"github.com/Maksim-Kot/Movie-application/gen"
	"github.com/Maksim-Kot/Movie-application/rating/internal/controller/rating"
	grpchandler "github.com/Maksim-Kot/Movie-application/rating/internal/handler/grpc"
	"github.com/Maksim-Kot/Movie-application/rating/internal/repository/memory"
)

// NewTestRatingGRPCServer creates a new rating gRPC server to be used in tests.
func NewTestRatingGRPCServer() gen.RatingServiceServer {
	r := memory.New()
	ctrl := rating.New(r, nil)
	return grpchandler.New(ctrl)
}
