package smtp

import (
	"context"
	"net"
	"net/textproto"
	"sync"

	"github.com/katallaxie/pkg/group"
)

// Session ...
type Session interface {
	// Create ...
	Create() group.RunFunc
}

type session struct {
	conn net.Conn
	text *textproto.Conn

	sync.Mutex
}

// Create
func (s *session) Create() group.RunFunc {
	return func(ctx context.Context) {
		defer func() {
			s.conn.Close()
		}()
	}
}

func newSession(conn net.Conn) *session {
	s := new(session)
	s.conn = conn
	s.text = textproto.NewConn(conn)

	return s
}
