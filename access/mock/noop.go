package mock

import (
	"context"

	pb "github.com/katallaxie/pkg/proto"
)

// NewNoop ...
func NewNoop() pb.AccessServer {
	n := new(noopServer)

	return n
}

type noopServer struct {
	pb.UnimplementedAccessServer
}

// Check ...
func (s *noopServer) Check(ctx context.Context, req *pb.Check_Request) (*pb.Check_Response, error) {
	return &pb.Check_Response{Allowed: true}, nil
}
