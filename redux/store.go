package redux

import (
	"fmt"
	"iter"
	"sync"

	"github.com/katallaxie/pkg/slices"
)

// Subscription is the type of the subscription of the store.
type Subscription[S State] interface {
	// Subscribe subscribes to the store.
	Subscribe() <-chan StateChange[S]
	// Cancel cancels the subscription.
	CancelSubscription(<-chan StateChange[S])
}

// ActionPayload is the type of the payload of the action.
type ActionPayload interface{}

// Action is the type of the action of the store.
type ActionType int

// State is the type of the state of the store.
type State interface{}

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

// Actionable is the interface that wraps the basic Action methods.
type Action interface {
	// Paylog gets the payload of the action
	Payload(p ...ActionPayload) ActionPayload
	// Type gets the type of the action
	Type(a ...ActionType) ActionType
}

// NewAction creates a new action.
func NewAction(actionType ActionType, payload ActionPayload) Action {
	return &action{
		actionType: actionType,
		payload:    payload,
	}
}

type action struct {
	actionType ActionType
	payload    ActionPayload
}

// Payload gets the payload of the action
func (a *action) Payload(payloads ...ActionPayload) ActionPayload {
	if slices.Len(payloads) > 0 {
		a.payload = slices.First(payloads...)
	}

	return a.payload
}

// Type gets the type of the action
func (a *action) Type(types ...ActionType) ActionType {
	if slices.Len(types) > 0 {
		a.actionType = slices.First(types...)
	}

	return a.actionType
}

// Reducer is the type of the reducer of the store.
type Reducer[S State] func(prev S, action Action) S

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

	sync.RWMutex
}

// New creates a new store.
func New[S State](initialState S, reducers ...Reducer[S]) Store[S] {
	s := new(store[S])
	s.state = initialState
	s.reducers = slices.Append(reducers, s.reducers...)

	return s
}

// Dispatch dispatches an event to the store.
func (s *store[S]) Dispatch(actions ...Action) {
	s.Lock()
	defer s.Unlock()

	for _, action := range actions {
		for _, reducer := range s.reducers {
			prev := s.state
			s.state = reducer(s.state, action)

			for _, sub := range s.subscribers {
				go func(sub chan<- StateChange[S]) { // background
					sub <- NewStateChange(prev, s.state)
				}(sub)
			}
		}
	}
}

// State gets the current state of the store.
func (s *store[S]) State(states ...S) S {
	s.Lock()
	defer s.Unlock()

	if slices.Len(states) > 0 {
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

	fmt.Println("cancel subscription")

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
