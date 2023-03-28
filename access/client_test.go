package access

import (
	"context"
	"net"
	"testing"

	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func TestNewClient(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	WithNoop(ctx, t, func(ctx context.Context, t *testing.T, dial func(context.Context, string) (net.Conn, error)) {
		conn, err := grpc.DialContext(ctx, "bufnet", grpc.WithContextDialer(dial), grpc.WithTransportCredentials(insecure.NewCredentials()))
		assert.NoError(t, err)

		defer conn.Close()

		c, err := NewClient(conn)
		assert.NoError(t, err)

		ok, err := c.Check(ctx, "urn:cloud:access:eu-central-1:12345678910:root", "urn:cloud:access:eu-central-1:12345678910:root", "access:changePassword")
		assert.NoError(t, err)
		assert.True(t, ok)
	})
}
