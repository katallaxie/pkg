package redux

import (
	"iter"
	"sync"

	"github.com/katallaxie/pkg/slices"
)

// Subscription is the type of the subscription of the FSM.
type Subscription[S State] interface {
	// Subscribe subscribes to the FSM.
	Subscribe() <-chan S
	// Cancel cancels the subscription.
	CancelSubscription(<-chan S)
}

// ActionPayload is the type of the payload of the action.
type ActionPayload interface{}

// Action is the type of the action of the FSM.
type ActionType int

// State is the type of the state of the FSM.
type State interface{}

// Reduce is the type of the reducer of the FSM.
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

// Reducer is the type of the reducer of the FSM.
type Reducer[S State] func(prev S, action Action) S

// Store is the type of the store of the FSM.
type Store[S State] interface {
	// Dispatch dispatches an event to the FSM.
	Dispatch(actions ...Action)
	// State gets the current state of the store.
	State(s ...S) S
	// Dispose disposes the store.
	Dispose()

	Subscription[S]
}

var _ Store[State] = (*store[State])(nil)

type store[S State] struct {
	state       S
	reducers    []Reducer[S]
	subscribers []chan S

	sync.RWMutex
}

// New creates a new store.
func New[S State](initialState S, reducers ...Reducer[S]) Store[S] {
	s := new(store[S])
	s.state = initialState
	s.reducers = slices.Append(reducers, s.reducers...)

	return s
}

// Dispatch dispatches an event to the FSM.
func (s *store[S]) Dispatch(actions ...Action) {
	s.Lock()
	defer s.Unlock()

	for _, action := range actions {
		for _, reducer := range s.reducers {
			s.state = reducer(s.state, action)

			for _, sub := range s.subscribers {
				go func(sub chan<- S) { // background
					sub <- s.state
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
func (s *store[S]) Subscribe() <-chan S {
	s.Lock()
	defer s.Unlock()

	newListener := make(chan S)
	s.subscribers = append(s.subscribers, newListener)

	return newListener
}

// CancelSubscription cancels the subscription.
func (s *store[S]) CancelSubscription(sub <-chan S) {
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
