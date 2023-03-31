package access

import (
	"context"
	"net"
	"testing"

	pb "github.com/katallaxie/pkg/proto"
	"github.com/katallaxie/pkg/urn"

	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
	"google.golang.org/grpc/test/bufconn"
)

const bufSize = 1024 * 1024

// WithNoop ...
func WithNoop(ctx context.Context, t *testing.T, access pb.AccessServer, f func(context.Context, *testing.T, func(context.Context, string) (net.Conn, error))) {
	lis := bufconn.Listen(bufSize)
	s := grpc.NewServer()

	pb.RegisterAccessServer(s, access)
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

type noopServer struct {
	Accessor
	pb.UnimplementedAccessServer
}

type noopPolicer []*Policy

func (n noopPolicer) Policies(principal *urn.URN) ([]*Policy, error) {
	return n, nil
}

type noopAccessor struct {
	Matcher
	Policer
}

func (n *noopAccessor) Allow(ctx context.Context, principal *urn.URN, ressource *urn.URN, action Action) (bool, error) {
	var allow bool // default to deny

	policies, err := n.Policies(principal)
	if err != nil {
		return allow, err
	}

	for _, p := range policies {
		for _, r := range p.Rules {
			for _, a := range r.Actions {
				if a == action {
					for _, rr := range r.Resources {
						u, err := urn.Parse(rr.String())
						if err != nil {
							return false, err
						}

						if !n.Matcher(u, ressource) {
							continue
						}

						allow = r.Effect == Allow
					}
				}
			}
		}
	}

	return allow, nil
}

func newNoopServer(matcher Matcher, policies ...*Policy) *noopServer {
	a := &noopAccessor{
		Matcher: matcher,
		Policer: noopPolicer(policies),
	}

	return &noopServer{a, pb.UnimplementedAccessServer{}}
}

// Check ...
func (s *noopServer) Check(ctx context.Context, req *pb.Check_Request) (*pb.Check_Response, error) {
	p, err := urn.FromProto(req.GetPrincipal())
	if err != nil {
		return nil, err
	}

	r, err := urn.FromProto(req.GetResource())
	if err != nil {
		return nil, err
	}

	allow, err := s.Allow(ctx, p, r, Action(req.GetAction()))
	if err != nil {
		return nil, err
	}

	return &pb.Check_Response{Allowed: allow}, nil
}
