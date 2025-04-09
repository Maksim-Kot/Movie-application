package testutil

import (
	"github.com/Maksim-Kot/Movie-application/gen"
	"github.com/Maksim-Kot/Movie-application/metadata/internal/controller/metadata"
	grpchandler "github.com/Maksim-Kot/Movie-application/metadata/internal/handler/grpc"
	"github.com/Maksim-Kot/Movie-application/metadata/internal/repository/memory"
)

// NewTestMetadataGRPCServer creates a new metadata gRPC server to be used in tests.
func NewTestMetadataGRPCServer() gen.MetadataServiceServer {
	r := memory.New()
	ctrl := metadata.New(r)
	return grpchandler.New(ctrl)
}
