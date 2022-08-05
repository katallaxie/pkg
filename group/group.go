package group

import (
	"context"
	"sync"
)

// Run is an abstraction on WaitGroup to run multiple functions concurrently.
// It mimics 'errgroup' to extend structs with functions to run concurrently with
// a root context.
type Run struct {
	wg  sync.WaitGroup
	err *Error

	sync.Mutex
}

// Func ...
type Func func(ctx context.Context) error

// Run is creating a new go routine to run a function concurrently.
func (g *Run) Run(ctx context.Context, fn Func) {
	g.wg.Add(1)

	go func() {
		defer g.wg.Done()

		err := fn(ctx)
		if err != nil {
			g.Lock()
			g.err = Append(g.err, err)
			g.Unlock()
		}
	}()
}

// Wait is waiting for all go routines to finish.
func (g *Run) Wait() error {
	g.wg.Wait()
	g.Lock()
	defer g.Unlock()

	return g.err
}
