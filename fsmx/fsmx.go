package fsmx

import (
	"sync"

	"github.com/katallaxie/pkg/slices"
)

// Payload is the type of the payload of the action.
type Payload interface{}

// Action is the type of the action of the FSM.
type Action int

// State is the type of the state of the FSM.
type State interface{}

// Actionable is the interface that wraps the basic Action methods.
type Actionable interface {
	// SetPayload sets the payload of the action
	SetPayload(payload Payload)
	// GetPayload gets the payload of the action
	GetPayload() Payload
	// SetType sets the type of the action
	SetType(actionType Action)
	// GetType gets the type of the action
	GetType() Action
}

// NewAction creates a new action.
func NewAction(actionType Action, payload Payload) Actionable {
	return &action{
		actionType: actionType,
		payload:    payload,
	}
}

type action struct {
	actionType Action
	payload    Payload
}

// GetPayload gets the payload of the action
func (a *action) GetPayload() Payload {
	return a.payload
}

// SetPayload sets the payload of the action
func (a *action) SetPayload(payload Payload) {
	a.payload = payload
}

// GetType gets the type of the action
func (a *action) GetType() Action {
	return a.actionType
}

// SetType sets the type of the action
func (a *action) SetType(actionType Action) {
	a.actionType = actionType
}

// Reducable is the type of the reducer of the FSM.
type Reducable func(prev State, action Actionable) State

// Storable is the interface that wraps the basic Store methods.
type Storable interface {
	// Dispatch dispatches an event to the FSM.
	Dispatch(action Actionable)
	// Subscribe subscribes to the store.
	Subscribe() <-chan State
	// Drain drains the store.
	Drain()
}

type store struct {
	state       State
	reducers    []Reducable
	subscribers []chan<- State

	sync.RWMutex
}

// New creates a new store.
func New(initialState State, reducers ...Reducable) Storable {
	s := new(store)
	s.state = initialState
	s.reducers = slices.Append(reducers, s.reducers...)

	return s
}

// Dispatch dispatches an event to the FSM.
func (s *store) Dispatch(action Actionable) {
	s.Lock()
	defer s.Unlock()

	for _, reducer := range s.reducers {
		s.state = reducer(s.state, action)

		for _, sub := range s.subscribers {
			go func(sub chan<- State) { // background
				sub <- s.state
			}(sub)
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
