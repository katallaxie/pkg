package access

import (
	"context"
	"net"
	"testing"

	pb "github.com/katallaxie/pkg/proto"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
	"google.golang.org/grpc/test/bufconn"
)

const bufSize = 1024 * 1024

// WithNoop ...
func WithNoop(ctx context.Context, t *testing.T, f func(context.Context, *testing.T, func(context.Context, string) (net.Conn, error))) {
	lis := bufconn.Listen(bufSize)
	s := grpc.NewServer()

	pb.RegisterAccessServer(s, NewNoop())
	ctx, cancel := context.WithCancel(ctx)

	g, ctx := errgroup.WithContext(ctx)

	g.Go(func() error {
		if err := s.Serve(lis); err != nil {
			return err
		}

		return nil
	})

	g.Go(func() error {
		<-ctx.Done()
		s.GracefulStop()

		return nil
	})

	g.Go(func() error {
		f(ctx, t, func(context.Context, string) (net.Conn, error) { return lis.Dial() })

		cancel()

		return nil
	})

	err := g.Wait()
	if err != nil {
		t.Fatal(err)
	}
}

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
