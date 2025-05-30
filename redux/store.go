package redux

import (
	"context"
	"iter"
	"sync"

	"github.com/katallaxie/pkg/slices"
	"github.com/katallaxie/pkg/utilx"
)

// Update is the type of the update of the store.
type Update interface{}

// Action is the type of the action of the store.
type Action func() Update

// Subscription is the type of the subscription of the store.
type Subscription[S State] interface {
	// Subscribe subscribes to the store.
	Subscribe() <-chan StateChange[S]
	// Cancel cancels the subscription.
	CancelSubscription(<-chan StateChange[S])
}

// State is the type of the state of the store.
type State interface{}

// QuitAction is the type of the quit action of the store.
func Quit() Update {
	return QuitUpdate{}
}

// QuitUpdate is the type of the quit update of the store.
type QuitUpdate struct{}

// Reduce is the type of the reducer of the store.
func Reduce[S State](curr S, reducers []Reducer[S], actions ...Action) iter.Seq[State] {
	return func(yield func(State) bool) {
		for _, action := range actions {
			for _, reducer := range reducers {
				curr = reducer(curr, action)
				if !yield(curr) {
					return
				}
			}
		}
	}
}

// Reducer is the type of the reducer of the store.
type Reducer[S State] func(curr S, update Update) S

// Store is the type of the store.
type Store[S State] interface {
	// Dispatch dispatches an event to the store.
	Dispatch(actions ...Action)
	// State gets the current state of the store.
	State(s ...S) S
	// Dispose disposes the store.
	Dispose()

	Subscription[S]
}

// StateChange is the type of the state change of the store.
type StateChange[S State] interface {
	// Prev gets the previous state of the store.
	Prev() S
	// Curr gets the current state of the store.
	Curr() S
}

var _ StateChange[any] = (*stateChange[any])(nil)

type stateChange[S State] struct {
	prev S
	curr S
}

// NewStateChange creates a new state change.
func NewStateChange[S State](prev, curr S) StateChange[S] {
	return &stateChange[S]{
		prev: prev,
		curr: curr,
	}
}

// Curr gets the current state of the store.
func (s *stateChange[S]) Curr() S {
	return s.curr
}

// Prev gets the previous state of the store.
func (s *stateChange[S]) Prev() S {
	return s.prev
}

var _ Store[State] = (*store[State])(nil)

type store[S State] struct {
	state       S
	reducers    []Reducer[S]
	subscribers []chan StateChange[S]
	updates     chan Update
	actions     chan Action
	ctx         context.Context

	sync.RWMutex
}

// New creates a new store.
func New[S State](ctx context.Context, initialState S, reducers ...Reducer[S]) Store[S] {
	s := new(store[S])
	s.ctx = ctx
	s.state = initialState
	s.actions = make(chan Action)
	s.updates = make(chan Update)
	s.reducers = slices.Append(reducers, s.reducers...)

	s.handleActions(s.actions)
	s.handleUpdates(s.updates)

	return s
}

// Dispatch dispatches an event to the store.
func (s *store[S]) Dispatch(actions ...Action) {
	go func() {
		for _, action := range actions {
			if utilx.IsNil(action) {
				continue
			}

			s.actions <- action
		}
	}()
}

// State gets the current state of the store.
func (s *store[S]) State(states ...S) S {
	s.Lock()
	defer s.Unlock()

	if slices.GreaterThen(0, states...) {
		s.state = slices.First(states...)
	}

	return s.state
}

// Subscribe subscribes to the store.
func (s *store[S]) Subscribe() <-chan StateChange[S] {
	s.Lock()
	defer s.Unlock()

	newListener := make(chan StateChange[S])
	s.subscribers = append(s.subscribers, newListener)

	return newListener
}

// CancelSubscription cancels the subscription.
func (s *store[S]) CancelSubscription(sub <-chan StateChange[S]) {
	s.Lock()
	defer s.Unlock()

	for i, l := range s.subscribers {
		if l == sub {
			s.subscribers = slices.Delete(i, s.subscribers...)
			close(l)
		}
	}
}

// Dispose disposes the store.
func (s *store[S]) Dispose() {
	s.Lock()
	defer s.Unlock()

	for _, sub := range s.subscribers {
		close(sub)
	}
}

func (s *store[S]) handleUpdates(updates chan Update) chan struct{} {
	ch := make(chan struct{})

	go func() {
		defer close(ch)

		for {
			select {
			case <-s.ctx.Done():
				return
			case update := <-updates:
				if utilx.Empty(update) {
					continue
				}

				prev := s.state

				for _, reducer := range s.reducers {
					s.state = reducer(s.state, update)
				}

				for _, sub := range s.subscribers {
					go func() {
						change := NewStateChange(prev, s.state)
						sub <- change
					}()
				}
			}
		}
	}()

	return ch
}

func (s *store[S]) handleActions(actions chan Action) chan struct{} {
	ch := make(chan struct{})

	go func() {
		defer close(ch)

		for {
			select {
			case <-s.ctx.Done():
				return
			case action := <-actions:
				if action == nil {
					continue
				}

				update := action()
				s.updates <- update
			}
		}
	}()

	return ch
}
