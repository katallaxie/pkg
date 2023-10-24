package access

import (
	"context"
	"net"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func TestNewClient(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	p := &Policy{
		Version: DefaultVersion,
		Rules: Rules{
			{
				ID:     "1",
				Effect: Allow,
				Resources: Resources{
					"urn:cloud:access:eu-central-1:12345678910:root",
				},
				Actions: Actions{
					"access:changePassword",
				},
			},
		},
	}

	noop := newNoopServer(IdentityBasedMatcher, p)

	WithNoop(ctx, t, noop, func(ctx context.Context, t *testing.T, dial func(context.Context, string) (net.Conn, error)) {
		conn, err := grpc.DialContext(ctx, "bufnet", grpc.WithContextDialer(dial), grpc.WithTransportCredentials(insecure.NewCredentials()))
		require.NoError(t, err)

		defer conn.Close()

		c, err := NewClient(conn)
		require.NoError(t, err)

		ok, err := c.Check(ctx, "urn:cloud:access:eu-central-1:12345678910:root", "urn:cloud:access:eu-central-1:12345678910:root", "access:changePassword")
		require.NoError(t, err)
		assert.True(t, ok)
	})
}

func BenchmarkNewClient(b *testing.B) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	p := &Policy{
		Version: DefaultVersion,
		Rules: Rules{
			{
				ID:     "1",
				Effect: Allow,
				Resources: Resources{
					"urn:cloud:access:eu-central-1:12345678910:root",
				},
				Actions: Actions{
					"access:changePassword",
				},
			},
		},
	}

	noop := newNoopServer(IdentityBasedMatcher, p)

	b.ReportAllocs()

	WithNoop(ctx, &testing.T{}, noop, func(ctx context.Context, t *testing.T, dial func(context.Context, string) (net.Conn, error)) {
		conn, err := grpc.DialContext(ctx, "bufnet", grpc.WithContextDialer(dial), grpc.WithTransportCredentials(insecure.NewCredentials()))
		require.NoError(t, err)

		defer conn.Close()

		c, err := NewClient(conn)
		require.NoError(t, err)

		for i := 0; i < b.N; i++ {
			ok, err := c.Check(ctx, "urn:cloud:access:eu-central-1:12345678910:root", "urn:cloud:access:eu-central-1:12345678910:root", "access:changePassword")
			require.NoError(t, err)
			assert.True(t, ok)
		}
	})
}
