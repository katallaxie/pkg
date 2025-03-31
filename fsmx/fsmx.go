package fsmx

import (
	"sync"

	"github.com/katallaxie/pkg/slices"
)

// ActionPayload is the type of the payload of the action.
type ActionPayload interface{}

// Action is the type of the action of the FSM.
type ActionType int

// State is the type of the state of the FSM.
type State interface{}

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
type Reducer func(prev State, action Action) State

// Store is the type of the store of the FSM.
type Store interface {
	// Dispatch dispatches an event to the FSM.
	Dispatch(actions ...Action)
	// Subscribe subscribes to the store.
	Subscribe() <-chan State
	// Drain drains the store.
	Drain()
}

type store struct {
	state       State
	reducers    []Reducer
	subscribers []chan<- State

	sync.RWMutex
}

// New creates a new store.
func New(initialState State, reducers ...Reducer) Store {
	s := new(store)
	s.state = initialState
	s.reducers = slices.Append(reducers, s.reducers...)

	return s
}

// Dispatch dispatches an event to the FSM.
func (s *store) Dispatch(actions ...Action) {
	s.Lock()
	defer s.Unlock()

	for _, action := range actions {
		for _, reducer := range s.reducers {
			s.state = reducer(s.state, action)

			for _, sub := range s.subscribers {
				go func(sub chan<- State) { // background
					sub <- s.state
				}(sub)
			}
		}
	}
}

// Subscribe subscribes to the store.
func (s *store) Subscribe() <-chan State {
	s.Lock()
	defer s.Unlock()

	sub := make(chan State)
	s.subscribers = slices.Append(s.subscribers, sub)

	return sub
}

// Drain drains the store.
func (s *store) Drain() {
	s.Lock()
	defer s.Unlock()

	for _, sub := range s.subscribers {
		close(sub)
	}
}
