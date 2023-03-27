package access

import (
	"context"

	pb "github.com/katallaxie/pkg/proto"
	"github.com/katallaxie/pkg/urn"

	"google.golang.org/grpc"
)

// Client is the interface for the access client.
type Access interface{}

// Client is the access client.
type Client struct {
	c pb.AccessClient
}

// Opt is the option for the access client.
type Opt func(*Client)

// NewClient creates a new access client.
func NewClient(conn *grpc.ClientConn, opts ...Opt) (*Client, error) {
	client := new(Client)

	for _, opt := range opts {
		opt(client)
	}

	c := pb.NewAccessClient(conn)
	client.c = c

	return client, nil
}

// Check is checking the access of a principal to a resource by an action.
func (c *Client) Check(ctx context.Context, principal, resource, action string) (bool, error) {
	p, err := urn.Parse(principal)
	if err != nil {
		return false, err
	}

	r, err := urn.Parse(resource)
	if err != nil {
		return false, err
	}

	res, err := c.c.Check(ctx, &pb.Check_Request{Principal: p.ProtoMessage(), Resource: r.ProtoMessage(), Action: action})
	if err != nil {
		return false, err
	}

	return res.GetAllowed(), nil
}
