package access

import (
	"context"
	"log"
	"net"
	"testing"

	"github.com/katallaxie/pkg/access/mock"
	pb "github.com/katallaxie/pkg/proto"

	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/test/bufconn"
)

const bufSize = 1024 * 1024

var lis *bufconn.Listener

func init() {
	lis = bufconn.Listen(bufSize)
	s := grpc.NewServer()
	pb.RegisterAccessServer(s, mock.NewNoop())

	go func() {
		if err := s.Serve(lis); err != nil {
			log.Fatalf("Server exited with error: %v", err)
		}
	}()
}

func bufDialer(context.Context, string) (net.Conn, error) {
	return lis.Dial()
}

func TestNewClient(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	conn, err := grpc.DialContext(ctx, "bufnet", grpc.WithContextDialer(bufDialer), grpc.WithTransportCredentials(insecure.NewCredentials()))
	assert.NoError(t, err)

	defer conn.Close()

	c, err := NewClient(conn)
	assert.NoError(t, err)

	ok, err := c.Check(ctx, "urn:cloud:access:eu-central-1:12345678910:root", "urn:cloud:access:eu-central-1:12345678910:root", "access:changePassword")
	assert.NoError(t, err)
	assert.True(t, ok)
}
