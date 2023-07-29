package smtp_test

import (
	"context"
	"fmt"
	"net"
	s "net/smtp"
	"testing"
	"time"

	"github.com/katallaxie/pkg/group"
	"github.com/katallaxie/pkg/smtp"
	"github.com/stretchr/testify/assert"
)

func TestSessionServe(t *testing.T) {
	l, err := net.Listen("tcp", "127.0.0.1:0")
	assert.NoError(t, err)
	assert.NotNil(t, l)

	ctx, cancel := context.WithCancel(context.Background())
	g, _ := group.WithContext(ctx)

	ready := make(chan struct{})

	go func() {
		<-ready

		conn, err := net.DialTimeout("tcp", l.Addr().String(), 1*time.Second)
		assert.NoError(t, err)

		c, err := s.NewClient(conn, "localhost")
		fmt.Println("hello")

		assert.NoError(t, err)
		assert.NotNil(t, c)

		c.Close()
		cancel()
	}()

	close(ready)

	conn, err := l.Accept()
	assert.NoError(t, err)
	assert.NotNil(t, conn)

	s := smtp.NewSession(conn)
	assert.NotNil(t, s)

	g.Run(s.Serve())
	g.Wait()
}
