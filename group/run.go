package group

import (
	"context"
	"errors"
	"fmt"
	"sync"
)

type token struct{}

// ErrUnimplemented is returned when a listener is not implemented.
var ErrUnimplemented = errors.New("server: unimplemented")

// Run is an abstraction on WaitGroup to run multiple functions concurrently.
// It mimics 'errgroup' to extend structs with functions to run concurrently with
// a root context.
type Run struct {
	ctx    context.Context
	cancel context.CancelFunc

	wg  sync.WaitGroup
	sem chan token

	sync.Mutex
}

// RunFunc is a function that is called to attach more routines
// to the server.
type RunFunc func(context.Context)

// Unimplemented is the default implementation.
type Unimplemented struct{}

// Run is running a new go routine
func (s *Unimplemented) Run(context.Context, RunFunc) error {
	return ErrUnimplemented
}

// WithContext is creating a new group with a context.
func WithContext(ctx context.Context) (*Run, context.Context) {
	ctx, cancel := context.WithCancel(ctx)

	s := new(Run)
	s.cancel = cancel
	s.ctx = ctx

	return s, ctx
}

// Run is creating a new go routine to run a function concurrently.
func (g *Run) Run(f RunFunc) error {
	if g.sem != nil {
		g.sem <- token{}
	}

	g.wg.Add(1)

	fn := func() {
		defer g.done()

		f(g.ctx)
	}

	go fn()

	return nil
}

// SetLimit limits the number of active listeners in this server
func (g *Run) SetLimit(n int) {
	if n < 0 {
		g.sem = nil
		return
	}

	if len(g.sem) != 0 {
		panic(fmt.Errorf("group: modify limit while %v routines run", len(g.sem)))
	}

	g.sem = make(chan token, n)
}

// Wait is waiting for all go routines to finish.
func (g *Run) Wait() {
	g.wg.Wait()
}

// Shutdown is waiting for all go routines to finish.
func (g *Run) Shutdown() error {
	g.cancel()
	g.wg.Wait() // wait for all routines to finish

	return nil
}

func (g *Run) done() {
	if g.sem != nil {
		<-g.sem
	}

	g.wg.Done()
}
